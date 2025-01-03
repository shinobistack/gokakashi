// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/agents"
	"github.com/shinobistack/gokakashi/ent/agenttasks"
	"github.com/shinobistack/gokakashi/ent/predicate"
	"github.com/shinobistack/gokakashi/ent/scans"
)

// AgentTasksUpdate is the builder for updating AgentTasks entities.
type AgentTasksUpdate struct {
	config
	hooks    []Hook
	mutation *AgentTasksMutation
}

// Where appends a list predicates to the AgentTasksUpdate builder.
func (atu *AgentTasksUpdate) Where(ps ...predicate.AgentTasks) *AgentTasksUpdate {
	atu.mutation.Where(ps...)
	return atu
}

// SetAgentID sets the "agent_id" field.
func (atu *AgentTasksUpdate) SetAgentID(i int) *AgentTasksUpdate {
	atu.mutation.SetAgentID(i)
	return atu
}

// SetNillableAgentID sets the "agent_id" field if the given value is not nil.
func (atu *AgentTasksUpdate) SetNillableAgentID(i *int) *AgentTasksUpdate {
	if i != nil {
		atu.SetAgentID(*i)
	}
	return atu
}

// SetScanID sets the "scan_id" field.
func (atu *AgentTasksUpdate) SetScanID(u uuid.UUID) *AgentTasksUpdate {
	atu.mutation.SetScanID(u)
	return atu
}

// SetNillableScanID sets the "scan_id" field if the given value is not nil.
func (atu *AgentTasksUpdate) SetNillableScanID(u *uuid.UUID) *AgentTasksUpdate {
	if u != nil {
		atu.SetScanID(*u)
	}
	return atu
}

// SetStatus sets the "status" field.
func (atu *AgentTasksUpdate) SetStatus(s string) *AgentTasksUpdate {
	atu.mutation.SetStatus(s)
	return atu
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (atu *AgentTasksUpdate) SetNillableStatus(s *string) *AgentTasksUpdate {
	if s != nil {
		atu.SetStatus(*s)
	}
	return atu
}

// SetAgent sets the "agent" edge to the Agents entity.
func (atu *AgentTasksUpdate) SetAgent(a *Agents) *AgentTasksUpdate {
	return atu.SetAgentID(a.ID)
}

// SetScan sets the "scan" edge to the Scans entity.
func (atu *AgentTasksUpdate) SetScan(s *Scans) *AgentTasksUpdate {
	return atu.SetScanID(s.ID)
}

// Mutation returns the AgentTasksMutation object of the builder.
func (atu *AgentTasksUpdate) Mutation() *AgentTasksMutation {
	return atu.mutation
}

// ClearAgent clears the "agent" edge to the Agents entity.
func (atu *AgentTasksUpdate) ClearAgent() *AgentTasksUpdate {
	atu.mutation.ClearAgent()
	return atu
}

// ClearScan clears the "scan" edge to the Scans entity.
func (atu *AgentTasksUpdate) ClearScan() *AgentTasksUpdate {
	atu.mutation.ClearScan()
	return atu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (atu *AgentTasksUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, atu.sqlSave, atu.mutation, atu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (atu *AgentTasksUpdate) SaveX(ctx context.Context) int {
	affected, err := atu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (atu *AgentTasksUpdate) Exec(ctx context.Context) error {
	_, err := atu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (atu *AgentTasksUpdate) ExecX(ctx context.Context) {
	if err := atu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (atu *AgentTasksUpdate) check() error {
	if v, ok := atu.mutation.Status(); ok {
		if err := agenttasks.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "AgentTasks.status": %w`, err)}
		}
	}
	if atu.mutation.AgentCleared() && len(atu.mutation.AgentIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "AgentTasks.agent"`)
	}
	if atu.mutation.ScanCleared() && len(atu.mutation.ScanIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "AgentTasks.scan"`)
	}
	return nil
}

func (atu *AgentTasksUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := atu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(agenttasks.Table, agenttasks.Columns, sqlgraph.NewFieldSpec(agenttasks.FieldID, field.TypeUUID))
	if ps := atu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := atu.mutation.Status(); ok {
		_spec.SetField(agenttasks.FieldStatus, field.TypeString, value)
	}
	if atu.mutation.AgentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   agenttasks.AgentTable,
			Columns: []string{agenttasks.AgentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(agents.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := atu.mutation.AgentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   agenttasks.AgentTable,
			Columns: []string{agenttasks.AgentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(agents.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if atu.mutation.ScanCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   agenttasks.ScanTable,
			Columns: []string{agenttasks.ScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(scans.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := atu.mutation.ScanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   agenttasks.ScanTable,
			Columns: []string{agenttasks.ScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(scans.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, atu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{agenttasks.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	atu.mutation.done = true
	return n, nil
}

// AgentTasksUpdateOne is the builder for updating a single AgentTasks entity.
type AgentTasksUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *AgentTasksMutation
}

// SetAgentID sets the "agent_id" field.
func (atuo *AgentTasksUpdateOne) SetAgentID(i int) *AgentTasksUpdateOne {
	atuo.mutation.SetAgentID(i)
	return atuo
}

// SetNillableAgentID sets the "agent_id" field if the given value is not nil.
func (atuo *AgentTasksUpdateOne) SetNillableAgentID(i *int) *AgentTasksUpdateOne {
	if i != nil {
		atuo.SetAgentID(*i)
	}
	return atuo
}

// SetScanID sets the "scan_id" field.
func (atuo *AgentTasksUpdateOne) SetScanID(u uuid.UUID) *AgentTasksUpdateOne {
	atuo.mutation.SetScanID(u)
	return atuo
}

// SetNillableScanID sets the "scan_id" field if the given value is not nil.
func (atuo *AgentTasksUpdateOne) SetNillableScanID(u *uuid.UUID) *AgentTasksUpdateOne {
	if u != nil {
		atuo.SetScanID(*u)
	}
	return atuo
}

// SetStatus sets the "status" field.
func (atuo *AgentTasksUpdateOne) SetStatus(s string) *AgentTasksUpdateOne {
	atuo.mutation.SetStatus(s)
	return atuo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (atuo *AgentTasksUpdateOne) SetNillableStatus(s *string) *AgentTasksUpdateOne {
	if s != nil {
		atuo.SetStatus(*s)
	}
	return atuo
}

// SetAgent sets the "agent" edge to the Agents entity.
func (atuo *AgentTasksUpdateOne) SetAgent(a *Agents) *AgentTasksUpdateOne {
	return atuo.SetAgentID(a.ID)
}

// SetScan sets the "scan" edge to the Scans entity.
func (atuo *AgentTasksUpdateOne) SetScan(s *Scans) *AgentTasksUpdateOne {
	return atuo.SetScanID(s.ID)
}

// Mutation returns the AgentTasksMutation object of the builder.
func (atuo *AgentTasksUpdateOne) Mutation() *AgentTasksMutation {
	return atuo.mutation
}

// ClearAgent clears the "agent" edge to the Agents entity.
func (atuo *AgentTasksUpdateOne) ClearAgent() *AgentTasksUpdateOne {
	atuo.mutation.ClearAgent()
	return atuo
}

// ClearScan clears the "scan" edge to the Scans entity.
func (atuo *AgentTasksUpdateOne) ClearScan() *AgentTasksUpdateOne {
	atuo.mutation.ClearScan()
	return atuo
}

// Where appends a list predicates to the AgentTasksUpdate builder.
func (atuo *AgentTasksUpdateOne) Where(ps ...predicate.AgentTasks) *AgentTasksUpdateOne {
	atuo.mutation.Where(ps...)
	return atuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (atuo *AgentTasksUpdateOne) Select(field string, fields ...string) *AgentTasksUpdateOne {
	atuo.fields = append([]string{field}, fields...)
	return atuo
}

// Save executes the query and returns the updated AgentTasks entity.
func (atuo *AgentTasksUpdateOne) Save(ctx context.Context) (*AgentTasks, error) {
	return withHooks(ctx, atuo.sqlSave, atuo.mutation, atuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (atuo *AgentTasksUpdateOne) SaveX(ctx context.Context) *AgentTasks {
	node, err := atuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (atuo *AgentTasksUpdateOne) Exec(ctx context.Context) error {
	_, err := atuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (atuo *AgentTasksUpdateOne) ExecX(ctx context.Context) {
	if err := atuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (atuo *AgentTasksUpdateOne) check() error {
	if v, ok := atuo.mutation.Status(); ok {
		if err := agenttasks.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "AgentTasks.status": %w`, err)}
		}
	}
	if atuo.mutation.AgentCleared() && len(atuo.mutation.AgentIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "AgentTasks.agent"`)
	}
	if atuo.mutation.ScanCleared() && len(atuo.mutation.ScanIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "AgentTasks.scan"`)
	}
	return nil
}

func (atuo *AgentTasksUpdateOne) sqlSave(ctx context.Context) (_node *AgentTasks, err error) {
	if err := atuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(agenttasks.Table, agenttasks.Columns, sqlgraph.NewFieldSpec(agenttasks.FieldID, field.TypeUUID))
	id, ok := atuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "AgentTasks.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := atuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, agenttasks.FieldID)
		for _, f := range fields {
			if !agenttasks.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != agenttasks.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := atuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := atuo.mutation.Status(); ok {
		_spec.SetField(agenttasks.FieldStatus, field.TypeString, value)
	}
	if atuo.mutation.AgentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   agenttasks.AgentTable,
			Columns: []string{agenttasks.AgentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(agents.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := atuo.mutation.AgentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   agenttasks.AgentTable,
			Columns: []string{agenttasks.AgentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(agents.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if atuo.mutation.ScanCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   agenttasks.ScanTable,
			Columns: []string{agenttasks.ScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(scans.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := atuo.mutation.ScanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   agenttasks.ScanTable,
			Columns: []string{agenttasks.ScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(scans.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &AgentTasks{config: atuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, atuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{agenttasks.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	atuo.mutation.done = true
	return _node, nil
}
