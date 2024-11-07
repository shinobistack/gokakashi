package schema

import "entgo.io/ent"

// ScanPolicy holds the schema definition for the ScanPolicy entity.
type ScanPolicy struct {
	ent.Schema
}

// Fields of the ScanPolicy.
func (ScanPolicy) Fields() []ent.Field {
	return nil
}

// Edges of the ScanPolicy.
func (ScanPolicy) Edges() []ent.Edge {
	return nil
}
