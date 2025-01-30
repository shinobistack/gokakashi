package db

import (
	"context"
	"github.com/shinobistack/gokakashi/ent/integrations"
	"github.com/shinobistack/gokakashi/ent/policies"
	"github.com/shinobistack/gokakashi/ent/policylabels"
	"github.com/shinobistack/gokakashi/ent/scans"
	"github.com/shinobistack/gokakashi/ent/schema"
	"log"

	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/internal/config/v1"
)

func RunMigrations(client *ent.Client) {
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("Failed to create database schema: %v", err)
	}
	log.Println("Database schema created successfully")
}

func PopulateDatabase(client *ent.Client, cfg *v1.Config) {
	log.Println("Populating database from configuration...")

	// Populate integrations
	for _, integration := range cfg.Integrations {
		existing, err := client.Integrations.Query().
			Where(integrations.Name(integration.Name)).
			Only(context.Background())
		if err == nil && existing != nil {
			// Update existing integration
			_, err := client.Integrations.UpdateOne(existing).
				SetType(integration.Type).
				SetConfig(integration.Config).
				Save(context.Background())
			if err != nil {
				log.Printf("Failed to update integration %s: %v", integration.Name, err)
			} else {
				log.Printf("Integration %s updated successfully.", integration.Name)
			}
			continue
		} else if !ent.IsNotFound(err) {
			log.Printf("Error querying integration %s: %v", integration.Name, err)
			continue
		}

		// Create new integration
		_, err = client.Integrations.
			Create().
			SetName(integration.Name).
			SetType(integration.Type).
			SetConfig(integration.Config).
			Save(context.Background())
		if err != nil {
			log.Printf("Failed to add integration %s: %v", integration.Name, err)
		} else {
			log.Printf("Integration %s added successfully.", integration.Name)
		}
	}

	// Populate policies and associated scans
	for _, policy := range cfg.Policies {
		existing, err := client.Policies.Query().Where(policies.Name(policy.Name)).Only(context.Background())
		// Dynamically populate Type: cron|ci
		triggerData := map[string]interface{}{"type": policy.Trigger.Type}
		if policy.Trigger.Type == "cron" {
			triggerData["schedule"] = policy.Trigger.Schedule
		}

		if err == nil && existing != nil {
			// Update existing policy
			_, err := client.Policies.UpdateOne(existing).
				SetImage(schema.Image(policy.Image)).
				SetTrigger(triggerData).
				// ToDo: to update the scanner field to take in tools and tool's argument
				SetScanner(policy.Scanner).
				SetNotify(policy.Notify).
				Save(context.Background())
			if err != nil {
				log.Printf("Failed to update policy %s: %v", policy.Name, err)
			} else {
				log.Printf("Policy %s updated successfully.", policy.Name)
			}
			continue
		} else if !ent.IsNotFound(err) {
			log.Printf("Error querying policy %s: %v", policy.Name, err)
			continue
		}

		// Create new policy
		policyRecord, err := client.Policies.
			Create().
			SetName(policy.Name).
			SetImage(schema.Image(policy.Image)).
			SetTrigger(triggerData).
			SetScanner(policy.Scanner).
			SetNotify(policy.Notify).
			Save(context.Background())
		if err != nil {
			log.Printf("Failed to add policy %s: %v", policy.Name, err)
			continue
		}

		// Populate Policy Labels
		for key, value := range policy.Labels {
			existingLabel, err := client.PolicyLabels.Query().
				Where(policylabels.PolicyID(policyRecord.ID), policylabels.Key(key)).Only(context.Background())
			if err == nil && existingLabel != nil {
				// Update existing label
				_, err := client.PolicyLabels.UpdateOne(existingLabel).
					SetValue(value).
					Save(context.Background())
				if err != nil {
					log.Printf("Failed to update label %s:%s for policy %s: %v", key, value, policy.Name, err)
				} else {
					log.Printf("Label %s:%s updated for policy %s.", key, value, policy.Name)
				}
				continue
			} else if !ent.IsNotFound(err) {
				log.Printf("Error querying label %s for policy %s: %v", key, policy.Name, err)
				continue
			}

			// Create new label
			_, err = client.PolicyLabels.
				Create().
				SetPolicyID(policyRecord.ID).
				SetKey(key).
				SetValue(value).
				Save(context.Background())
			if err != nil {
				log.Printf("Failed to add label %s:%s for policy %s: %v", key, value, policy.Name, err)
			} else {
				log.Printf("Label %s:%s added for policy %s.", key, value, policy.Name)
			}
		}

		if policy.Trigger.Type == "cron" {
			// Populate scans for each policy image tag
			// ToDo: Runs periodic scans on a predefined set of images
			for _, tag := range policy.Image.Tags {
				existingScan, err := client.Scans.Query().
					Where(scans.PolicyID(policyRecord.ID), scans.Image(policy.Image.Name+":"+tag)).
					Only(context.Background())
				if err == nil && existingScan != nil {
					log.Printf("Scan for policy %s, image:tag %s:%s already exists. Skipping.", policy.Name, policy.Image.Name, tag)
					continue
				} else if !ent.IsNotFound(err) {
					log.Printf("Error querying scan for policy %s, image:tag %s:%s: %v", policy.Name, policy.Image.Name, tag, err)
					continue
				}

				existingIntegration, err := client.Integrations.Query().
					Where(integrations.Name(policy.Image.Registry)).
					OnlyID(context.Background())
				if err != nil {
					log.Printf("Error querying integrationID for policy %s and integrationName %s : %v", policy.Name, policy.Image.Registry, err)
					continue
				}

				scanCreate := client.Scans.Create().
					SetPolicyID(policyRecord.ID).
					SetImage(policy.Image.Name + ":" + tag).
					SetScanner(policy.Scanner).
					SetIntegrationID(existingIntegration)

				// Map notify.to to IntegrationID
				for _, notify := range policy.Notify {
					notifyIntegration, err := client.Integrations.Query().
						Where(integrations.Name(notify.To)).
						OnlyID(context.Background())
					//log.Printf("CHECK: notifyIntegration: %s", notifyIntegration)
					//log.Printf("CHECK: notify:%s", notify)
					if err != nil {
						log.Printf("Error finding integration for notify.to: %s, policy: %s: %v", notify.To, policy.Name, err)
						continue
					}
					scanCreate.SetNotify([]schema.Notify{
						{To: notifyIntegration.String(), When: notify.When, Format: notify.Format},
					})
				}

				_, err = scanCreate.Save(context.Background())
				if err != nil {
					log.Printf("Failed to add scan for policy %s, tag %s: %v", policy.Name, tag, err)
				} else {
					log.Printf("Scan for policy %s, tag %s added successfully.", policy.Name, tag)
				}
			}
		}
	}
}
