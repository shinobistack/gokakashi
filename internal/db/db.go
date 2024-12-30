package db

import (
	"context"
	"github.com/shinobistack/gokakashi/ent/schema"
	"log"

	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/internal/config/v1"
)

func PopulateDatabase(client *ent.Client, cfg *v1.Config) {
	log.Println("Populating database from configuration...")

	// Populate integrations
	for _, integration := range cfg.Integrations {
		_, err := client.Integrations.
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
		// ToDo: To remove post ent/policies.go line 48
		// Convert Trigger to map[string]interface{}
		triggerMap := map[string]interface{}{
			"type":     policy.Trigger.Type,
			"schedule": policy.Trigger.Schedule,
		}

		// Create Policy record
		policyRecord, err := client.Policies.
			Create().
			SetName(policy.Name).
			SetImage(schema.Image(policy.Image)).
			SetTrigger(triggerMap).
			SetCheck(schema.Check(policy.Check)).
			Save(context.Background())
		if err != nil {
			log.Printf("Failed to add policy %s: %v", policy.Name, err)
			continue
		}

		// Populate Policy Labels
		for key, value := range policy.Labels {
			_, err := client.PolicyLabels.
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

		// Populate scans for each policy image tag
		for _, tag := range policy.Image.Tags {
			_, err := client.Scans.
				Create().
				SetPolicyID(policyRecord.ID).
				SetImage(policy.Image.Registry + "/" + policy.Image.Name + ":" + tag).
				// SetStatus defaults to scan_pending
				Save(context.Background())
			if err != nil {
				log.Printf("Failed to add scan for policy %s, tag %s: %v", policy.Name, tag, err)
			} else {
				log.Printf("Scan for policy %s, tag %s added successfully.", policy.Name, tag)
			}
		}
	}
}
