// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/shinobistack/gokakashi/ent/agentlabels"
	"github.com/shinobistack/gokakashi/ent/predicate"
)

// AgentLabelsDelete is the builder for deleting a AgentLabels entity.
type AgentLabelsDelete struct {
	config
	hooks    []Hook
	mutation *AgentLabelsMutation
}

// Where appends a list predicates to the AgentLabelsDelete builder.
func (ald *AgentLabelsDelete) Where(ps ...predicate.AgentLabels) *AgentLabelsDelete {
	ald.mutation.Where(ps...)
	return ald
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ald *AgentLabelsDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, ald.sqlExec, ald.mutation, ald.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (ald *AgentLabelsDelete) ExecX(ctx context.Context) int {
	n, err := ald.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ald *AgentLabelsDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(agentlabels.Table, sqlgraph.NewFieldSpec(agentlabels.FieldID, field.TypeInt))
	if ps := ald.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, ald.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	ald.mutation.done = true
	return affected, err
}

// AgentLabelsDeleteOne is the builder for deleting a single AgentLabels entity.
type AgentLabelsDeleteOne struct {
	ald *AgentLabelsDelete
}

// Where appends a list predicates to the AgentLabelsDelete builder.
func (aldo *AgentLabelsDeleteOne) Where(ps ...predicate.AgentLabels) *AgentLabelsDeleteOne {
	aldo.ald.mutation.Where(ps...)
	return aldo
}

// Exec executes the deletion query.
func (aldo *AgentLabelsDeleteOne) Exec(ctx context.Context) error {
	n, err := aldo.ald.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{agentlabels.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (aldo *AgentLabelsDeleteOne) ExecX(ctx context.Context) {
	if err := aldo.Exec(ctx); err != nil {
		panic(err)
	}
}
