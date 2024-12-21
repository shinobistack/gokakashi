package integrationtype

import (
	"context"
	"errors"
	"fmt"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/swaggest/usecase/status"
	"regexp"
	"strings"
)

type CreateIntegrationTypeRequest struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

func CreateIntegrationType(client *ent.Client) func(ctx context.Context, req CreateIntegrationTypeRequest, res *GetIntegrationTypeResponse) error {
	return func(ctx context.Context, req CreateIntegrationTypeRequest, res *GetIntegrationTypeResponse) error {
		// Validate fields
		if req.ID == "" || req.DisplayName == "" {
			return status.Wrap(errors.New("missing required fields: id and/or display_name"), status.InvalidArgument)
		}

		if !isValidID(req.ID) {
			return status.Wrap(errors.New("invalid id format; must be lowercase, alphanumeric, or dashes"), status.InvalidArgument) // 400 Bad Request
		}

		// Create integration type
		it, err := client.IntegrationType.Create().
			SetID(req.ID).
			SetDisplayName(req.DisplayName).
			Save(ctx)
		if err != nil {
			if ent.IsConstraintError(err) {
				return status.Wrap(errors.New("integration type already exists"), status.AlreadyExists)
			}
			return status.Wrap(fmt.Errorf("failed to create integration type: %v", err), status.Internal)
		}

		res.ID = it.ID
		res.DisplayName = it.DisplayName
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
