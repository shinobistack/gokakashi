// Code generated by ent, DO NOT EDIT.

package ent

import (
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/integrations"
	"github.com/shinobistack/gokakashi/ent/integrationtype"
	"github.com/shinobistack/gokakashi/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	integrationtypeFields := schema.IntegrationType{}.Fields()
	_ = integrationtypeFields
	// integrationtypeDescDisplayName is the schema descriptor for display_name field.
	integrationtypeDescDisplayName := integrationtypeFields[1].Descriptor()
	// integrationtype.DisplayNameValidator is a validator for the "display_name" field. It is called by the builders before save.
	integrationtype.DisplayNameValidator = integrationtypeDescDisplayName.Validators[0].(func(string) error)
	integrationsFields := schema.Integrations{}.Fields()
	_ = integrationsFields
	// integrationsDescName is the schema descriptor for name field.
	integrationsDescName := integrationsFields[1].Descriptor()
	// integrations.NameValidator is a validator for the "name" field. It is called by the builders before save.
	integrations.NameValidator = integrationsDescName.Validators[0].(func(string) error)
	// integrationsDescType is the schema descriptor for type field.
	integrationsDescType := integrationsFields[2].Descriptor()
	// integrations.TypeValidator is a validator for the "type" field. It is called by the builders before save.
	integrations.TypeValidator = integrationsDescType.Validators[0].(func(string) error)
	// integrationsDescID is the schema descriptor for id field.
	integrationsDescID := integrationsFields[0].Descriptor()
	// integrations.DefaultID holds the default value on creation for the id field.
	integrations.DefaultID = integrationsDescID.Default.(func() uuid.UUID)
}
