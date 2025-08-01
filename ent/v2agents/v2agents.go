// Code generated by ent, DO NOT EDIT.

package v2agents

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the v2agents type in the database.
	Label = "v2agents"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldLastHeartbeatAt holds the string denoting the last_heartbeat_at field in the database.
	FieldLastHeartbeatAt = "last_heartbeat_at"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// Table holds the table name of the v2agents in the database.
	Table = "v2_agents"
)

// Columns holds all SQL columns for v2agents fields.
var Columns = []string{
	FieldID,
	FieldStatus,
	FieldLastHeartbeatAt,
	FieldCreatedAt,
	FieldUpdatedAt,
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
	// StatusValidator is a validator for the "status" field. It is called by the builders before save.
	StatusValidator func(string) error
	// DefaultLastHeartbeatAt holds the default value on creation for the "last_heartbeat_at" field.
	DefaultLastHeartbeatAt func() time.Time
	// UpdateDefaultLastHeartbeatAt holds the default value on update for the "last_heartbeat_at" field.
	UpdateDefaultLastHeartbeatAt func() time.Time
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the V2Agents queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByStatus orders the results by the status field.
func ByStatus(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStatus, opts...).ToFunc()
}

// ByLastHeartbeatAt orders the results by the last_heartbeat_at field.
func ByLastHeartbeatAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLastHeartbeatAt, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}
