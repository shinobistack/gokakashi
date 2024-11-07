package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// ScanPolicy holds the schema definition for the ScanPolicy entity.
type ScanPolicy struct {
	ent.Schema
}

// Fields of the ScanPolicy.
func (ScanPolicy) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("trigger_type"),
	}
}

// Edges of the ScanPolicy.
func (ScanPolicy) Edges() []ent.Edge {
	return nil
}
