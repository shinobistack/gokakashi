// Code generated by ent, DO NOT EDIT.

package scannotify

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the scannotify type in the database.
	Label = "scan_notify"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldScanID holds the string denoting the scan_id field in the database.
	FieldScanID = "scan_id"
	// FieldHash holds the string denoting the hash field in the database.
	FieldHash = "hash"
	// EdgeScan holds the string denoting the scan edge name in mutations.
	EdgeScan = "scan"
	// Table holds the table name of the scannotify in the database.
	Table = "scan_notifies"
	// ScanTable is the table that holds the scan relation/edge.
	ScanTable = "scan_notifies"
	// ScanInverseTable is the table name for the Scans entity.
	// It exists in this package in order to avoid circular dependency with the "scans" package.
	ScanInverseTable = "scans"
	// ScanColumn is the table column denoting the scan relation/edge.
	ScanColumn = "scan_id"
)

// Columns holds all SQL columns for scannotify fields.
var Columns = []string{
	FieldID,
	FieldScanID,
	FieldHash,
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
	// HashValidator is a validator for the "hash" field. It is called by the builders before save.
	HashValidator func(string) error
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the ScanNotify queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByScanID orders the results by the scan_id field.
func ByScanID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldScanID, opts...).ToFunc()
}

// ByHash orders the results by the hash field.
func ByHash(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHash, opts...).ToFunc()
}

// ByScanField orders the results by scan field.
func ByScanField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newScanStep(), sql.OrderByField(field, opts...))
	}
}
func newScanStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ScanInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, ScanTable, ScanColumn),
	)
}