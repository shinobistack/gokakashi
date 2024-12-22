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
	Name  string       `json:"name"`
	Image schema.Image `json:"image"`
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
		// Todo: Name validation
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

		// Save policy to database
		policy, err := client.Policies.Create().
			SetName(req.Name).
			SetImage(req.Image).
			SetTrigger(req.Trigger).
			SetNillableCheck(req.Check).
			Save(ctx)
		if err != nil {
			return status.Wrap(fmt.Errorf("failed to create policy: %v", err), status.Internal)
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
