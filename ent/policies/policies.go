// Code generated by ent, DO NOT EDIT.

package policies

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the policies type in the database.
	Label = "policies"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldImage holds the string denoting the image field in the database.
	FieldImage = "image"
	// FieldTrigger holds the string denoting the trigger field in the database.
	FieldTrigger = "trigger"
	// FieldCheck holds the string denoting the check field in the database.
	FieldCheck = "check"
	// EdgePolicyLabels holds the string denoting the policy_labels edge name in mutations.
	EdgePolicyLabels = "policy_labels"
	// Table holds the table name of the policies in the database.
	Table = "policies"
	// PolicyLabelsTable is the table that holds the policy_labels relation/edge.
	PolicyLabelsTable = "policy_labels"
	// PolicyLabelsInverseTable is the table name for the PolicyLabels entity.
	// It exists in this package in order to avoid circular dependency with the "policylabels" package.
	PolicyLabelsInverseTable = "policy_labels"
	// PolicyLabelsColumn is the table column denoting the policy_labels relation/edge.
	PolicyLabelsColumn = "policy_id"
)

// Columns holds all SQL columns for policies fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldImage,
	FieldTrigger,
	FieldCheck,
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
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the Policies queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByPolicyLabelsCount orders the results by policy_labels count.
func ByPolicyLabelsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newPolicyLabelsStep(), opts...)
	}
}

// ByPolicyLabels orders the results by policy_labels terms.
func ByPolicyLabels(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPolicyLabelsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newPolicyLabelsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PolicyLabelsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, PolicyLabelsTable, PolicyLabelsColumn),
	)
}
