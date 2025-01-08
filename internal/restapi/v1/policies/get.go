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
)

type GetPolicyRequests struct {
	// ToDO: Fix it to UUID
	ID uuid.UUID `path:"id"`
}

type GetPolicyResponse struct {
	ID      uuid.UUID              `json:"id"`
	Name    string                 `json:"name"`
	Scanner string                 `json:"scanner"`
	Image   schema.Image           `json:"image"`
	Labels  []schema.PolicyLabel   `json:"labels"`
	Trigger map[string]interface{} `json:"trigger,omitempty"`
	Notify  *[]schema.Notify       `json:"notify"`
}

type ListPoliciesRequest struct{}

type ListPoliciesResponse struct {
	Policies []GetPolicyResponse `json:"policies"`
}

func ListPolicies(client *ent.Client) func(ctx context.Context, req ListPoliciesRequest, res *[]GetPolicyResponse) error {
	return func(ctx context.Context, req ListPoliciesRequest, res *[]GetPolicyResponse) error {
		// Query all policies with their labels
		policies, err := client.Policies.Query().
			WithPolicyLabels().
			All(ctx)
		if err != nil {
			return status.Wrap(errors.New("failed to fetch policies"), status.Internal)
		}

		// Build the response
		*res = make([]GetPolicyResponse, len(policies))
		for i, policy := range policies {
			labels := make([]schema.PolicyLabel, len(policy.Edges.PolicyLabels))
			for j, label := range policy.Edges.PolicyLabels {
				labels[j] = schema.PolicyLabel{Key: label.Key, Value: label.Value}
			}

			(*res)[i] = GetPolicyResponse{
				ID:      policy.ID,
				Name:    policy.Name,
				Image:   policy.Image,
				Scanner: policy.Scanner,
				Labels:  labels,
				Trigger: policy.Trigger,
				Notify:  convertToPointer(policy.Notify),
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

		// Fetch the policy by ID with its labels
		policy, err := client.Policies.Query().
			Where(policies.ID(req.ID)).
			WithPolicyLabels(). // Include related policy labels
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("policy not found"), status.NotFound)
			}
			return status.Wrap(fmt.Errorf("unexpected error: %v", err), status.Internal)
		}

		// Map the labels
		labels := make([]schema.PolicyLabel, len(policy.Edges.PolicyLabels))
		for i, label := range policy.Edges.PolicyLabels {
			labels[i] = schema.PolicyLabel{Key: label.Key, Value: label.Value}
		}

		// Build the response
		res.ID = policy.ID
		res.Name = policy.Name
		res.Image = policy.Image
		res.Scanner = policy.Scanner
		res.Labels = labels
		res.Trigger = policy.Trigger
		res.Notify = convertToPointer(policy.Notify)

		return nil
	}
}

// Utility function to convert nullable fields to pointers
func convertToPointer(data []schema.Notify) *[]schema.Notify {
	return &data
}
