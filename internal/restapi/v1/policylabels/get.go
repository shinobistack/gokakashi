package policylabels

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/policylabels"
	"github.com/swaggest/usecase/status"
)

type ListPolicyLabelsRequest struct {
	PolicyID uuid.UUID `path:"policy_id"`
	Keys     []string  `query:"key"`
}

type PolicyLabel struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ListPolicyLabelsResponse struct {
	Labels []PolicyLabel `json:"labels"`
}

type GetPolicyLabelRequest struct {
	PolicyID uuid.UUID `path:"policy_id"`
	Key      string    `path:"key"`
}

type GetPolicyLabelResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ListPolicyLabels(client *ent.Client) func(ctx context.Context, req ListPolicyLabelsRequest, res *ListPolicyLabelsResponse) error {
	return func(ctx context.Context, req ListPolicyLabelsRequest, res *ListPolicyLabelsResponse) error {
		// Validate inputs
		if req.PolicyID == uuid.Nil {
			return status.Wrap(errors.New("invalid Policy ID"), status.InvalidArgument)
		}

		query := client.PolicyLabels.Query().Where(policylabels.PolicyID(req.PolicyID))

		// Filter by keys if provided
		if len(req.Keys) > 0 {
			query = query.Where(policylabels.KeyIn(req.Keys...))
		}

		labels, err := query.All(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		// Build the response
		res.Labels = make([]PolicyLabel, len(labels))
		for i, label := range labels {
			res.Labels[i] = PolicyLabel{Key: label.Key, Value: label.Value}
		}
		return nil
	}
}

func GetPolicyLabel(client *ent.Client) func(ctx context.Context, req GetPolicyLabelRequest, res *GetPolicyLabelResponse) error {
	return func(ctx context.Context, req GetPolicyLabelRequest, res *GetPolicyLabelResponse) error {
		// Validate inputs
		if req.PolicyID == uuid.Nil || req.Key == "" {
			return status.Wrap(errors.New("invalid Policy ID or Key"), status.InvalidArgument)
		}

		// Query the label
		label, err := client.PolicyLabels.Query().
			Where(policylabels.PolicyID(req.PolicyID), policylabels.Key(req.Key)).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Wrap(errors.New("label not found"), status.NotFound)
			}
			return status.Wrap(err, status.Internal)
		}

		res.Key = label.Key
		res.Value = label.Value
		return nil
	}
}
