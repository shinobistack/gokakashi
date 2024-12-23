// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/shinobistack/gokakashi/ent/policies"
	"github.com/shinobistack/gokakashi/ent/policylabels"
	"github.com/shinobistack/gokakashi/ent/predicate"
	"github.com/shinobistack/gokakashi/ent/schema"
)

// PoliciesUpdate is the builder for updating Policies entities.
type PoliciesUpdate struct {
	config
	hooks    []Hook
	mutation *PoliciesMutation
}

// Where appends a list predicates to the PoliciesUpdate builder.
func (pu *PoliciesUpdate) Where(ps ...predicate.Policies) *PoliciesUpdate {
	pu.mutation.Where(ps...)
	return pu
}

// SetName sets the "name" field.
func (pu *PoliciesUpdate) SetName(s string) *PoliciesUpdate {
	pu.mutation.SetName(s)
	return pu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (pu *PoliciesUpdate) SetNillableName(s *string) *PoliciesUpdate {
	if s != nil {
		pu.SetName(*s)
	}
	return pu
}

// SetImage sets the "image" field.
func (pu *PoliciesUpdate) SetImage(s schema.Image) *PoliciesUpdate {
	pu.mutation.SetImage(s)
	return pu
}

// SetNillableImage sets the "image" field if the given value is not nil.
func (pu *PoliciesUpdate) SetNillableImage(s *schema.Image) *PoliciesUpdate {
	if s != nil {
		pu.SetImage(*s)
	}
	return pu
}

// SetLabels sets the "labels" field.
func (pu *PoliciesUpdate) SetLabels(sl schema.PolicyLabels) *PoliciesUpdate {
	pu.mutation.SetLabels(sl)
	return pu
}

// SetNillableLabels sets the "labels" field if the given value is not nil.
func (pu *PoliciesUpdate) SetNillableLabels(sl *schema.PolicyLabels) *PoliciesUpdate {
	if sl != nil {
		pu.SetLabels(*sl)
	}
	return pu
}

// ClearLabels clears the value of the "labels" field.
func (pu *PoliciesUpdate) ClearLabels() *PoliciesUpdate {
	pu.mutation.ClearLabels()
	return pu
}

// SetTrigger sets the "trigger" field.
func (pu *PoliciesUpdate) SetTrigger(m map[string]interface{}) *PoliciesUpdate {
	pu.mutation.SetTrigger(m)
	return pu
}

// ClearTrigger clears the value of the "trigger" field.
func (pu *PoliciesUpdate) ClearTrigger() *PoliciesUpdate {
	pu.mutation.ClearTrigger()
	return pu
}

// SetCheck sets the "check" field.
func (pu *PoliciesUpdate) SetCheck(s schema.Check) *PoliciesUpdate {
	pu.mutation.SetCheck(s)
	return pu
}

// SetNillableCheck sets the "check" field if the given value is not nil.
func (pu *PoliciesUpdate) SetNillableCheck(s *schema.Check) *PoliciesUpdate {
	if s != nil {
		pu.SetCheck(*s)
	}
	return pu
}

// ClearCheck clears the value of the "check" field.
func (pu *PoliciesUpdate) ClearCheck() *PoliciesUpdate {
	pu.mutation.ClearCheck()
	return pu
}

// AddPolicyLabelIDs adds the "policy_labels" edge to the PolicyLabels entity by IDs.
func (pu *PoliciesUpdate) AddPolicyLabelIDs(ids ...int) *PoliciesUpdate {
	pu.mutation.AddPolicyLabelIDs(ids...)
	return pu
}

// AddPolicyLabels adds the "policy_labels" edges to the PolicyLabels entity.
func (pu *PoliciesUpdate) AddPolicyLabels(p ...*PolicyLabels) *PoliciesUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pu.AddPolicyLabelIDs(ids...)
}

// Mutation returns the PoliciesMutation object of the builder.
func (pu *PoliciesUpdate) Mutation() *PoliciesMutation {
	return pu.mutation
}

// ClearPolicyLabels clears all "policy_labels" edges to the PolicyLabels entity.
func (pu *PoliciesUpdate) ClearPolicyLabels() *PoliciesUpdate {
	pu.mutation.ClearPolicyLabels()
	return pu
}

// RemovePolicyLabelIDs removes the "policy_labels" edge to PolicyLabels entities by IDs.
func (pu *PoliciesUpdate) RemovePolicyLabelIDs(ids ...int) *PoliciesUpdate {
	pu.mutation.RemovePolicyLabelIDs(ids...)
	return pu
}

// RemovePolicyLabels removes "policy_labels" edges to PolicyLabels entities.
func (pu *PoliciesUpdate) RemovePolicyLabels(p ...*PolicyLabels) *PoliciesUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pu.RemovePolicyLabelIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (pu *PoliciesUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, pu.sqlSave, pu.mutation, pu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (pu *PoliciesUpdate) SaveX(ctx context.Context) int {
	affected, err := pu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pu *PoliciesUpdate) Exec(ctx context.Context) error {
	_, err := pu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pu *PoliciesUpdate) ExecX(ctx context.Context) {
	if err := pu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (pu *PoliciesUpdate) check() error {
	if v, ok := pu.mutation.Name(); ok {
		if err := policies.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Policies.name": %w`, err)}
		}
	}
	return nil
}

func (pu *PoliciesUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := pu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(policies.Table, policies.Columns, sqlgraph.NewFieldSpec(policies.FieldID, field.TypeUUID))
	if ps := pu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pu.mutation.Name(); ok {
		_spec.SetField(policies.FieldName, field.TypeString, value)
	}
	if value, ok := pu.mutation.Image(); ok {
		_spec.SetField(policies.FieldImage, field.TypeJSON, value)
	}
	if value, ok := pu.mutation.Labels(); ok {
		_spec.SetField(policies.FieldLabels, field.TypeJSON, value)
	}
	if pu.mutation.LabelsCleared() {
		_spec.ClearField(policies.FieldLabels, field.TypeJSON)
	}
	if value, ok := pu.mutation.Trigger(); ok {
		_spec.SetField(policies.FieldTrigger, field.TypeJSON, value)
	}
	if pu.mutation.TriggerCleared() {
		_spec.ClearField(policies.FieldTrigger, field.TypeJSON)
	}
	if value, ok := pu.mutation.Check(); ok {
		_spec.SetField(policies.FieldCheck, field.TypeJSON, value)
	}
	if pu.mutation.CheckCleared() {
		_spec.ClearField(policies.FieldCheck, field.TypeJSON)
	}
	if pu.mutation.PolicyLabelsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   policies.PolicyLabelsTable,
			Columns: []string{policies.PolicyLabelsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(policylabels.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.RemovedPolicyLabelsIDs(); len(nodes) > 0 && !pu.mutation.PolicyLabelsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   policies.PolicyLabelsTable,
			Columns: []string{policies.PolicyLabelsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(policylabels.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.PolicyLabelsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   policies.PolicyLabelsTable,
			Columns: []string{policies.PolicyLabelsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(policylabels.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, pu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{policies.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	pu.mutation.done = true
	return n, nil
}

// PoliciesUpdateOne is the builder for updating a single Policies entity.
type PoliciesUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *PoliciesMutation
}

// SetName sets the "name" field.
func (puo *PoliciesUpdateOne) SetName(s string) *PoliciesUpdateOne {
	puo.mutation.SetName(s)
	return puo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (puo *PoliciesUpdateOne) SetNillableName(s *string) *PoliciesUpdateOne {
	if s != nil {
		puo.SetName(*s)
	}
	return puo
}

// SetImage sets the "image" field.
func (puo *PoliciesUpdateOne) SetImage(s schema.Image) *PoliciesUpdateOne {
	puo.mutation.SetImage(s)
	return puo
}

// SetNillableImage sets the "image" field if the given value is not nil.
func (puo *PoliciesUpdateOne) SetNillableImage(s *schema.Image) *PoliciesUpdateOne {
	if s != nil {
		puo.SetImage(*s)
	}
	return puo
}

// SetLabels sets the "labels" field.
func (puo *PoliciesUpdateOne) SetLabels(sl schema.PolicyLabels) *PoliciesUpdateOne {
	puo.mutation.SetLabels(sl)
	return puo
}

// SetNillableLabels sets the "labels" field if the given value is not nil.
func (puo *PoliciesUpdateOne) SetNillableLabels(sl *schema.PolicyLabels) *PoliciesUpdateOne {
	if sl != nil {
		puo.SetLabels(*sl)
	}
	return puo
}

// ClearLabels clears the value of the "labels" field.
func (puo *PoliciesUpdateOne) ClearLabels() *PoliciesUpdateOne {
	puo.mutation.ClearLabels()
	return puo
}

// SetTrigger sets the "trigger" field.
func (puo *PoliciesUpdateOne) SetTrigger(m map[string]interface{}) *PoliciesUpdateOne {
	puo.mutation.SetTrigger(m)
	return puo
}

// ClearTrigger clears the value of the "trigger" field.
func (puo *PoliciesUpdateOne) ClearTrigger() *PoliciesUpdateOne {
	puo.mutation.ClearTrigger()
	return puo
}

// SetCheck sets the "check" field.
func (puo *PoliciesUpdateOne) SetCheck(s schema.Check) *PoliciesUpdateOne {
	puo.mutation.SetCheck(s)
	return puo
}

// SetNillableCheck sets the "check" field if the given value is not nil.
func (puo *PoliciesUpdateOne) SetNillableCheck(s *schema.Check) *PoliciesUpdateOne {
	if s != nil {
		puo.SetCheck(*s)
	}
	return puo
}

// ClearCheck clears the value of the "check" field.
func (puo *PoliciesUpdateOne) ClearCheck() *PoliciesUpdateOne {
	puo.mutation.ClearCheck()
	return puo
}

// AddPolicyLabelIDs adds the "policy_labels" edge to the PolicyLabels entity by IDs.
func (puo *PoliciesUpdateOne) AddPolicyLabelIDs(ids ...int) *PoliciesUpdateOne {
	puo.mutation.AddPolicyLabelIDs(ids...)
	return puo
}

// AddPolicyLabels adds the "policy_labels" edges to the PolicyLabels entity.
func (puo *PoliciesUpdateOne) AddPolicyLabels(p ...*PolicyLabels) *PoliciesUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return puo.AddPolicyLabelIDs(ids...)
}

// Mutation returns the PoliciesMutation object of the builder.
func (puo *PoliciesUpdateOne) Mutation() *PoliciesMutation {
	return puo.mutation
}

// ClearPolicyLabels clears all "policy_labels" edges to the PolicyLabels entity.
func (puo *PoliciesUpdateOne) ClearPolicyLabels() *PoliciesUpdateOne {
	puo.mutation.ClearPolicyLabels()
	return puo
}

// RemovePolicyLabelIDs removes the "policy_labels" edge to PolicyLabels entities by IDs.
func (puo *PoliciesUpdateOne) RemovePolicyLabelIDs(ids ...int) *PoliciesUpdateOne {
	puo.mutation.RemovePolicyLabelIDs(ids...)
	return puo
}

// RemovePolicyLabels removes "policy_labels" edges to PolicyLabels entities.
func (puo *PoliciesUpdateOne) RemovePolicyLabels(p ...*PolicyLabels) *PoliciesUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return puo.RemovePolicyLabelIDs(ids...)
}

// Where appends a list predicates to the PoliciesUpdate builder.
func (puo *PoliciesUpdateOne) Where(ps ...predicate.Policies) *PoliciesUpdateOne {
	puo.mutation.Where(ps...)
	return puo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (puo *PoliciesUpdateOne) Select(field string, fields ...string) *PoliciesUpdateOne {
	puo.fields = append([]string{field}, fields...)
	return puo
}

// Save executes the query and returns the updated Policies entity.
func (puo *PoliciesUpdateOne) Save(ctx context.Context) (*Policies, error) {
	return withHooks(ctx, puo.sqlSave, puo.mutation, puo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (puo *PoliciesUpdateOne) SaveX(ctx context.Context) *Policies {
	node, err := puo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (puo *PoliciesUpdateOne) Exec(ctx context.Context) error {
	_, err := puo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (puo *PoliciesUpdateOne) ExecX(ctx context.Context) {
	if err := puo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (puo *PoliciesUpdateOne) check() error {
	if v, ok := puo.mutation.Name(); ok {
		if err := policies.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Policies.name": %w`, err)}
		}
	}
	return nil
}

func (puo *PoliciesUpdateOne) sqlSave(ctx context.Context) (_node *Policies, err error) {
	if err := puo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(policies.Table, policies.Columns, sqlgraph.NewFieldSpec(policies.FieldID, field.TypeUUID))
	id, ok := puo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Policies.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := puo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, policies.FieldID)
		for _, f := range fields {
			if !policies.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != policies.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := puo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := puo.mutation.Name(); ok {
		_spec.SetField(policies.FieldName, field.TypeString, value)
	}
	if value, ok := puo.mutation.Image(); ok {
		_spec.SetField(policies.FieldImage, field.TypeJSON, value)
	}
	if value, ok := puo.mutation.Labels(); ok {
		_spec.SetField(policies.FieldLabels, field.TypeJSON, value)
	}
	if puo.mutation.LabelsCleared() {
		_spec.ClearField(policies.FieldLabels, field.TypeJSON)
	}
	if value, ok := puo.mutation.Trigger(); ok {
		_spec.SetField(policies.FieldTrigger, field.TypeJSON, value)
	}
	if puo.mutation.TriggerCleared() {
		_spec.ClearField(policies.FieldTrigger, field.TypeJSON)
	}
	if value, ok := puo.mutation.Check(); ok {
		_spec.SetField(policies.FieldCheck, field.TypeJSON, value)
	}
	if puo.mutation.CheckCleared() {
		_spec.ClearField(policies.FieldCheck, field.TypeJSON)
	}
	if puo.mutation.PolicyLabelsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   policies.PolicyLabelsTable,
			Columns: []string{policies.PolicyLabelsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(policylabels.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.RemovedPolicyLabelsIDs(); len(nodes) > 0 && !puo.mutation.PolicyLabelsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   policies.PolicyLabelsTable,
			Columns: []string{policies.PolicyLabelsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(policylabels.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.PolicyLabelsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   policies.PolicyLabelsTable,
			Columns: []string{policies.PolicyLabelsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(policylabels.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Policies{config: puo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, puo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{policies.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	puo.mutation.done = true
	return _node, nil
}
