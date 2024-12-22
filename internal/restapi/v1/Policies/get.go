package policies

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/schema"
	"github.com/swaggest/usecase/status"
)

type GetPolicyRequests struct {
	// ToDO: Fix it to UUID
	ID uuid.UUID `path:"id"`
}

type GetPolicyResponse struct {
	ID      uuid.UUID              `json:"id"`
	Name    string                 `json:"name"`
	Image   schema.Image           `json:"image"`
	Trigger map[string]interface{} `json:"trigger,omitempty"`
	Check   *schema.Check          `json:"check,omitempty"`
}

type ListPoliciesRequest struct{}

type ListPoliciesResponse struct {
	Policies []GetPolicyResponse `json:"policies"`
}

func ListPolicies(client *ent.Client) func(ctx context.Context, req ListPoliciesRequest, res *[]GetPolicyResponse) error {
	return func(ctx context.Context, req ListPoliciesRequest, res *[]GetPolicyResponse) error {
		// Query all policies
		policies, err := client.Policies.Query().All(ctx)
		if err != nil {
			return status.Wrap(errors.New("failed to fetch policies"), status.Internal)
		}

		// Build the response
		*res = make([]GetPolicyResponse, len(policies))
		for i, policy := range policies {
			(*res)[i] = GetPolicyResponse{
				ID:      policy.ID,
				Name:    policy.Name,
				Image:   policy.Image,
				Trigger: policy.Trigger,
				Check:   convertToPointer(policy.Check),
			}
		}
		return nil
	}
}

func GetPolicy(client *ent.Client) func(ctx context.Context, req GetPolicyRequests, res *GetPolicyResponse) error {
	return func(ctx context.Context, req GetPolicyRequests, res *GetPolicyResponse) error {
		// Validate ID should not be nil and should be correct UUID
		if req.ID == uuid.Nil {
			return status.Wrap(errors.New("invalid UUID: cannot be nil"), status.InvalidArgument)
		}

		// Fetch the policy by ID
		policy, err := client.Policies.Get(ctx, req.ID)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("integration not found"), status.NotFound)
			}
			return status.Wrap(fmt.Errorf("unexpected error: %v", err), status.Internal)
		}

		// Build the response
		res.ID = policy.ID
		res.Name = policy.Name
		res.Image = policy.Image
		res.Trigger = policy.Trigger
		res.Check = convertToPointer(policy.Check)

		return nil
	}
}

// Utility function to convert nullable fields to pointers
func convertToPointer(data schema.Check) *schema.Check {
	return &data
}
