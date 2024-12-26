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
	"github.com/shinobistack/gokakashi/ent/agents"
	"github.com/shinobistack/gokakashi/ent/agenttasks"
	"github.com/shinobistack/gokakashi/ent/predicate"
)

// AgentsQuery is the builder for querying Agents entities.
type AgentsQuery struct {
	config
	ctx            *QueryContext
	order          []agents.OrderOption
	inters         []Interceptor
	predicates     []predicate.Agents
	withAgentTasks *AgentTasksQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the AgentsQuery builder.
func (aq *AgentsQuery) Where(ps ...predicate.Agents) *AgentsQuery {
	aq.predicates = append(aq.predicates, ps...)
	return aq
}

// Limit the number of records to be returned by this query.
func (aq *AgentsQuery) Limit(limit int) *AgentsQuery {
	aq.ctx.Limit = &limit
	return aq
}

// Offset to start from.
func (aq *AgentsQuery) Offset(offset int) *AgentsQuery {
	aq.ctx.Offset = &offset
	return aq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (aq *AgentsQuery) Unique(unique bool) *AgentsQuery {
	aq.ctx.Unique = &unique
	return aq
}

// Order specifies how the records should be ordered.
func (aq *AgentsQuery) Order(o ...agents.OrderOption) *AgentsQuery {
	aq.order = append(aq.order, o...)
	return aq
}

// QueryAgentTasks chains the current query on the "agent_tasks" edge.
func (aq *AgentsQuery) QueryAgentTasks() *AgentTasksQuery {
	query := (&AgentTasksClient{config: aq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := aq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := aq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(agents.Table, agents.FieldID, selector),
			sqlgraph.To(agenttasks.Table, agenttasks.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, agents.AgentTasksTable, agents.AgentTasksColumn),
		)
		fromU = sqlgraph.SetNeighbors(aq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Agents entity from the query.
// Returns a *NotFoundError when no Agents was found.
func (aq *AgentsQuery) First(ctx context.Context) (*Agents, error) {
	nodes, err := aq.Limit(1).All(setContextOp(ctx, aq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{agents.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (aq *AgentsQuery) FirstX(ctx context.Context) *Agents {
	node, err := aq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Agents ID from the query.
// Returns a *NotFoundError when no Agents ID was found.
func (aq *AgentsQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = aq.Limit(1).IDs(setContextOp(ctx, aq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{agents.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (aq *AgentsQuery) FirstIDX(ctx context.Context) int {
	id, err := aq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Agents entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Agents entity is found.
// Returns a *NotFoundError when no Agents entities are found.
func (aq *AgentsQuery) Only(ctx context.Context) (*Agents, error) {
	nodes, err := aq.Limit(2).All(setContextOp(ctx, aq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{agents.Label}
	default:
		return nil, &NotSingularError{agents.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (aq *AgentsQuery) OnlyX(ctx context.Context) *Agents {
	node, err := aq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Agents ID in the query.
// Returns a *NotSingularError when more than one Agents ID is found.
// Returns a *NotFoundError when no entities are found.
func (aq *AgentsQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = aq.Limit(2).IDs(setContextOp(ctx, aq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{agents.Label}
	default:
		err = &NotSingularError{agents.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (aq *AgentsQuery) OnlyIDX(ctx context.Context) int {
	id, err := aq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of AgentsSlice.
func (aq *AgentsQuery) All(ctx context.Context) ([]*Agents, error) {
	ctx = setContextOp(ctx, aq.ctx, ent.OpQueryAll)
	if err := aq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Agents, *AgentsQuery]()
	return withInterceptors[[]*Agents](ctx, aq, qr, aq.inters)
}

// AllX is like All, but panics if an error occurs.
func (aq *AgentsQuery) AllX(ctx context.Context) []*Agents {
	nodes, err := aq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Agents IDs.
func (aq *AgentsQuery) IDs(ctx context.Context) (ids []int, err error) {
	if aq.ctx.Unique == nil && aq.path != nil {
		aq.Unique(true)
	}
	ctx = setContextOp(ctx, aq.ctx, ent.OpQueryIDs)
	if err = aq.Select(agents.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (aq *AgentsQuery) IDsX(ctx context.Context) []int {
	ids, err := aq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (aq *AgentsQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, aq.ctx, ent.OpQueryCount)
	if err := aq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, aq, querierCount[*AgentsQuery](), aq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (aq *AgentsQuery) CountX(ctx context.Context) int {
	count, err := aq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (aq *AgentsQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, aq.ctx, ent.OpQueryExist)
	switch _, err := aq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (aq *AgentsQuery) ExistX(ctx context.Context) bool {
	exist, err := aq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the AgentsQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (aq *AgentsQuery) Clone() *AgentsQuery {
	if aq == nil {
		return nil
	}
	return &AgentsQuery{
		config:         aq.config,
		ctx:            aq.ctx.Clone(),
		order:          append([]agents.OrderOption{}, aq.order...),
		inters:         append([]Interceptor{}, aq.inters...),
		predicates:     append([]predicate.Agents{}, aq.predicates...),
		withAgentTasks: aq.withAgentTasks.Clone(),
		// clone intermediate query.
		sql:  aq.sql.Clone(),
		path: aq.path,
	}
}

// WithAgentTasks tells the query-builder to eager-load the nodes that are connected to
// the "agent_tasks" edge. The optional arguments are used to configure the query builder of the edge.
func (aq *AgentsQuery) WithAgentTasks(opts ...func(*AgentTasksQuery)) *AgentsQuery {
	query := (&AgentTasksClient{config: aq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	aq.withAgentTasks = query
	return aq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Status string `json:"status,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Agents.Query().
//		GroupBy(agents.FieldStatus).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (aq *AgentsQuery) GroupBy(field string, fields ...string) *AgentsGroupBy {
	aq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &AgentsGroupBy{build: aq}
	grbuild.flds = &aq.ctx.Fields
	grbuild.label = agents.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Status string `json:"status,omitempty"`
//	}
//
//	client.Agents.Query().
//		Select(agents.FieldStatus).
//		Scan(ctx, &v)
func (aq *AgentsQuery) Select(fields ...string) *AgentsSelect {
	aq.ctx.Fields = append(aq.ctx.Fields, fields...)
	sbuild := &AgentsSelect{AgentsQuery: aq}
	sbuild.label = agents.Label
	sbuild.flds, sbuild.scan = &aq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a AgentsSelect configured with the given aggregations.
func (aq *AgentsQuery) Aggregate(fns ...AggregateFunc) *AgentsSelect {
	return aq.Select().Aggregate(fns...)
}

func (aq *AgentsQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range aq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, aq); err != nil {
				return err
			}
		}
	}
	for _, f := range aq.ctx.Fields {
		if !agents.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if aq.path != nil {
		prev, err := aq.path(ctx)
		if err != nil {
			return err
		}
		aq.sql = prev
	}
	return nil
}

func (aq *AgentsQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Agents, error) {
	var (
		nodes       = []*Agents{}
		_spec       = aq.querySpec()
		loadedTypes = [1]bool{
			aq.withAgentTasks != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Agents).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Agents{config: aq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, aq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := aq.withAgentTasks; query != nil {
		if err := aq.loadAgentTasks(ctx, query, nodes,
			func(n *Agents) { n.Edges.AgentTasks = []*AgentTasks{} },
			func(n *Agents, e *AgentTasks) { n.Edges.AgentTasks = append(n.Edges.AgentTasks, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (aq *AgentsQuery) loadAgentTasks(ctx context.Context, query *AgentTasksQuery, nodes []*Agents, init func(*Agents), assign func(*Agents, *AgentTasks)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Agents)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(agenttasks.FieldAgentID)
	}
	query.Where(predicate.AgentTasks(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(agents.AgentTasksColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.AgentID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "agent_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (aq *AgentsQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := aq.querySpec()
	_spec.Node.Columns = aq.ctx.Fields
	if len(aq.ctx.Fields) > 0 {
		_spec.Unique = aq.ctx.Unique != nil && *aq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, aq.driver, _spec)
}

func (aq *AgentsQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(agents.Table, agents.Columns, sqlgraph.NewFieldSpec(agents.FieldID, field.TypeInt))
	_spec.From = aq.sql
	if unique := aq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if aq.path != nil {
		_spec.Unique = true
	}
	if fields := aq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, agents.FieldID)
		for i := range fields {
			if fields[i] != agents.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := aq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := aq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := aq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := aq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (aq *AgentsQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(aq.driver.Dialect())
	t1 := builder.Table(agents.Table)
	columns := aq.ctx.Fields
	if len(columns) == 0 {
		columns = agents.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if aq.sql != nil {
		selector = aq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if aq.ctx.Unique != nil && *aq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range aq.predicates {
		p(selector)
	}
	for _, p := range aq.order {
		p(selector)
	}
	if offset := aq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := aq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// AgentsGroupBy is the group-by builder for Agents entities.
type AgentsGroupBy struct {
	selector
	build *AgentsQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (agb *AgentsGroupBy) Aggregate(fns ...AggregateFunc) *AgentsGroupBy {
	agb.fns = append(agb.fns, fns...)
	return agb
}

// Scan applies the selector query and scans the result into the given value.
func (agb *AgentsGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, agb.build.ctx, ent.OpQueryGroupBy)
	if err := agb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*AgentsQuery, *AgentsGroupBy](ctx, agb.build, agb, agb.build.inters, v)
}

func (agb *AgentsGroupBy) sqlScan(ctx context.Context, root *AgentsQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(agb.fns))
	for _, fn := range agb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*agb.flds)+len(agb.fns))
		for _, f := range *agb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*agb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := agb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// AgentsSelect is the builder for selecting fields of Agents entities.
type AgentsSelect struct {
	*AgentsQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (as *AgentsSelect) Aggregate(fns ...AggregateFunc) *AgentsSelect {
	as.fns = append(as.fns, fns...)
	return as
}

// Scan applies the selector query and scans the result into the given value.
func (as *AgentsSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, as.ctx, ent.OpQuerySelect)
	if err := as.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*AgentsQuery, *AgentsSelect](ctx, as.AgentsQuery, as, as.inters, v)
}

func (as *AgentsSelect) sqlScan(ctx context.Context, root *AgentsQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(as.fns))
	for _, fn := range as.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*as.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := as.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
