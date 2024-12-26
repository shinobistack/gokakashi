// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/agents"
	"github.com/shinobistack/gokakashi/ent/agenttasks"
	"github.com/shinobistack/gokakashi/ent/scans"
)

// AgentTasks is the model entity for the AgentTasks schema.
type AgentTasks struct {
	config `json:"-"`
	// ID of the ent.
	// Primary key, unique identifier.
	ID int `json:"id,omitempty"`
	// Foreign key to Agents.ID.
	AgentID int `json:"agent_id,omitempty"`
	// Foreign key to Scans.ID.
	ScanID uuid.UUID `json:"scan_id,omitempty"`
	// Enum: { pending, in_progress, complete }.
	Status string `json:"status,omitempty"`
	// Timestamp for task creation.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the AgentTasksQuery when eager-loading is set.
	Edges        AgentTasksEdges `json:"edges"`
	selectValues sql.SelectValues
}

// AgentTasksEdges holds the relations/edges for other nodes in the graph.
type AgentTasksEdges struct {
	// Agent holds the value of the agent edge.
	Agent *Agents `json:"agent,omitempty"`
	// Scan holds the value of the scan edge.
	Scan *Scans `json:"scan,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// AgentOrErr returns the Agent value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AgentTasksEdges) AgentOrErr() (*Agents, error) {
	if e.Agent != nil {
		return e.Agent, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: agents.Label}
	}
	return nil, &NotLoadedError{edge: "agent"}
}

// ScanOrErr returns the Scan value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AgentTasksEdges) ScanOrErr() (*Scans, error) {
	if e.Scan != nil {
		return e.Scan, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: scans.Label}
	}
	return nil, &NotLoadedError{edge: "scan"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*AgentTasks) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case agenttasks.FieldID, agenttasks.FieldAgentID:
			values[i] = new(sql.NullInt64)
		case agenttasks.FieldStatus:
			values[i] = new(sql.NullString)
		case agenttasks.FieldCreatedAt:
			values[i] = new(sql.NullTime)
		case agenttasks.FieldScanID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the AgentTasks fields.
func (at *AgentTasks) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case agenttasks.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			at.ID = int(value.Int64)
		case agenttasks.FieldAgentID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field agent_id", values[i])
			} else if value.Valid {
				at.AgentID = int(value.Int64)
			}
		case agenttasks.FieldScanID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field scan_id", values[i])
			} else if value != nil {
				at.ScanID = *value
			}
		case agenttasks.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				at.Status = value.String
			}
		case agenttasks.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				at.CreatedAt = value.Time
			}
		default:
			at.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the AgentTasks.
// This includes values selected through modifiers, order, etc.
func (at *AgentTasks) Value(name string) (ent.Value, error) {
	return at.selectValues.Get(name)
}

// QueryAgent queries the "agent" edge of the AgentTasks entity.
func (at *AgentTasks) QueryAgent() *AgentsQuery {
	return NewAgentTasksClient(at.config).QueryAgent(at)
}

// QueryScan queries the "scan" edge of the AgentTasks entity.
func (at *AgentTasks) QueryScan() *ScansQuery {
	return NewAgentTasksClient(at.config).QueryScan(at)
}

// Update returns a builder for updating this AgentTasks.
// Note that you need to call AgentTasks.Unwrap() before calling this method if this AgentTasks
// was returned from a transaction, and the transaction was committed or rolled back.
func (at *AgentTasks) Update() *AgentTasksUpdateOne {
	return NewAgentTasksClient(at.config).UpdateOne(at)
}

// Unwrap unwraps the AgentTasks entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (at *AgentTasks) Unwrap() *AgentTasks {
	_tx, ok := at.config.driver.(*txDriver)
	if !ok {
		panic("ent: AgentTasks is not a transactional entity")
	}
	at.config.driver = _tx.drv
	return at
}

// String implements the fmt.Stringer.
func (at *AgentTasks) String() string {
	var builder strings.Builder
	builder.WriteString("AgentTasks(")
	builder.WriteString(fmt.Sprintf("id=%v, ", at.ID))
	builder.WriteString("agent_id=")
	builder.WriteString(fmt.Sprintf("%v", at.AgentID))
	builder.WriteString(", ")
	builder.WriteString("scan_id=")
	builder.WriteString(fmt.Sprintf("%v", at.ScanID))
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(at.Status)
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(at.CreatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// AgentTasksSlice is a parsable slice of AgentTasks.
type AgentTasksSlice []*AgentTasks