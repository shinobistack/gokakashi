// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/policies"
	"github.com/shinobistack/gokakashi/ent/policylabels"
)

// PolicyLabels is the model entity for the PolicyLabels schema.
type PolicyLabels struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Foreign key to policies
	PolicyID uuid.UUID `json:"policy_id,omitempty"`
	// Key holds the value of the "key" field.
	Key string `json:"key,omitempty"`
	// Value holds the value of the "value" field.
	Value string `json:"value,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the PolicyLabelsQuery when eager-loading is set.
	Edges        PolicyLabelsEdges `json:"edges"`
	selectValues sql.SelectValues
}

// PolicyLabelsEdges holds the relations/edges for other nodes in the graph.
type PolicyLabelsEdges struct {
	// Policy holds the value of the policy edge.
	Policy *Policies `json:"policy,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// PolicyOrErr returns the Policy value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PolicyLabelsEdges) PolicyOrErr() (*Policies, error) {
	if e.Policy != nil {
		return e.Policy, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: policies.Label}
	}
	return nil, &NotLoadedError{edge: "policy"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*PolicyLabels) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case policylabels.FieldID:
			values[i] = new(sql.NullInt64)
		case policylabels.FieldKey, policylabels.FieldValue:
			values[i] = new(sql.NullString)
		case policylabels.FieldPolicyID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the PolicyLabels fields.
func (pl *PolicyLabels) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case policylabels.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			pl.ID = int(value.Int64)
		case policylabels.FieldPolicyID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field policy_id", values[i])
			} else if value != nil {
				pl.PolicyID = *value
			}
		case policylabels.FieldKey:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field key", values[i])
			} else if value.Valid {
				pl.Key = value.String
			}
		case policylabels.FieldValue:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field value", values[i])
			} else if value.Valid {
				pl.Value = value.String
			}
		default:
			pl.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// GetValue returns the ent.Value that was dynamically selected and assigned to the PolicyLabels.
// This includes values selected through modifiers, order, etc.
func (pl *PolicyLabels) GetValue(name string) (ent.Value, error) {
	return pl.selectValues.Get(name)
}

// QueryPolicy queries the "policy" edge of the PolicyLabels entity.
func (pl *PolicyLabels) QueryPolicy() *PoliciesQuery {
	return NewPolicyLabelsClient(pl.config).QueryPolicy(pl)
}

// Update returns a builder for updating this PolicyLabels.
// Note that you need to call PolicyLabels.Unwrap() before calling this method if this PolicyLabels
// was returned from a transaction, and the transaction was committed or rolled back.
func (pl *PolicyLabels) Update() *PolicyLabelsUpdateOne {
	return NewPolicyLabelsClient(pl.config).UpdateOne(pl)
}

// Unwrap unwraps the PolicyLabels entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pl *PolicyLabels) Unwrap() *PolicyLabels {
	_tx, ok := pl.config.driver.(*txDriver)
	if !ok {
		panic("ent: PolicyLabels is not a transactional entity")
	}
	pl.config.driver = _tx.drv
	return pl
}

// String implements the fmt.Stringer.
func (pl *PolicyLabels) String() string {
	var builder strings.Builder
	builder.WriteString("PolicyLabels(")
	builder.WriteString(fmt.Sprintf("id=%v, ", pl.ID))
	builder.WriteString("policy_id=")
	builder.WriteString(fmt.Sprintf("%v", pl.PolicyID))
	builder.WriteString(", ")
	builder.WriteString("key=")
	builder.WriteString(pl.Key)
	builder.WriteString(", ")
	builder.WriteString("value=")
	builder.WriteString(pl.Value)
	builder.WriteByte(')')
	return builder.String()
}

// PolicyLabelsSlice is a parsable slice of PolicyLabels.
type PolicyLabelsSlice []*PolicyLabels
