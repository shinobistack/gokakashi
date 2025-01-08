// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/scannotify"
	"github.com/shinobistack/gokakashi/ent/scans"
)

// ScanNotify is the model entity for the ScanNotify schema.
type ScanNotify struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// Foreign key to the scans table
	ScanID uuid.UUID `json:"scan_id,omitempty"`
	// Unique hash for condition evaluation and vulnerabilities
	Hash string `json:"hash,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ScanNotifyQuery when eager-loading is set.
	Edges        ScanNotifyEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ScanNotifyEdges holds the relations/edges for other nodes in the graph.
type ScanNotifyEdges struct {
	// Links the notification to its corresponding scan
	Scan *Scans `json:"scan,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// ScanOrErr returns the Scan value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ScanNotifyEdges) ScanOrErr() (*Scans, error) {
	if e.Scan != nil {
		return e.Scan, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: scans.Label}
	}
	return nil, &NotLoadedError{edge: "scan"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ScanNotify) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case scannotify.FieldHash:
			values[i] = new(sql.NullString)
		case scannotify.FieldID, scannotify.FieldScanID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ScanNotify fields.
func (sn *ScanNotify) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case scannotify.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				sn.ID = *value
			}
		case scannotify.FieldScanID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field scan_id", values[i])
			} else if value != nil {
				sn.ScanID = *value
			}
		case scannotify.FieldHash:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field hash", values[i])
			} else if value.Valid {
				sn.Hash = value.String
			}
		default:
			sn.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the ScanNotify.
// This includes values selected through modifiers, order, etc.
func (sn *ScanNotify) Value(name string) (ent.Value, error) {
	return sn.selectValues.Get(name)
}

// QueryScan queries the "scan" edge of the ScanNotify entity.
func (sn *ScanNotify) QueryScan() *ScansQuery {
	return NewScanNotifyClient(sn.config).QueryScan(sn)
}

// Update returns a builder for updating this ScanNotify.
// Note that you need to call ScanNotify.Unwrap() before calling this method if this ScanNotify
// was returned from a transaction, and the transaction was committed or rolled back.
func (sn *ScanNotify) Update() *ScanNotifyUpdateOne {
	return NewScanNotifyClient(sn.config).UpdateOne(sn)
}

// Unwrap unwraps the ScanNotify entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (sn *ScanNotify) Unwrap() *ScanNotify {
	_tx, ok := sn.config.driver.(*txDriver)
	if !ok {
		panic("ent: ScanNotify is not a transactional entity")
	}
	sn.config.driver = _tx.drv
	return sn
}

// String implements the fmt.Stringer.
func (sn *ScanNotify) String() string {
	var builder strings.Builder
	builder.WriteString("ScanNotify(")
	builder.WriteString(fmt.Sprintf("id=%v, ", sn.ID))
	builder.WriteString("scan_id=")
	builder.WriteString(fmt.Sprintf("%v", sn.ScanID))
	builder.WriteString(", ")
	builder.WriteString("hash=")
	builder.WriteString(sn.Hash)
	builder.WriteByte(')')
	return builder.String()
}

// ScanNotifies is a parsable slice of ScanNotify.
type ScanNotifies []*ScanNotify