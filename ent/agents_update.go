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
)

// AgentsUpdate is the builder for updating Agents entities.
type AgentsUpdate struct {
	config
	hooks    []Hook
	mutation *AgentsMutation
}

// Where appends a list predicates to the AgentsUpdate builder.
func (au *AgentsUpdate) Where(ps ...predicate.Agents) *AgentsUpdate {
	au.mutation.Where(ps...)
	return au
}

// SetStatus sets the "status" field.
func (au *AgentsUpdate) SetStatus(s string) *AgentsUpdate {
	au.mutation.SetStatus(s)
	return au
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (au *AgentsUpdate) SetNillableStatus(s *string) *AgentsUpdate {
	if s != nil {
		au.SetStatus(*s)
	}
	return au
}

// AddAgentTaskIDs adds the "agent_tasks" edge to the AgentTasks entity by IDs.
func (au *AgentsUpdate) AddAgentTaskIDs(ids ...uuid.UUID) *AgentsUpdate {
	au.mutation.AddAgentTaskIDs(ids...)
	return au
}

// AddAgentTasks adds the "agent_tasks" edges to the AgentTasks entity.
func (au *AgentsUpdate) AddAgentTasks(a ...*AgentTasks) *AgentsUpdate {
	ids := make([]uuid.UUID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return au.AddAgentTaskIDs(ids...)
}

// Mutation returns the AgentsMutation object of the builder.
func (au *AgentsUpdate) Mutation() *AgentsMutation {
	return au.mutation
}

// ClearAgentTasks clears all "agent_tasks" edges to the AgentTasks entity.
func (au *AgentsUpdate) ClearAgentTasks() *AgentsUpdate {
	au.mutation.ClearAgentTasks()
	return au
}

// RemoveAgentTaskIDs removes the "agent_tasks" edge to AgentTasks entities by IDs.
func (au *AgentsUpdate) RemoveAgentTaskIDs(ids ...uuid.UUID) *AgentsUpdate {
	au.mutation.RemoveAgentTaskIDs(ids...)
	return au
}

// RemoveAgentTasks removes "agent_tasks" edges to AgentTasks entities.
func (au *AgentsUpdate) RemoveAgentTasks(a ...*AgentTasks) *AgentsUpdate {
	ids := make([]uuid.UUID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return au.RemoveAgentTaskIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (au *AgentsUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, au.sqlSave, au.mutation, au.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (au *AgentsUpdate) SaveX(ctx context.Context) int {
	affected, err := au.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (au *AgentsUpdate) Exec(ctx context.Context) error {
	_, err := au.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (au *AgentsUpdate) ExecX(ctx context.Context) {
	if err := au.Exec(ctx); err != nil {
		panic(err)
	}
}

func (au *AgentsUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(agents.Table, agents.Columns, sqlgraph.NewFieldSpec(agents.FieldID, field.TypeInt))
	if ps := au.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := au.mutation.Status(); ok {
		_spec.SetField(agents.FieldStatus, field.TypeString, value)
	}
	if au.mutation.AgentTasksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   agents.AgentTasksTable,
			Columns: []string{agents.AgentTasksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(agenttasks.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := au.mutation.RemovedAgentTasksIDs(); len(nodes) > 0 && !au.mutation.AgentTasksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   agents.AgentTasksTable,
			Columns: []string{agents.AgentTasksColumn},
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
	if nodes := au.mutation.AgentTasksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   agents.AgentTasksTable,
			Columns: []string{agents.AgentTasksColumn},
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
	if n, err = sqlgraph.UpdateNodes(ctx, au.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{agents.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	au.mutation.done = true
	return n, nil
}

// AgentsUpdateOne is the builder for updating a single Agents entity.
type AgentsUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *AgentsMutation
}

// SetStatus sets the "status" field.
func (auo *AgentsUpdateOne) SetStatus(s string) *AgentsUpdateOne {
	auo.mutation.SetStatus(s)
	return auo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (auo *AgentsUpdateOne) SetNillableStatus(s *string) *AgentsUpdateOne {
	if s != nil {
		auo.SetStatus(*s)
	}
	return auo
}

// AddAgentTaskIDs adds the "agent_tasks" edge to the AgentTasks entity by IDs.
func (auo *AgentsUpdateOne) AddAgentTaskIDs(ids ...uuid.UUID) *AgentsUpdateOne {
	auo.mutation.AddAgentTaskIDs(ids...)
	return auo
}

// AddAgentTasks adds the "agent_tasks" edges to the AgentTasks entity.
func (auo *AgentsUpdateOne) AddAgentTasks(a ...*AgentTasks) *AgentsUpdateOne {
	ids := make([]uuid.UUID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return auo.AddAgentTaskIDs(ids...)
}

// Mutation returns the AgentsMutation object of the builder.
func (auo *AgentsUpdateOne) Mutation() *AgentsMutation {
	return auo.mutation
}

// ClearAgentTasks clears all "agent_tasks" edges to the AgentTasks entity.
func (auo *AgentsUpdateOne) ClearAgentTasks() *AgentsUpdateOne {
	auo.mutation.ClearAgentTasks()
	return auo
}

// RemoveAgentTaskIDs removes the "agent_tasks" edge to AgentTasks entities by IDs.
func (auo *AgentsUpdateOne) RemoveAgentTaskIDs(ids ...uuid.UUID) *AgentsUpdateOne {
	auo.mutation.RemoveAgentTaskIDs(ids...)
	return auo
}

// RemoveAgentTasks removes "agent_tasks" edges to AgentTasks entities.
func (auo *AgentsUpdateOne) RemoveAgentTasks(a ...*AgentTasks) *AgentsUpdateOne {
	ids := make([]uuid.UUID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return auo.RemoveAgentTaskIDs(ids...)
}

// Where appends a list predicates to the AgentsUpdate builder.
func (auo *AgentsUpdateOne) Where(ps ...predicate.Agents) *AgentsUpdateOne {
	auo.mutation.Where(ps...)
	return auo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (auo *AgentsUpdateOne) Select(field string, fields ...string) *AgentsUpdateOne {
	auo.fields = append([]string{field}, fields...)
	return auo
}

// Save executes the query and returns the updated Agents entity.
func (auo *AgentsUpdateOne) Save(ctx context.Context) (*Agents, error) {
	return withHooks(ctx, auo.sqlSave, auo.mutation, auo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (auo *AgentsUpdateOne) SaveX(ctx context.Context) *Agents {
	node, err := auo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (auo *AgentsUpdateOne) Exec(ctx context.Context) error {
	_, err := auo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (auo *AgentsUpdateOne) ExecX(ctx context.Context) {
	if err := auo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (auo *AgentsUpdateOne) sqlSave(ctx context.Context) (_node *Agents, err error) {
	_spec := sqlgraph.NewUpdateSpec(agents.Table, agents.Columns, sqlgraph.NewFieldSpec(agents.FieldID, field.TypeInt))
	id, ok := auo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Agents.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := auo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, agents.FieldID)
		for _, f := range fields {
			if !agents.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != agents.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := auo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := auo.mutation.Status(); ok {
		_spec.SetField(agents.FieldStatus, field.TypeString, value)
	}
	if auo.mutation.AgentTasksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   agents.AgentTasksTable,
			Columns: []string{agents.AgentTasksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(agenttasks.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := auo.mutation.RemovedAgentTasksIDs(); len(nodes) > 0 && !auo.mutation.AgentTasksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   agents.AgentTasksTable,
			Columns: []string{agents.AgentTasksColumn},
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
	if nodes := auo.mutation.AgentTasksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   agents.AgentTasksTable,
			Columns: []string{agents.AgentTasksColumn},
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
	_node = &Agents{config: auo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, auo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{agents.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	auo.mutation.done = true
	return _node, nil
}
