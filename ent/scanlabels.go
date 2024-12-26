// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/scanlabels"
	"github.com/shinobistack/gokakashi/ent/scans"
)

// ScanLabels is the model entity for the ScanLabels schema.
type ScanLabels struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Foreign key to Scans.ID.
	ScanID uuid.UUID `json:"scan_id,omitempty"`
	// Label key.
	Key string `json:"key,omitempty"`
	// Label value.
	Value string `json:"value,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ScanLabelsQuery when eager-loading is set.
	Edges        ScanLabelsEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ScanLabelsEdges holds the relations/edges for other nodes in the graph.
type ScanLabelsEdges struct {
	// Scan holds the value of the scan edge.
	Scan *Scans `json:"scan,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// ScanOrErr returns the Scan value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ScanLabelsEdges) ScanOrErr() (*Scans, error) {
	if e.Scan != nil {
		return e.Scan, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: scans.Label}
	}
	return nil, &NotLoadedError{edge: "scan"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ScanLabels) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case scanlabels.FieldID:
			values[i] = new(sql.NullInt64)
		case scanlabels.FieldKey, scanlabels.FieldValue:
			values[i] = new(sql.NullString)
		case scanlabels.FieldScanID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ScanLabels fields.
func (sl *ScanLabels) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case scanlabels.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			sl.ID = int(value.Int64)
		case scanlabels.FieldScanID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field scan_id", values[i])
			} else if value != nil {
				sl.ScanID = *value
			}
		case scanlabels.FieldKey:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field key", values[i])
			} else if value.Valid {
				sl.Key = value.String
			}
		case scanlabels.FieldValue:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field value", values[i])
			} else if value.Valid {
				sl.Value = value.String
			}
		default:
			sl.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// GetValue returns the ent.Value that was dynamically selected and assigned to the ScanLabels.
// This includes values selected through modifiers, order, etc.
func (sl *ScanLabels) GetValue(name string) (ent.Value, error) {
	return sl.selectValues.Get(name)
}

// QueryScan queries the "scan" edge of the ScanLabels entity.
func (sl *ScanLabels) QueryScan() *ScansQuery {
	return NewScanLabelsClient(sl.config).QueryScan(sl)
}

// Update returns a builder for updating this ScanLabels.
// Note that you need to call ScanLabels.Unwrap() before calling this method if this ScanLabels
// was returned from a transaction, and the transaction was committed or rolled back.
func (sl *ScanLabels) Update() *ScanLabelsUpdateOne {
	return NewScanLabelsClient(sl.config).UpdateOne(sl)
}

// Unwrap unwraps the ScanLabels entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (sl *ScanLabels) Unwrap() *ScanLabels {
	_tx, ok := sl.config.driver.(*txDriver)
	if !ok {
		panic("ent: ScanLabels is not a transactional entity")
	}
	sl.config.driver = _tx.drv
	return sl
}

// String implements the fmt.Stringer.
func (sl *ScanLabels) String() string {
	var builder strings.Builder
	builder.WriteString("ScanLabels(")
	builder.WriteString(fmt.Sprintf("id=%v, ", sl.ID))
	builder.WriteString("scan_id=")
	builder.WriteString(fmt.Sprintf("%v", sl.ScanID))
	builder.WriteString(", ")
	builder.WriteString("key=")
	builder.WriteString(sl.Key)
	builder.WriteString(", ")
	builder.WriteString("value=")
	builder.WriteString(sl.Value)
	builder.WriteByte(')')
	return builder.String()
}

// ScanLabelsSlice is a parsable slice of ScanLabels.
type ScanLabelsSlice []*ScanLabels