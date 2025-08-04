package schema

import (
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/internal/scan/v2"
)

// Scans holds the schema definition for the Scans entity.
type V2Scans struct {
	ent.Schema
}

// Fields of the Scans.
func (V2Scans) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Comment("Primary key, unique identifier."),

		field.String("status").
			Default(string(scan.Pending)).
			Validate(scan.ValidateStatus).
			Comment("Enum: " + strings.Join(scan.Statuses(), ", ")),

		field.String("image").
			Comment("Details of the image being scanned."),

		field.JSON("labels", map[string]string{}).
			Optional().
			Comment("Scan labels key:value"),

		field.Time("created_at").
			Default(time.Now).
			Comment("Scan creation time"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Scan update time"),
	}
}

// Edges of the Scans.
func (V2Scans) Edges() []ent.Edge {
	return []ent.Edge{}
}

// Annotations of the Scans.
func (V2Scans) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "v2_scans"},
	}
}
