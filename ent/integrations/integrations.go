// Code generated by ent, DO NOT EDIT.

package integrations

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the integrations type in the database.
	Label = "integrations"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldConfig holds the string denoting the config field in the database.
	FieldConfig = "config"
	// EdgeScans holds the string denoting the scans edge name in mutations.
	EdgeScans = "scans"
	// Table holds the table name of the integrations in the database.
	Table = "integrations"
	// ScansTable is the table that holds the scans relation/edge.
	ScansTable = "scans"
	// ScansInverseTable is the table name for the Scans entity.
	// It exists in this package in order to avoid circular dependency with the "scans" package.
	ScansInverseTable = "scans"
	// ScansColumn is the table column denoting the scans relation/edge.
	ScansColumn = "integration_id"
)

// Columns holds all SQL columns for integrations fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldType,
	FieldConfig,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "integrations"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"integration_type_integrations",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// TypeValidator is a validator for the "type" field. It is called by the builders before save.
	TypeValidator func(string) error
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the Integrations queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByType orders the results by the type field.
func ByType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldType, opts...).ToFunc()
}

// ByScansCount orders the results by scans count.
func ByScansCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newScansStep(), opts...)
	}
}

// ByScans orders the results by scans terms.
func ByScans(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newScansStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newScansStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ScansInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ScansTable, ScansColumn),
	)
}
