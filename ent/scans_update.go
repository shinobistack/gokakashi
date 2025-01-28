// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/dialect/sql/sqljson"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/agenttasks"
	"github.com/shinobistack/gokakashi/ent/integrations"
	"github.com/shinobistack/gokakashi/ent/policies"
	"github.com/shinobistack/gokakashi/ent/predicate"
	"github.com/shinobistack/gokakashi/ent/scanlabels"
	"github.com/shinobistack/gokakashi/ent/scannotify"
	"github.com/shinobistack/gokakashi/ent/scans"
	"github.com/shinobistack/gokakashi/ent/schema"
)

// ScansUpdate is the builder for updating Scans entities.
type ScansUpdate struct {
	config
	hooks    []Hook
	mutation *ScansMutation
}

// Where appends a list predicates to the ScansUpdate builder.
func (su *ScansUpdate) Where(ps ...predicate.Scans) *ScansUpdate {
	su.mutation.Where(ps...)
	return su
}

// SetPolicyID sets the "policy_id" field.
func (su *ScansUpdate) SetPolicyID(u uuid.UUID) *ScansUpdate {
	su.mutation.SetPolicyID(u)
	return su
}

// SetNillablePolicyID sets the "policy_id" field if the given value is not nil.
func (su *ScansUpdate) SetNillablePolicyID(u *uuid.UUID) *ScansUpdate {
	if u != nil {
		su.SetPolicyID(*u)
	}
	return su
}

// SetStatus sets the "status" field.
func (su *ScansUpdate) SetStatus(s string) *ScansUpdate {
	su.mutation.SetStatus(s)
	return su
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (su *ScansUpdate) SetNillableStatus(s *string) *ScansUpdate {
	if s != nil {
		su.SetStatus(*s)
	}
	return su
}

// SetImage sets the "image" field.
func (su *ScansUpdate) SetImage(s string) *ScansUpdate {
	su.mutation.SetImage(s)
	return su
}

// SetNillableImage sets the "image" field if the given value is not nil.
func (su *ScansUpdate) SetNillableImage(s *string) *ScansUpdate {
	if s != nil {
		su.SetImage(*s)
	}
	return su
}

// SetIntegrationID sets the "integration_id" field.
func (su *ScansUpdate) SetIntegrationID(u uuid.UUID) *ScansUpdate {
	su.mutation.SetIntegrationID(u)
	return su
}

// SetNillableIntegrationID sets the "integration_id" field if the given value is not nil.
func (su *ScansUpdate) SetNillableIntegrationID(u *uuid.UUID) *ScansUpdate {
	if u != nil {
		su.SetIntegrationID(*u)
	}
	return su
}

// SetScanner sets the "scanner" field.
func (su *ScansUpdate) SetScanner(s string) *ScansUpdate {
	su.mutation.SetScanner(s)
	return su
}

// SetNillableScanner sets the "scanner" field if the given value is not nil.
func (su *ScansUpdate) SetNillableScanner(s *string) *ScansUpdate {
	if s != nil {
		su.SetScanner(*s)
	}
	return su
}

// SetNotify sets the "notify" field.
func (su *ScansUpdate) SetNotify(s []schema.Notify) *ScansUpdate {
	su.mutation.SetNotify(s)
	return su
}

// AppendNotify appends s to the "notify" field.
func (su *ScansUpdate) AppendNotify(s []schema.Notify) *ScansUpdate {
	su.mutation.AppendNotify(s)
	return su
}

// ClearNotify clears the value of the "notify" field.
func (su *ScansUpdate) ClearNotify() *ScansUpdate {
	su.mutation.ClearNotify()
	return su
}

// SetLabels sets the "labels" field.
func (su *ScansUpdate) SetLabels(sl schema.CommonLabels) *ScansUpdate {
	su.mutation.SetLabels(sl)
	return su
}

// SetNillableLabels sets the "labels" field if the given value is not nil.
func (su *ScansUpdate) SetNillableLabels(sl *schema.CommonLabels) *ScansUpdate {
	if sl != nil {
		su.SetLabels(*sl)
	}
	return su
}

// ClearLabels clears the value of the "labels" field.
func (su *ScansUpdate) ClearLabels() *ScansUpdate {
	su.mutation.ClearLabels()
	return su
}

// SetReport sets the "report" field.
func (su *ScansUpdate) SetReport(jm json.RawMessage) *ScansUpdate {
	su.mutation.SetReport(jm)
	return su
}

// AppendReport appends jm to the "report" field.
func (su *ScansUpdate) AppendReport(jm json.RawMessage) *ScansUpdate {
	su.mutation.AppendReport(jm)
	return su
}

// ClearReport clears the value of the "report" field.
func (su *ScansUpdate) ClearReport() *ScansUpdate {
	su.mutation.ClearReport()
	return su
}

// SetPolicy sets the "policy" edge to the Policies entity.
func (su *ScansUpdate) SetPolicy(p *Policies) *ScansUpdate {
	return su.SetPolicyID(p.ID)
}

// SetIntegrationsID sets the "integrations" edge to the Integrations entity by ID.
func (su *ScansUpdate) SetIntegrationsID(id uuid.UUID) *ScansUpdate {
	su.mutation.SetIntegrationsID(id)
	return su
}

// SetIntegrations sets the "integrations" edge to the Integrations entity.
func (su *ScansUpdate) SetIntegrations(i *Integrations) *ScansUpdate {
	return su.SetIntegrationsID(i.ID)
}

// AddScanLabelIDs adds the "scan_labels" edge to the ScanLabels entity by IDs.
func (su *ScansUpdate) AddScanLabelIDs(ids ...int) *ScansUpdate {
	su.mutation.AddScanLabelIDs(ids...)
	return su
}

// AddScanLabels adds the "scan_labels" edges to the ScanLabels entity.
func (su *ScansUpdate) AddScanLabels(s ...*ScanLabels) *ScansUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.AddScanLabelIDs(ids...)
}

// AddAgentTaskIDs adds the "agent_tasks" edge to the AgentTasks entity by IDs.
func (su *ScansUpdate) AddAgentTaskIDs(ids ...uuid.UUID) *ScansUpdate {
	su.mutation.AddAgentTaskIDs(ids...)
	return su
}

// AddAgentTasks adds the "agent_tasks" edges to the AgentTasks entity.
func (su *ScansUpdate) AddAgentTasks(a ...*AgentTasks) *ScansUpdate {
	ids := make([]uuid.UUID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return su.AddAgentTaskIDs(ids...)
}

// AddScanNotificationIDs adds the "scan_notifications" edge to the ScanNotify entity by IDs.
func (su *ScansUpdate) AddScanNotificationIDs(ids ...uuid.UUID) *ScansUpdate {
	su.mutation.AddScanNotificationIDs(ids...)
	return su
}

// AddScanNotifications adds the "scan_notifications" edges to the ScanNotify entity.
func (su *ScansUpdate) AddScanNotifications(s ...*ScanNotify) *ScansUpdate {
	ids := make([]uuid.UUID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.AddScanNotificationIDs(ids...)
}

// Mutation returns the ScansMutation object of the builder.
func (su *ScansUpdate) Mutation() *ScansMutation {
	return su.mutation
}

// ClearPolicy clears the "policy" edge to the Policies entity.
func (su *ScansUpdate) ClearPolicy() *ScansUpdate {
	su.mutation.ClearPolicy()
	return su
}

// ClearIntegrations clears the "integrations" edge to the Integrations entity.
func (su *ScansUpdate) ClearIntegrations() *ScansUpdate {
	su.mutation.ClearIntegrations()
	return su
}

// ClearScanLabels clears all "scan_labels" edges to the ScanLabels entity.
func (su *ScansUpdate) ClearScanLabels() *ScansUpdate {
	su.mutation.ClearScanLabels()
	return su
}

// RemoveScanLabelIDs removes the "scan_labels" edge to ScanLabels entities by IDs.
func (su *ScansUpdate) RemoveScanLabelIDs(ids ...int) *ScansUpdate {
	su.mutation.RemoveScanLabelIDs(ids...)
	return su
}

// RemoveScanLabels removes "scan_labels" edges to ScanLabels entities.
func (su *ScansUpdate) RemoveScanLabels(s ...*ScanLabels) *ScansUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.RemoveScanLabelIDs(ids...)
}

// ClearAgentTasks clears all "agent_tasks" edges to the AgentTasks entity.
func (su *ScansUpdate) ClearAgentTasks() *ScansUpdate {
	su.mutation.ClearAgentTasks()
	return su
}

// RemoveAgentTaskIDs removes the "agent_tasks" edge to AgentTasks entities by IDs.
func (su *ScansUpdate) RemoveAgentTaskIDs(ids ...uuid.UUID) *ScansUpdate {
	su.mutation.RemoveAgentTaskIDs(ids...)
	return su
}

// RemoveAgentTasks removes "agent_tasks" edges to AgentTasks entities.
func (su *ScansUpdate) RemoveAgentTasks(a ...*AgentTasks) *ScansUpdate {
	ids := make([]uuid.UUID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return su.RemoveAgentTaskIDs(ids...)
}

// ClearScanNotifications clears all "scan_notifications" edges to the ScanNotify entity.
func (su *ScansUpdate) ClearScanNotifications() *ScansUpdate {
	su.mutation.ClearScanNotifications()
	return su
}

// RemoveScanNotificationIDs removes the "scan_notifications" edge to ScanNotify entities by IDs.
func (su *ScansUpdate) RemoveScanNotificationIDs(ids ...uuid.UUID) *ScansUpdate {
	su.mutation.RemoveScanNotificationIDs(ids...)
	return su
}

// RemoveScanNotifications removes "scan_notifications" edges to ScanNotify entities.
func (su *ScansUpdate) RemoveScanNotifications(s ...*ScanNotify) *ScansUpdate {
	ids := make([]uuid.UUID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.RemoveScanNotificationIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *ScansUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, su.sqlSave, su.mutation, su.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (su *ScansUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *ScansUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *ScansUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (su *ScansUpdate) check() error {
	if v, ok := su.mutation.Status(); ok {
		if err := scans.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Scans.status": %w`, err)}
		}
	}
	if su.mutation.PolicyCleared() && len(su.mutation.PolicyIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Scans.policy"`)
	}
	if su.mutation.IntegrationsCleared() && len(su.mutation.IntegrationsIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Scans.integrations"`)
	}
	return nil
}

func (su *ScansUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := su.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(scans.Table, scans.Columns, sqlgraph.NewFieldSpec(scans.FieldID, field.TypeUUID))
	if ps := su.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := su.mutation.Status(); ok {
		_spec.SetField(scans.FieldStatus, field.TypeString, value)
	}
	if value, ok := su.mutation.Image(); ok {
		_spec.SetField(scans.FieldImage, field.TypeString, value)
	}
	if value, ok := su.mutation.Scanner(); ok {
		_spec.SetField(scans.FieldScanner, field.TypeString, value)
	}
	if value, ok := su.mutation.Notify(); ok {
		_spec.SetField(scans.FieldNotify, field.TypeJSON, value)
	}
	if value, ok := su.mutation.AppendedNotify(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, scans.FieldNotify, value)
		})
	}
	if su.mutation.NotifyCleared() {
		_spec.ClearField(scans.FieldNotify, field.TypeJSON)
	}
	if value, ok := su.mutation.Labels(); ok {
		_spec.SetField(scans.FieldLabels, field.TypeJSON, value)
	}
	if su.mutation.LabelsCleared() {
		_spec.ClearField(scans.FieldLabels, field.TypeJSON)
	}
	if value, ok := su.mutation.Report(); ok {
		_spec.SetField(scans.FieldReport, field.TypeJSON, value)
	}
	if value, ok := su.mutation.AppendedReport(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, scans.FieldReport, value)
		})
	}
	if su.mutation.ReportCleared() {
		_spec.ClearField(scans.FieldReport, field.TypeJSON)
	}
	if su.mutation.PolicyCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   scans.PolicyTable,
			Columns: []string{scans.PolicyColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(policies.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.PolicyIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   scans.PolicyTable,
			Columns: []string{scans.PolicyColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(policies.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if su.mutation.IntegrationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   scans.IntegrationsTable,
			Columns: []string{scans.IntegrationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(integrations.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.IntegrationsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   scans.IntegrationsTable,
			Columns: []string{scans.IntegrationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(integrations.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if su.mutation.ScanLabelsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   scans.ScanLabelsTable,
			Columns: []string{scans.ScanLabelsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(scanlabels.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.RemovedScanLabelsIDs(); len(nodes) > 0 && !su.mutation.ScanLabelsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   scans.ScanLabelsTable,
			Columns: []string{scans.ScanLabelsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(scanlabels.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.ScanLabelsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   scans.ScanLabelsTable,
			Columns: []string{scans.ScanLabelsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(scanlabels.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if su.mutation.AgentTasksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   scans.AgentTasksTable,
			Columns: []string{scans.AgentTasksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(agenttasks.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.RemovedAgentTasksIDs(); len(nodes) > 0 && !su.mutation.AgentTasksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   scans.AgentTasksTable,
			Columns: []string{scans.AgentTasksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(agenttasks.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.AgentTasksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   scans.AgentTasksTable,
			Columns: []string{scans.AgentTasksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(agenttasks.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if su.mutation.ScanNotificationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   scans.ScanNotificationsTable,
			Columns: []string{scans.ScanNotificationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(scannotify.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.RemovedScanNotificationsIDs(); len(nodes) > 0 && !su.mutation.ScanNotificationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   scans.ScanNotificationsTable,
			Columns: []string{scans.ScanNotificationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(scannotify.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.ScanNotificationsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   scans.ScanNotificationsTable,
			Columns: []string{scans.ScanNotificationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(scannotify.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{scans.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	su.mutation.done = true
	return n, nil
}

// ScansUpdateOne is the builder for updating a single Scans entity.
type ScansUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ScansMutation
}

// SetPolicyID sets the "policy_id" field.
func (suo *ScansUpdateOne) SetPolicyID(u uuid.UUID) *ScansUpdateOne {
	suo.mutation.SetPolicyID(u)
	return suo
}

// SetNillablePolicyID sets the "policy_id" field if the given value is not nil.
func (suo *ScansUpdateOne) SetNillablePolicyID(u *uuid.UUID) *ScansUpdateOne {
	if u != nil {
		suo.SetPolicyID(*u)
	}
	return suo
}

// SetStatus sets the "status" field.
func (suo *ScansUpdateOne) SetStatus(s string) *ScansUpdateOne {
	suo.mutation.SetStatus(s)
	return suo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (suo *ScansUpdateOne) SetNillableStatus(s *string) *ScansUpdateOne {
	if s != nil {
		suo.SetStatus(*s)
	}
	return suo
}

// SetImage sets the "image" field.
func (suo *ScansUpdateOne) SetImage(s string) *ScansUpdateOne {
	suo.mutation.SetImage(s)
	return suo
}

// SetNillableImage sets the "image" field if the given value is not nil.
func (suo *ScansUpdateOne) SetNillableImage(s *string) *ScansUpdateOne {
	if s != nil {
		suo.SetImage(*s)
	}
	return suo
}

// SetIntegrationID sets the "integration_id" field.
func (suo *ScansUpdateOne) SetIntegrationID(u uuid.UUID) *ScansUpdateOne {
	suo.mutation.SetIntegrationID(u)
	return suo
}

// SetNillableIntegrationID sets the "integration_id" field if the given value is not nil.
func (suo *ScansUpdateOne) SetNillableIntegrationID(u *uuid.UUID) *ScansUpdateOne {
	if u != nil {
		suo.SetIntegrationID(*u)
	}
	return suo
}

// SetScanner sets the "scanner" field.
func (suo *ScansUpdateOne) SetScanner(s string) *ScansUpdateOne {
	suo.mutation.SetScanner(s)
	return suo
}

// SetNillableScanner sets the "scanner" field if the given value is not nil.
func (suo *ScansUpdateOne) SetNillableScanner(s *string) *ScansUpdateOne {
	if s != nil {
		suo.SetScanner(*s)
	}
	return suo
}

// SetNotify sets the "notify" field.
func (suo *ScansUpdateOne) SetNotify(s []schema.Notify) *ScansUpdateOne {
	suo.mutation.SetNotify(s)
	return suo
}

// AppendNotify appends s to the "notify" field.
func (suo *ScansUpdateOne) AppendNotify(s []schema.Notify) *ScansUpdateOne {
	suo.mutation.AppendNotify(s)
	return suo
}

// ClearNotify clears the value of the "notify" field.
func (suo *ScansUpdateOne) ClearNotify() *ScansUpdateOne {
	suo.mutation.ClearNotify()
	return suo
}

// SetLabels sets the "labels" field.
func (suo *ScansUpdateOne) SetLabels(sl schema.CommonLabels) *ScansUpdateOne {
	suo.mutation.SetLabels(sl)
	return suo
}

// SetNillableLabels sets the "labels" field if the given value is not nil.
func (suo *ScansUpdateOne) SetNillableLabels(sl *schema.CommonLabels) *ScansUpdateOne {
	if sl != nil {
		suo.SetLabels(*sl)
	}
	return suo
}

// ClearLabels clears the value of the "labels" field.
func (suo *ScansUpdateOne) ClearLabels() *ScansUpdateOne {
	suo.mutation.ClearLabels()
	return suo
}

// SetReport sets the "report" field.
func (suo *ScansUpdateOne) SetReport(jm json.RawMessage) *ScansUpdateOne {
	suo.mutation.SetReport(jm)
	return suo
}

// AppendReport appends jm to the "report" field.
func (suo *ScansUpdateOne) AppendReport(jm json.RawMessage) *ScansUpdateOne {
	suo.mutation.AppendReport(jm)
	return suo
}

// ClearReport clears the value of the "report" field.
func (suo *ScansUpdateOne) ClearReport() *ScansUpdateOne {
	suo.mutation.ClearReport()
	return suo
}

// SetPolicy sets the "policy" edge to the Policies entity.
func (suo *ScansUpdateOne) SetPolicy(p *Policies) *ScansUpdateOne {
	return suo.SetPolicyID(p.ID)
}

// SetIntegrationsID sets the "integrations" edge to the Integrations entity by ID.
func (suo *ScansUpdateOne) SetIntegrationsID(id uuid.UUID) *ScansUpdateOne {
	suo.mutation.SetIntegrationsID(id)
	return suo
}

// SetIntegrations sets the "integrations" edge to the Integrations entity.
func (suo *ScansUpdateOne) SetIntegrations(i *Integrations) *ScansUpdateOne {
	return suo.SetIntegrationsID(i.ID)
}

// AddScanLabelIDs adds the "scan_labels" edge to the ScanLabels entity by IDs.
func (suo *ScansUpdateOne) AddScanLabelIDs(ids ...int) *ScansUpdateOne {
	suo.mutation.AddScanLabelIDs(ids...)
	return suo
}

// AddScanLabels adds the "scan_labels" edges to the ScanLabels entity.
func (suo *ScansUpdateOne) AddScanLabels(s ...*ScanLabels) *ScansUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.AddScanLabelIDs(ids...)
}

// AddAgentTaskIDs adds the "agent_tasks" edge to the AgentTasks entity by IDs.
func (suo *ScansUpdateOne) AddAgentTaskIDs(ids ...uuid.UUID) *ScansUpdateOne {
	suo.mutation.AddAgentTaskIDs(ids...)
	return suo
}

// AddAgentTasks adds the "agent_tasks" edges to the AgentTasks entity.
func (suo *ScansUpdateOne) AddAgentTasks(a ...*AgentTasks) *ScansUpdateOne {
	ids := make([]uuid.UUID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return suo.AddAgentTaskIDs(ids...)
}

// AddScanNotificationIDs adds the "scan_notifications" edge to the ScanNotify entity by IDs.
func (suo *ScansUpdateOne) AddScanNotificationIDs(ids ...uuid.UUID) *ScansUpdateOne {
	suo.mutation.AddScanNotificationIDs(ids...)
	return suo
}

// AddScanNotifications adds the "scan_notifications" edges to the ScanNotify entity.
func (suo *ScansUpdateOne) AddScanNotifications(s ...*ScanNotify) *ScansUpdateOne {
	ids := make([]uuid.UUID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.AddScanNotificationIDs(ids...)
}

// Mutation returns the ScansMutation object of the builder.
func (suo *ScansUpdateOne) Mutation() *ScansMutation {
	return suo.mutation
}

// ClearPolicy clears the "policy" edge to the Policies entity.
func (suo *ScansUpdateOne) ClearPolicy() *ScansUpdateOne {
	suo.mutation.ClearPolicy()
	return suo
}

// ClearIntegrations clears the "integrations" edge to the Integrations entity.
func (suo *ScansUpdateOne) ClearIntegrations() *ScansUpdateOne {
	suo.mutation.ClearIntegrations()
	return suo
}

// ClearScanLabels clears all "scan_labels" edges to the ScanLabels entity.
func (suo *ScansUpdateOne) ClearScanLabels() *ScansUpdateOne {
	suo.mutation.ClearScanLabels()
	return suo
}

// RemoveScanLabelIDs removes the "scan_labels" edge to ScanLabels entities by IDs.
func (suo *ScansUpdateOne) RemoveScanLabelIDs(ids ...int) *ScansUpdateOne {
	suo.mutation.RemoveScanLabelIDs(ids...)
	return suo
}

// RemoveScanLabels removes "scan_labels" edges to ScanLabels entities.
func (suo *ScansUpdateOne) RemoveScanLabels(s ...*ScanLabels) *ScansUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.RemoveScanLabelIDs(ids...)
}

// ClearAgentTasks clears all "agent_tasks" edges to the AgentTasks entity.
func (suo *ScansUpdateOne) ClearAgentTasks() *ScansUpdateOne {
	suo.mutation.ClearAgentTasks()
	return suo
}

// RemoveAgentTaskIDs removes the "agent_tasks" edge to AgentTasks entities by IDs.
func (suo *ScansUpdateOne) RemoveAgentTaskIDs(ids ...uuid.UUID) *ScansUpdateOne {
	suo.mutation.RemoveAgentTaskIDs(ids...)
	return suo
}

// RemoveAgentTasks removes "agent_tasks" edges to AgentTasks entities.
func (suo *ScansUpdateOne) RemoveAgentTasks(a ...*AgentTasks) *ScansUpdateOne {
	ids := make([]uuid.UUID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return suo.RemoveAgentTaskIDs(ids...)
}

// ClearScanNotifications clears all "scan_notifications" edges to the ScanNotify entity.
func (suo *ScansUpdateOne) ClearScanNotifications() *ScansUpdateOne {
	suo.mutation.ClearScanNotifications()
	return suo
}

// RemoveScanNotificationIDs removes the "scan_notifications" edge to ScanNotify entities by IDs.
func (suo *ScansUpdateOne) RemoveScanNotificationIDs(ids ...uuid.UUID) *ScansUpdateOne {
	suo.mutation.RemoveScanNotificationIDs(ids...)
	return suo
}

// RemoveScanNotifications removes "scan_notifications" edges to ScanNotify entities.
func (suo *ScansUpdateOne) RemoveScanNotifications(s ...*ScanNotify) *ScansUpdateOne {
	ids := make([]uuid.UUID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.RemoveScanNotificationIDs(ids...)
}

// Where appends a list predicates to the ScansUpdate builder.
func (suo *ScansUpdateOne) Where(ps ...predicate.Scans) *ScansUpdateOne {
	suo.mutation.Where(ps...)
	return suo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (suo *ScansUpdateOne) Select(field string, fields ...string) *ScansUpdateOne {
	suo.fields = append([]string{field}, fields...)
	return suo
}

// Save executes the query and returns the updated Scans entity.
func (suo *ScansUpdateOne) Save(ctx context.Context) (*Scans, error) {
	return withHooks(ctx, suo.sqlSave, suo.mutation, suo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (suo *ScansUpdateOne) SaveX(ctx context.Context) *Scans {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *ScansUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *ScansUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (suo *ScansUpdateOne) check() error {
	if v, ok := suo.mutation.Status(); ok {
		if err := scans.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Scans.status": %w`, err)}
		}
	}
	if suo.mutation.PolicyCleared() && len(suo.mutation.PolicyIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Scans.policy"`)
	}
	if suo.mutation.IntegrationsCleared() && len(suo.mutation.IntegrationsIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Scans.integrations"`)
	}
	return nil
}

func (suo *ScansUpdateOne) sqlSave(ctx context.Context) (_node *Scans, err error) {
	if err := suo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(scans.Table, scans.Columns, sqlgraph.NewFieldSpec(scans.FieldID, field.TypeUUID))
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Scans.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := suo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, scans.FieldID)
		for _, f := range fields {
			if !scans.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != scans.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := suo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := suo.mutation.Status(); ok {
		_spec.SetField(scans.FieldStatus, field.TypeString, value)
	}
	if value, ok := suo.mutation.Image(); ok {
		_spec.SetField(scans.FieldImage, field.TypeString, value)
	}
	if value, ok := suo.mutation.Scanner(); ok {
		_spec.SetField(scans.FieldScanner, field.TypeString, value)
	}
	if value, ok := suo.mutation.Notify(); ok {
		_spec.SetField(scans.FieldNotify, field.TypeJSON, value)
	}
	if value, ok := suo.mutation.AppendedNotify(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, scans.FieldNotify, value)
		})
	}
	if suo.mutation.NotifyCleared() {
		_spec.ClearField(scans.FieldNotify, field.TypeJSON)
	}
	if value, ok := suo.mutation.Labels(); ok {
		_spec.SetField(scans.FieldLabels, field.TypeJSON, value)
	}
	if suo.mutation.LabelsCleared() {
		_spec.ClearField(scans.FieldLabels, field.TypeJSON)
	}
	if value, ok := suo.mutation.Report(); ok {
		_spec.SetField(scans.FieldReport, field.TypeJSON, value)
	}
	if value, ok := suo.mutation.AppendedReport(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, scans.FieldReport, value)
		})
	}
	if suo.mutation.ReportCleared() {
		_spec.ClearField(scans.FieldReport, field.TypeJSON)
	}
	if suo.mutation.PolicyCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   scans.PolicyTable,
			Columns: []string{scans.PolicyColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(policies.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.PolicyIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   scans.PolicyTable,
			Columns: []string{scans.PolicyColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(policies.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if suo.mutation.IntegrationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   scans.IntegrationsTable,
			Columns: []string{scans.IntegrationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(integrations.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.IntegrationsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   scans.IntegrationsTable,
			Columns: []string{scans.IntegrationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(integrations.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if suo.mutation.ScanLabelsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   scans.ScanLabelsTable,
			Columns: []string{scans.ScanLabelsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(scanlabels.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.RemovedScanLabelsIDs(); len(nodes) > 0 && !suo.mutation.ScanLabelsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   scans.ScanLabelsTable,
			Columns: []string{scans.ScanLabelsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(scanlabels.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.ScanLabelsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   scans.ScanLabelsTable,
			Columns: []string{scans.ScanLabelsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(scanlabels.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if suo.mutation.AgentTasksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   scans.AgentTasksTable,
			Columns: []string{scans.AgentTasksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(agenttasks.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.RemovedAgentTasksIDs(); len(nodes) > 0 && !suo.mutation.AgentTasksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   scans.AgentTasksTable,
			Columns: []string{scans.AgentTasksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(agenttasks.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.AgentTasksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   scans.AgentTasksTable,
			Columns: []string{scans.AgentTasksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(agenttasks.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if suo.mutation.ScanNotificationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   scans.ScanNotificationsTable,
			Columns: []string{scans.ScanNotificationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(scannotify.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.RemovedScanNotificationsIDs(); len(nodes) > 0 && !suo.mutation.ScanNotificationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   scans.ScanNotificationsTable,
			Columns: []string{scans.ScanNotificationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(scannotify.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.ScanNotificationsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   scans.ScanNotificationsTable,
			Columns: []string{scans.ScanNotificationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(scannotify.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Scans{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{scans.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	suo.mutation.done = true
	return _node, nil
}
