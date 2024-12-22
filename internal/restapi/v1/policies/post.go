package policies

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/policies"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/swaggest/usecase/status"
	"regexp"
	"strings"
)

type CreatePolicyRequest struct {
	Name   string               `json:"name"`
	Image  schema.Image         `json:"image"`
	Labels []schema.PolicyLabel `json:"labels"`
	// Todo: Implement the logic of Type:cron etc
	Trigger map[string]interface{} `json:"trigger"`
	Check   *schema.Check          `json:"check"`
}

type CreatePolicyResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func CreatePolicy(client *ent.Client) func(ctx context.Context, req CreatePolicyRequest, res *CreatePolicyResponse) error {
	return func(ctx context.Context, req CreatePolicyRequest, res *CreatePolicyResponse) error {
		// Validate image fields
		if req.Image.Registry == "" || req.Image.Name == "" || len(req.Image.Tags) == 0 {
			return status.Wrap(errors.New("invalid image: missing required fields"), status.InvalidArgument)
		}

		if !isValidID(req.Name) {
			return status.Wrap(errors.New("invalid id format; must be lowercase, alphanumeric, or dashes"), status.InvalidArgument)
		}

		// Check for duplicate name
		exists, err := client.Policies.Query().
			Where(policies.Name(req.Name)).
			Exist(ctx)
		if err != nil {
			return status.Wrap(fmt.Errorf("failed to check for duplicate policies name: %v", err), status.Internal)
		}
		if exists {
			return status.Wrap(errors.New("policy with the same name already exists"), status.AlreadyExists)
		}

		// Validate trigger
		// ToDo: Valid cron for type: cron

		// Validate check fields
		// ToDo: Discuss, if check is mentioned then validate the check and accordingly publish the report by default or publish complete report.
		// ToDo: Discuss, if check is present then post satisfying the condition execute Notify accordingly
		// ToDO: Discuss the report server
		// ToDo: check.condition format for evaluation
		if req.Check.Condition == "" || len(req.Check.Notify) == 0 {
			return status.Wrap(errors.New("invalid check: missing required fields"), status.InvalidArgument)
		}

		tx, err := client.Tx(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		// Save the policy
		policy, err := tx.Policies.Create().
			SetName(req.Name).
			SetImage(req.Image).
			SetTrigger(req.Trigger).
			SetNillableCheck(req.Check).
			Save(ctx)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				fmt.Printf("rollback failed: %v\n", rollbackErr)
			}
			return status.Wrap(err, status.Internal)
		}

		// Save policy labels
		if len(req.Labels) > 0 {
			bulk := make([]*ent.PolicyLabelsCreate, len(req.Labels))
			for i, label := range req.Labels {
				bulk[i] = tx.PolicyLabels.Create().
					SetPolicyID(policy.ID).
					SetKey(label.Key).
					SetValue(label.Value)
			}

			if _, err := tx.PolicyLabels.CreateBulk(bulk...).Save(ctx); err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					fmt.Printf("rollback failed: %v\n", rollbackErr)
				}
				return status.Wrap(err, status.Internal)
			}
		}

		if err := tx.Commit(); err != nil {
			return status.Wrap(err, status.Internal)
		}

		res.ID = policy.ID
		res.Status = "created"
		return nil
	}
}

// isValidID validates the ID format:
// - All lowercase letters.
// - Multiple words separated by dashes (`-`).
// - No spaces at the beginning or end.
// - No special characters other than hyphen.
func isValidID(id string) bool {
	id = strings.TrimSpace(id)
	regex := regexp.MustCompile(`^[a-z]+(-[a-z]+)*$`)
	return regex.MatchString(id)
}
