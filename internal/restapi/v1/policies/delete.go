package policies

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/ent/policies"
	"github.com/shinobistack/gokakashi/ent/policylabels"
	"github.com/swaggest/usecase/status"
)

type DeletePolicyRequest struct {
	ID uuid.UUID `path:"id"`
}

type DeletePolicyResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func DeletePolicy(client *ent.Client) func(ctx context.Context, req DeletePolicyRequest, res *DeletePolicyResponse) error {
	return func(ctx context.Context, req DeletePolicyRequest, res *DeletePolicyResponse) error {
		// Validate ID
		if req.ID == uuid.Nil {
			return status.Wrap(errors.New("invalid UUID: cannot be nil"), status.InvalidArgument)
		}

		// Start a transaction
		tx, err := client.Tx(ctx)
		if err != nil {
			return status.Wrap(fmt.Errorf("failed to start transaction: %v", err), status.Internal)
		}

		// Check if the policy exists
		exists, err := tx.Policies.Query().Where(policies.ID(req.ID)).Exist(ctx)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				fmt.Printf("rollback failed: %v\n", rollbackErr)
			}
			return status.Wrap(fmt.Errorf("failed to query policy: %v", err), status.Internal)
		}
		if !exists {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				fmt.Printf("rollback failed: %v\n", rollbackErr)
			}
			return status.Wrap(errors.New("policy not found"), status.NotFound)
		}

		// Delete associated labels
		_, err = tx.PolicyLabels.Delete().
			Where(policylabels.PolicyID(req.ID)).
			Exec(ctx)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				fmt.Printf("rollback failed: %v\n", rollbackErr)
			}
			return status.Wrap(fmt.Errorf("failed to delete policy labels: %v", err), status.Internal)
		}

		// Delete the policy
		err = tx.Policies.DeleteOneID(req.ID).Exec(ctx)
		if err != nil {
			tx.Rollback()
			return status.Wrap(fmt.Errorf("failed to delete policy: %v", err), status.Internal)
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			return status.Wrap(fmt.Errorf("failed to commit transaction: %v", err), status.Internal)
		}

		// Build the response
		res.ID = req.ID
		res.Status = "deleted"
		return nil
	}
}
