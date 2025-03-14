// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/scannotify"
	"github.com/shinobistack/gokakashi/ent/scans"
)

// ScanNotifyCreate is the builder for creating a ScanNotify entity.
type ScanNotifyCreate struct {
	config
	mutation *ScanNotifyMutation
	hooks    []Hook
}

// SetScanID sets the "scan_id" field.
func (snc *ScanNotifyCreate) SetScanID(u uuid.UUID) *ScanNotifyCreate {
	snc.mutation.SetScanID(u)
	return snc
}

// SetHash sets the "hash" field.
func (snc *ScanNotifyCreate) SetHash(s string) *ScanNotifyCreate {
	snc.mutation.SetHash(s)
	return snc
}

// SetID sets the "id" field.
func (snc *ScanNotifyCreate) SetID(u uuid.UUID) *ScanNotifyCreate {
	snc.mutation.SetID(u)
	return snc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (snc *ScanNotifyCreate) SetNillableID(u *uuid.UUID) *ScanNotifyCreate {
	if u != nil {
		snc.SetID(*u)
	}
	return snc
}

// SetScan sets the "scan" edge to the Scans entity.
func (snc *ScanNotifyCreate) SetScan(s *Scans) *ScanNotifyCreate {
	return snc.SetScanID(s.ID)
}

// Mutation returns the ScanNotifyMutation object of the builder.
func (snc *ScanNotifyCreate) Mutation() *ScanNotifyMutation {
	return snc.mutation
}

// Save creates the ScanNotify in the database.
func (snc *ScanNotifyCreate) Save(ctx context.Context) (*ScanNotify, error) {
	snc.defaults()
	return withHooks(ctx, snc.sqlSave, snc.mutation, snc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (snc *ScanNotifyCreate) SaveX(ctx context.Context) *ScanNotify {
	v, err := snc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (snc *ScanNotifyCreate) Exec(ctx context.Context) error {
	_, err := snc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (snc *ScanNotifyCreate) ExecX(ctx context.Context) {
	if err := snc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (snc *ScanNotifyCreate) defaults() {
	if _, ok := snc.mutation.ID(); !ok {
		v := scannotify.DefaultID()
		snc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (snc *ScanNotifyCreate) check() error {
	if _, ok := snc.mutation.ScanID(); !ok {
		return &ValidationError{Name: "scan_id", err: errors.New(`ent: missing required field "ScanNotify.scan_id"`)}
	}
	if _, ok := snc.mutation.Hash(); !ok {
		return &ValidationError{Name: "hash", err: errors.New(`ent: missing required field "ScanNotify.hash"`)}
	}
	if v, ok := snc.mutation.Hash(); ok {
		if err := scannotify.HashValidator(v); err != nil {
			return &ValidationError{Name: "hash", err: fmt.Errorf(`ent: validator failed for field "ScanNotify.hash": %w`, err)}
		}
	}
	if len(snc.mutation.ScanIDs()) == 0 {
		return &ValidationError{Name: "scan", err: errors.New(`ent: missing required edge "ScanNotify.scan"`)}
	}
	return nil
}

func (snc *ScanNotifyCreate) sqlSave(ctx context.Context) (*ScanNotify, error) {
	if err := snc.check(); err != nil {
		return nil, err
	}
	_node, _spec := snc.createSpec()
	if err := sqlgraph.CreateNode(ctx, snc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	snc.mutation.id = &_node.ID
	snc.mutation.done = true
	return _node, nil
}

func (snc *ScanNotifyCreate) createSpec() (*ScanNotify, *sqlgraph.CreateSpec) {
	var (
		_node = &ScanNotify{config: snc.config}
		_spec = sqlgraph.NewCreateSpec(scannotify.Table, sqlgraph.NewFieldSpec(scannotify.FieldID, field.TypeUUID))
	)
	if id, ok := snc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := snc.mutation.Hash(); ok {
		_spec.SetField(scannotify.FieldHash, field.TypeString, value)
		_node.Hash = value
	}
	if nodes := snc.mutation.ScanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   scannotify.ScanTable,
			Columns: []string{scannotify.ScanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(scans.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ScanID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// ScanNotifyCreateBulk is the builder for creating many ScanNotify entities in bulk.
type ScanNotifyCreateBulk struct {
	config
	err      error
	builders []*ScanNotifyCreate
}

// Save creates the ScanNotify entities in the database.
func (sncb *ScanNotifyCreateBulk) Save(ctx context.Context) ([]*ScanNotify, error) {
	if sncb.err != nil {
		return nil, sncb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(sncb.builders))
	nodes := make([]*ScanNotify, len(sncb.builders))
	mutators := make([]Mutator, len(sncb.builders))
	for i := range sncb.builders {
		func(i int, root context.Context) {
			builder := sncb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ScanNotifyMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, sncb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, sncb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, sncb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (sncb *ScanNotifyCreateBulk) SaveX(ctx context.Context) []*ScanNotify {
	v, err := sncb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sncb *ScanNotifyCreateBulk) Exec(ctx context.Context) error {
	_, err := sncb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sncb *ScanNotifyCreateBulk) ExecX(ctx context.Context) {
	if err := sncb.Exec(ctx); err != nil {
		panic(err)
	}
}
