// Code generated by ent, DO NOT EDIT.

package scannotify

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldLTE(FieldID, id))
}

// ScanID applies equality check predicate on the "scan_id" field. It's identical to ScanIDEQ.
func ScanID(v uuid.UUID) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldEQ(FieldScanID, v))
}

// Hash applies equality check predicate on the "hash" field. It's identical to HashEQ.
func Hash(v string) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldEQ(FieldHash, v))
}

// ScanIDEQ applies the EQ predicate on the "scan_id" field.
func ScanIDEQ(v uuid.UUID) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldEQ(FieldScanID, v))
}

// ScanIDNEQ applies the NEQ predicate on the "scan_id" field.
func ScanIDNEQ(v uuid.UUID) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldNEQ(FieldScanID, v))
}

// ScanIDIn applies the In predicate on the "scan_id" field.
func ScanIDIn(vs ...uuid.UUID) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldIn(FieldScanID, vs...))
}

// ScanIDNotIn applies the NotIn predicate on the "scan_id" field.
func ScanIDNotIn(vs ...uuid.UUID) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldNotIn(FieldScanID, vs...))
}

// HashEQ applies the EQ predicate on the "hash" field.
func HashEQ(v string) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldEQ(FieldHash, v))
}

// HashNEQ applies the NEQ predicate on the "hash" field.
func HashNEQ(v string) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldNEQ(FieldHash, v))
}

// HashIn applies the In predicate on the "hash" field.
func HashIn(vs ...string) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldIn(FieldHash, vs...))
}

// HashNotIn applies the NotIn predicate on the "hash" field.
func HashNotIn(vs ...string) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldNotIn(FieldHash, vs...))
}

// HashGT applies the GT predicate on the "hash" field.
func HashGT(v string) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldGT(FieldHash, v))
}

// HashGTE applies the GTE predicate on the "hash" field.
func HashGTE(v string) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldGTE(FieldHash, v))
}

// HashLT applies the LT predicate on the "hash" field.
func HashLT(v string) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldLT(FieldHash, v))
}

// HashLTE applies the LTE predicate on the "hash" field.
func HashLTE(v string) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldLTE(FieldHash, v))
}

// HashContains applies the Contains predicate on the "hash" field.
func HashContains(v string) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldContains(FieldHash, v))
}

// HashHasPrefix applies the HasPrefix predicate on the "hash" field.
func HashHasPrefix(v string) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldHasPrefix(FieldHash, v))
}

// HashHasSuffix applies the HasSuffix predicate on the "hash" field.
func HashHasSuffix(v string) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldHasSuffix(FieldHash, v))
}

// HashEqualFold applies the EqualFold predicate on the "hash" field.
func HashEqualFold(v string) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldEqualFold(FieldHash, v))
}

// HashContainsFold applies the ContainsFold predicate on the "hash" field.
func HashContainsFold(v string) predicate.ScanNotify {
	return predicate.ScanNotify(sql.FieldContainsFold(FieldHash, v))
}

// HasScan applies the HasEdge predicate on the "scan" edge.
func HasScan() predicate.ScanNotify {
	return predicate.ScanNotify(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ScanTable, ScanColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasScanWith applies the HasEdge predicate on the "scan" edge with a given conditions (other predicates).
func HasScanWith(preds ...predicate.Scans) predicate.ScanNotify {
	return predicate.ScanNotify(func(s *sql.Selector) {
		step := newScanStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ScanNotify) predicate.ScanNotify {
	return predicate.ScanNotify(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ScanNotify) predicate.ScanNotify {
	return predicate.ScanNotify(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.ScanNotify) predicate.ScanNotify {
	return predicate.ScanNotify(sql.NotPredicates(p))
}
