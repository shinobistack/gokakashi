// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/agentlabels"
	"github.com/shinobistack/gokakashi/ent/agents"
	"github.com/shinobistack/gokakashi/ent/agenttasks"
	"github.com/shinobistack/gokakashi/ent/integrations"
	"github.com/shinobistack/gokakashi/ent/integrationtype"
	"github.com/shinobistack/gokakashi/ent/policies"
	"github.com/shinobistack/gokakashi/ent/policylabels"
	"github.com/shinobistack/gokakashi/ent/scanlabels"
	"github.com/shinobistack/gokakashi/ent/scannotify"
	"github.com/shinobistack/gokakashi/ent/scans"
	"github.com/shinobistack/gokakashi/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	agentlabelsFields := schema.AgentLabels{}.Fields()
	_ = agentlabelsFields
	// agentlabelsDescKey is the schema descriptor for key field.
	agentlabelsDescKey := agentlabelsFields[1].Descriptor()
	// agentlabels.KeyValidator is a validator for the "key" field. It is called by the builders before save.
	agentlabels.KeyValidator = agentlabelsDescKey.Validators[0].(func(string) error)
	// agentlabelsDescValue is the schema descriptor for value field.
	agentlabelsDescValue := agentlabelsFields[2].Descriptor()
	// agentlabels.ValueValidator is a validator for the "value" field. It is called by the builders before save.
	agentlabels.ValueValidator = agentlabelsDescValue.Validators[0].(func(string) error)
	agenttasksFields := schema.AgentTasks{}.Fields()
	_ = agenttasksFields
	// agenttasksDescStatus is the schema descriptor for status field.
	agenttasksDescStatus := agenttasksFields[3].Descriptor()
	// agenttasks.DefaultStatus holds the default value on creation for the status field.
	agenttasks.DefaultStatus = agenttasksDescStatus.Default.(string)
	// agenttasks.StatusValidator is a validator for the "status" field. It is called by the builders before save.
	agenttasks.StatusValidator = agenttasksDescStatus.Validators[0].(func(string) error)
	// agenttasksDescCreatedAt is the schema descriptor for created_at field.
	agenttasksDescCreatedAt := agenttasksFields[4].Descriptor()
	// agenttasks.DefaultCreatedAt holds the default value on creation for the created_at field.
	agenttasks.DefaultCreatedAt = agenttasksDescCreatedAt.Default.(func() time.Time)
	// agenttasksDescID is the schema descriptor for id field.
	agenttasksDescID := agenttasksFields[0].Descriptor()
	// agenttasks.DefaultID holds the default value on creation for the id field.
	agenttasks.DefaultID = agenttasksDescID.Default.(func() uuid.UUID)
	agentsFields := schema.Agents{}.Fields()
	_ = agentsFields
	// agentsDescStatus is the schema descriptor for status field.
	agentsDescStatus := agentsFields[2].Descriptor()
	// agents.DefaultStatus holds the default value on creation for the status field.
	agents.DefaultStatus = agentsDescStatus.Default.(string)
	// agents.StatusValidator is a validator for the "status" field. It is called by the builders before save.
	agents.StatusValidator = agentsDescStatus.Validators[0].(func(string) error)
	// agentsDescLastSeen is the schema descriptor for last_seen field.
	agentsDescLastSeen := agentsFields[6].Descriptor()
	// agents.DefaultLastSeen holds the default value on creation for the last_seen field.
	agents.DefaultLastSeen = agentsDescLastSeen.Default.(func() time.Time)
	// agents.UpdateDefaultLastSeen holds the default value on update for the last_seen field.
	agents.UpdateDefaultLastSeen = agentsDescLastSeen.UpdateDefault.(func() time.Time)
	// agentsDescLastHeartbeat is the schema descriptor for last_heartbeat field.
	agentsDescLastHeartbeat := agentsFields[7].Descriptor()
	// agents.DefaultLastHeartbeat holds the default value on creation for the last_heartbeat field.
	agents.DefaultLastHeartbeat = agentsDescLastHeartbeat.Default.(func() time.Time)
	// agents.UpdateDefaultLastHeartbeat holds the default value on update for the last_heartbeat field.
	agents.UpdateDefaultLastHeartbeat = agentsDescLastHeartbeat.UpdateDefault.(func() time.Time)
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
	policiesFields := schema.Policies{}.Fields()
	_ = policiesFields
	// policiesDescName is the schema descriptor for name field.
	policiesDescName := policiesFields[1].Descriptor()
	// policies.NameValidator is a validator for the "name" field. It is called by the builders before save.
	policies.NameValidator = func() func(string) error {
		validators := policiesDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// policiesDescID is the schema descriptor for id field.
	policiesDescID := policiesFields[0].Descriptor()
	// policies.DefaultID holds the default value on creation for the id field.
	policies.DefaultID = policiesDescID.Default.(func() uuid.UUID)
	policylabelsFields := schema.PolicyLabels{}.Fields()
	_ = policylabelsFields
	// policylabelsDescKey is the schema descriptor for key field.
	policylabelsDescKey := policylabelsFields[1].Descriptor()
	// policylabels.KeyValidator is a validator for the "key" field. It is called by the builders before save.
	policylabels.KeyValidator = policylabelsDescKey.Validators[0].(func(string) error)
	// policylabelsDescValue is the schema descriptor for value field.
	policylabelsDescValue := policylabelsFields[2].Descriptor()
	// policylabels.ValueValidator is a validator for the "value" field. It is called by the builders before save.
	policylabels.ValueValidator = policylabelsDescValue.Validators[0].(func(string) error)
	scanlabelsFields := schema.ScanLabels{}.Fields()
	_ = scanlabelsFields
	// scanlabelsDescKey is the schema descriptor for key field.
	scanlabelsDescKey := scanlabelsFields[1].Descriptor()
	// scanlabels.KeyValidator is a validator for the "key" field. It is called by the builders before save.
	scanlabels.KeyValidator = scanlabelsDescKey.Validators[0].(func(string) error)
	// scanlabelsDescValue is the schema descriptor for value field.
	scanlabelsDescValue := scanlabelsFields[2].Descriptor()
	// scanlabels.ValueValidator is a validator for the "value" field. It is called by the builders before save.
	scanlabels.ValueValidator = scanlabelsDescValue.Validators[0].(func(string) error)
	scannotifyFields := schema.ScanNotify{}.Fields()
	_ = scannotifyFields
	// scannotifyDescHash is the schema descriptor for hash field.
	scannotifyDescHash := scannotifyFields[2].Descriptor()
	// scannotify.HashValidator is a validator for the "hash" field. It is called by the builders before save.
	scannotify.HashValidator = scannotifyDescHash.Validators[0].(func(string) error)
	// scannotifyDescID is the schema descriptor for id field.
	scannotifyDescID := scannotifyFields[0].Descriptor()
	// scannotify.DefaultID holds the default value on creation for the id field.
	scannotify.DefaultID = scannotifyDescID.Default.(func() uuid.UUID)
	scansFields := schema.Scans{}.Fields()
	_ = scansFields
	// scansDescStatus is the schema descriptor for status field.
	scansDescStatus := scansFields[2].Descriptor()
	// scans.DefaultStatus holds the default value on creation for the status field.
	scans.DefaultStatus = scansDescStatus.Default.(string)
	// scans.StatusValidator is a validator for the "status" field. It is called by the builders before save.
	scans.StatusValidator = scansDescStatus.Validators[0].(func(string) error)
	// scansDescID is the schema descriptor for id field.
	scansDescID := scansFields[0].Descriptor()
	// scans.DefaultID holds the default value on creation for the id field.
	scans.DefaultID = scansDescID.Default.(func() uuid.UUID)
}
