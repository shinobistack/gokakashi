// Code generated by ent, DO NOT EDIT.

package scans

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the scans type in the database.
	Label = "scans"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldPolicyID holds the string denoting the policy_id field in the database.
	FieldPolicyID = "policy_id"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldImage holds the string denoting the image field in the database.
	FieldImage = "image"
	// FieldCheck holds the string denoting the check field in the database.
	FieldCheck = "check"
	// FieldReport holds the string denoting the report field in the database.
	FieldReport = "report"
	// EdgePolicy holds the string denoting the policy edge name in mutations.
	EdgePolicy = "policy"
	// EdgeScanLabels holds the string denoting the scan_labels edge name in mutations.
	EdgeScanLabels = "scan_labels"
	// EdgeAgentTasks holds the string denoting the agent_tasks edge name in mutations.
	EdgeAgentTasks = "agent_tasks"
	// Table holds the table name of the scans in the database.
	Table = "scans"
	// PolicyTable is the table that holds the policy relation/edge.
	PolicyTable = "scans"
	// PolicyInverseTable is the table name for the Policies entity.
	// It exists in this package in order to avoid circular dependency with the "policies" package.
	PolicyInverseTable = "policies"
	// PolicyColumn is the table column denoting the policy relation/edge.
	PolicyColumn = "policy_id"
	// ScanLabelsTable is the table that holds the scan_labels relation/edge.
	ScanLabelsTable = "scan_labels"
	// ScanLabelsInverseTable is the table name for the ScanLabels entity.
	// It exists in this package in order to avoid circular dependency with the "scanlabels" package.
	ScanLabelsInverseTable = "scan_labels"
	// ScanLabelsColumn is the table column denoting the scan_labels relation/edge.
	ScanLabelsColumn = "scan_id"
	// AgentTasksTable is the table that holds the agent_tasks relation/edge.
	AgentTasksTable = "agent_tasks"
	// AgentTasksInverseTable is the table name for the AgentTasks entity.
	// It exists in this package in order to avoid circular dependency with the "agenttasks" package.
	AgentTasksInverseTable = "agent_tasks"
	// AgentTasksColumn is the table column denoting the agent_tasks relation/edge.
	AgentTasksColumn = "scan_id"
)

// Columns holds all SQL columns for scans fields.
var Columns = []string{
	FieldID,
	FieldPolicyID,
	FieldStatus,
	FieldImage,
	FieldCheck,
	FieldReport,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultStatus holds the default value on creation for the "status" field.
	DefaultStatus string
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the Scans queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByPolicyID orders the results by the policy_id field.
func ByPolicyID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPolicyID, opts...).ToFunc()
}

// ByStatus orders the results by the status field.
func ByStatus(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStatus, opts...).ToFunc()
}

// ByImage orders the results by the image field.
func ByImage(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldImage, opts...).ToFunc()
}

// ByReport orders the results by the report field.
func ByReport(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldReport, opts...).ToFunc()
}

// ByPolicyField orders the results by policy field.
func ByPolicyField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPolicyStep(), sql.OrderByField(field, opts...))
	}
}

// ByScanLabelsCount orders the results by scan_labels count.
func ByScanLabelsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newScanLabelsStep(), opts...)
	}
}

// ByScanLabels orders the results by scan_labels terms.
func ByScanLabels(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newScanLabelsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByAgentTasksCount orders the results by agent_tasks count.
func ByAgentTasksCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newAgentTasksStep(), opts...)
	}
}

// ByAgentTasks orders the results by agent_tasks terms.
func ByAgentTasks(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newAgentTasksStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newPolicyStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PolicyInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, PolicyTable, PolicyColumn),
	)
}
func newScanLabelsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ScanLabelsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ScanLabelsTable, ScanLabelsColumn),
	)
}
func newAgentTasksStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(AgentTasksInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, AgentTasksTable, AgentTasksColumn),
	)
}
