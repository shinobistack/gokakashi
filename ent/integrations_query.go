// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/integrations"
	"github.com/shinobistack/gokakashi/ent/predicate"
	"github.com/shinobistack/gokakashi/ent/scans"
)

// IntegrationsQuery is the builder for querying Integrations entities.
type IntegrationsQuery struct {
	config
	ctx        *QueryContext
	order      []integrations.OrderOption
	inters     []Interceptor
	predicates []predicate.Integrations
	withScans  *ScansQuery
	withFKs    bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the IntegrationsQuery builder.
func (iq *IntegrationsQuery) Where(ps ...predicate.Integrations) *IntegrationsQuery {
	iq.predicates = append(iq.predicates, ps...)
	return iq
}

// Limit the number of records to be returned by this query.
func (iq *IntegrationsQuery) Limit(limit int) *IntegrationsQuery {
	iq.ctx.Limit = &limit
	return iq
}

// Offset to start from.
func (iq *IntegrationsQuery) Offset(offset int) *IntegrationsQuery {
	iq.ctx.Offset = &offset
	return iq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (iq *IntegrationsQuery) Unique(unique bool) *IntegrationsQuery {
	iq.ctx.Unique = &unique
	return iq
}

// Order specifies how the records should be ordered.
func (iq *IntegrationsQuery) Order(o ...integrations.OrderOption) *IntegrationsQuery {
	iq.order = append(iq.order, o...)
	return iq
}

// QueryScans chains the current query on the "scans" edge.
func (iq *IntegrationsQuery) QueryScans() *ScansQuery {
	query := (&ScansClient{config: iq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := iq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := iq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(integrations.Table, integrations.FieldID, selector),
			sqlgraph.To(scans.Table, scans.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, integrations.ScansTable, integrations.ScansColumn),
		)
		fromU = sqlgraph.SetNeighbors(iq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Integrations entity from the query.
// Returns a *NotFoundError when no Integrations was found.
func (iq *IntegrationsQuery) First(ctx context.Context) (*Integrations, error) {
	nodes, err := iq.Limit(1).All(setContextOp(ctx, iq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{integrations.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (iq *IntegrationsQuery) FirstX(ctx context.Context) *Integrations {
	node, err := iq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Integrations ID from the query.
// Returns a *NotFoundError when no Integrations ID was found.
func (iq *IntegrationsQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = iq.Limit(1).IDs(setContextOp(ctx, iq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{integrations.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (iq *IntegrationsQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := iq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Integrations entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Integrations entity is found.
// Returns a *NotFoundError when no Integrations entities are found.
func (iq *IntegrationsQuery) Only(ctx context.Context) (*Integrations, error) {
	nodes, err := iq.Limit(2).All(setContextOp(ctx, iq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{integrations.Label}
	default:
		return nil, &NotSingularError{integrations.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (iq *IntegrationsQuery) OnlyX(ctx context.Context) *Integrations {
	node, err := iq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Integrations ID in the query.
// Returns a *NotSingularError when more than one Integrations ID is found.
// Returns a *NotFoundError when no entities are found.
func (iq *IntegrationsQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = iq.Limit(2).IDs(setContextOp(ctx, iq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{integrations.Label}
	default:
		err = &NotSingularError{integrations.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (iq *IntegrationsQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := iq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of IntegrationsSlice.
func (iq *IntegrationsQuery) All(ctx context.Context) ([]*Integrations, error) {
	ctx = setContextOp(ctx, iq.ctx, ent.OpQueryAll)
	if err := iq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Integrations, *IntegrationsQuery]()
	return withInterceptors[[]*Integrations](ctx, iq, qr, iq.inters)
}

// AllX is like All, but panics if an error occurs.
func (iq *IntegrationsQuery) AllX(ctx context.Context) []*Integrations {
	nodes, err := iq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Integrations IDs.
func (iq *IntegrationsQuery) IDs(ctx context.Context) (ids []uuid.UUID, err error) {
	if iq.ctx.Unique == nil && iq.path != nil {
		iq.Unique(true)
	}
	ctx = setContextOp(ctx, iq.ctx, ent.OpQueryIDs)
	if err = iq.Select(integrations.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (iq *IntegrationsQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := iq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (iq *IntegrationsQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, iq.ctx, ent.OpQueryCount)
	if err := iq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, iq, querierCount[*IntegrationsQuery](), iq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (iq *IntegrationsQuery) CountX(ctx context.Context) int {
	count, err := iq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (iq *IntegrationsQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, iq.ctx, ent.OpQueryExist)
	switch _, err := iq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (iq *IntegrationsQuery) ExistX(ctx context.Context) bool {
	exist, err := iq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the IntegrationsQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (iq *IntegrationsQuery) Clone() *IntegrationsQuery {
	if iq == nil {
		return nil
	}
	return &IntegrationsQuery{
		config:     iq.config,
		ctx:        iq.ctx.Clone(),
		order:      append([]integrations.OrderOption{}, iq.order...),
		inters:     append([]Interceptor{}, iq.inters...),
		predicates: append([]predicate.Integrations{}, iq.predicates...),
		withScans:  iq.withScans.Clone(),
		// clone intermediate query.
		sql:  iq.sql.Clone(),
		path: iq.path,
	}
}

// WithScans tells the query-builder to eager-load the nodes that are connected to
// the "scans" edge. The optional arguments are used to configure the query builder of the edge.
func (iq *IntegrationsQuery) WithScans(opts ...func(*ScansQuery)) *IntegrationsQuery {
	query := (&ScansClient{config: iq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	iq.withScans = query
	return iq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Integrations.Query().
//		GroupBy(integrations.FieldName).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (iq *IntegrationsQuery) GroupBy(field string, fields ...string) *IntegrationsGroupBy {
	iq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &IntegrationsGroupBy{build: iq}
	grbuild.flds = &iq.ctx.Fields
	grbuild.label = integrations.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//	}
//
//	client.Integrations.Query().
//		Select(integrations.FieldName).
//		Scan(ctx, &v)
func (iq *IntegrationsQuery) Select(fields ...string) *IntegrationsSelect {
	iq.ctx.Fields = append(iq.ctx.Fields, fields...)
	sbuild := &IntegrationsSelect{IntegrationsQuery: iq}
	sbuild.label = integrations.Label
	sbuild.flds, sbuild.scan = &iq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a IntegrationsSelect configured with the given aggregations.
func (iq *IntegrationsQuery) Aggregate(fns ...AggregateFunc) *IntegrationsSelect {
	return iq.Select().Aggregate(fns...)
}

func (iq *IntegrationsQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range iq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, iq); err != nil {
				return err
			}
		}
	}
	for _, f := range iq.ctx.Fields {
		if !integrations.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if iq.path != nil {
		prev, err := iq.path(ctx)
		if err != nil {
			return err
		}
		iq.sql = prev
	}
	return nil
}

func (iq *IntegrationsQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Integrations, error) {
	var (
		nodes       = []*Integrations{}
		withFKs     = iq.withFKs
		_spec       = iq.querySpec()
		loadedTypes = [1]bool{
			iq.withScans != nil,
		}
	)
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, integrations.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Integrations).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Integrations{config: iq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, iq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := iq.withScans; query != nil {
		if err := iq.loadScans(ctx, query, nodes,
			func(n *Integrations) { n.Edges.Scans = []*Scans{} },
			func(n *Integrations, e *Scans) { n.Edges.Scans = append(n.Edges.Scans, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (iq *IntegrationsQuery) loadScans(ctx context.Context, query *ScansQuery, nodes []*Integrations, init func(*Integrations), assign func(*Integrations, *Scans)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uuid.UUID]*Integrations)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(scans.FieldIntegrationID)
	}
	query.Where(predicate.Scans(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(integrations.ScansColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.IntegrationID
		if fk == nil {
			return fmt.Errorf(`foreign-key "integration_id" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "integration_id" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (iq *IntegrationsQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := iq.querySpec()
	_spec.Node.Columns = iq.ctx.Fields
	if len(iq.ctx.Fields) > 0 {
		_spec.Unique = iq.ctx.Unique != nil && *iq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, iq.driver, _spec)
}

func (iq *IntegrationsQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(integrations.Table, integrations.Columns, sqlgraph.NewFieldSpec(integrations.FieldID, field.TypeUUID))
	_spec.From = iq.sql
	if unique := iq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if iq.path != nil {
		_spec.Unique = true
	}
	if fields := iq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, integrations.FieldID)
		for i := range fields {
			if fields[i] != integrations.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := iq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := iq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := iq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := iq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (iq *IntegrationsQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(iq.driver.Dialect())
	t1 := builder.Table(integrations.Table)
	columns := iq.ctx.Fields
	if len(columns) == 0 {
		columns = integrations.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if iq.sql != nil {
		selector = iq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if iq.ctx.Unique != nil && *iq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range iq.predicates {
		p(selector)
	}
	for _, p := range iq.order {
		p(selector)
	}
	if offset := iq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := iq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// IntegrationsGroupBy is the group-by builder for Integrations entities.
type IntegrationsGroupBy struct {
	selector
	build *IntegrationsQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (igb *IntegrationsGroupBy) Aggregate(fns ...AggregateFunc) *IntegrationsGroupBy {
	igb.fns = append(igb.fns, fns...)
	return igb
}

// Scan applies the selector query and scans the result into the given value.
func (igb *IntegrationsGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, igb.build.ctx, ent.OpQueryGroupBy)
	if err := igb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*IntegrationsQuery, *IntegrationsGroupBy](ctx, igb.build, igb, igb.build.inters, v)
}

func (igb *IntegrationsGroupBy) sqlScan(ctx context.Context, root *IntegrationsQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(igb.fns))
	for _, fn := range igb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*igb.flds)+len(igb.fns))
		for _, f := range *igb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*igb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := igb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// IntegrationsSelect is the builder for selecting fields of Integrations entities.
type IntegrationsSelect struct {
	*IntegrationsQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (is *IntegrationsSelect) Aggregate(fns ...AggregateFunc) *IntegrationsSelect {
	is.fns = append(is.fns, fns...)
	return is
}

// Scan applies the selector query and scans the result into the given value.
func (is *IntegrationsSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, is.ctx, ent.OpQuerySelect)
	if err := is.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*IntegrationsQuery, *IntegrationsSelect](ctx, is.IntegrationsQuery, is, is.inters, v)
}

func (is *IntegrationsSelect) sqlScan(ctx context.Context, root *IntegrationsQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(is.fns))
	for _, fn := range is.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*is.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := is.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
