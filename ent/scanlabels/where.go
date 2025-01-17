// Code generated by ent, DO NOT EDIT.

package scanlabels

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldLTE(FieldID, id))
}

// ScanID applies equality check predicate on the "scan_id" field. It's identical to ScanIDEQ.
func ScanID(v uuid.UUID) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldEQ(FieldScanID, v))
}

// Key applies equality check predicate on the "key" field. It's identical to KeyEQ.
func Key(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldEQ(FieldKey, v))
}

// Value applies equality check predicate on the "value" field. It's identical to ValueEQ.
func Value(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldEQ(FieldValue, v))
}

// ScanIDEQ applies the EQ predicate on the "scan_id" field.
func ScanIDEQ(v uuid.UUID) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldEQ(FieldScanID, v))
}

// ScanIDNEQ applies the NEQ predicate on the "scan_id" field.
func ScanIDNEQ(v uuid.UUID) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldNEQ(FieldScanID, v))
}

// ScanIDIn applies the In predicate on the "scan_id" field.
func ScanIDIn(vs ...uuid.UUID) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldIn(FieldScanID, vs...))
}

// ScanIDNotIn applies the NotIn predicate on the "scan_id" field.
func ScanIDNotIn(vs ...uuid.UUID) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldNotIn(FieldScanID, vs...))
}

// KeyEQ applies the EQ predicate on the "key" field.
func KeyEQ(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldEQ(FieldKey, v))
}

// KeyNEQ applies the NEQ predicate on the "key" field.
func KeyNEQ(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldNEQ(FieldKey, v))
}

// KeyIn applies the In predicate on the "key" field.
func KeyIn(vs ...string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldIn(FieldKey, vs...))
}

// KeyNotIn applies the NotIn predicate on the "key" field.
func KeyNotIn(vs ...string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldNotIn(FieldKey, vs...))
}

// KeyGT applies the GT predicate on the "key" field.
func KeyGT(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldGT(FieldKey, v))
}

// KeyGTE applies the GTE predicate on the "key" field.
func KeyGTE(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldGTE(FieldKey, v))
}

// KeyLT applies the LT predicate on the "key" field.
func KeyLT(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldLT(FieldKey, v))
}

// KeyLTE applies the LTE predicate on the "key" field.
func KeyLTE(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldLTE(FieldKey, v))
}

// KeyContains applies the Contains predicate on the "key" field.
func KeyContains(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldContains(FieldKey, v))
}

// KeyHasPrefix applies the HasPrefix predicate on the "key" field.
func KeyHasPrefix(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldHasPrefix(FieldKey, v))
}

// KeyHasSuffix applies the HasSuffix predicate on the "key" field.
func KeyHasSuffix(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldHasSuffix(FieldKey, v))
}

// KeyEqualFold applies the EqualFold predicate on the "key" field.
func KeyEqualFold(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldEqualFold(FieldKey, v))
}

// KeyContainsFold applies the ContainsFold predicate on the "key" field.
func KeyContainsFold(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldContainsFold(FieldKey, v))
}

// ValueEQ applies the EQ predicate on the "value" field.
func ValueEQ(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldEQ(FieldValue, v))
}

// ValueNEQ applies the NEQ predicate on the "value" field.
func ValueNEQ(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldNEQ(FieldValue, v))
}

// ValueIn applies the In predicate on the "value" field.
func ValueIn(vs ...string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldIn(FieldValue, vs...))
}

// ValueNotIn applies the NotIn predicate on the "value" field.
func ValueNotIn(vs ...string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldNotIn(FieldValue, vs...))
}

// ValueGT applies the GT predicate on the "value" field.
func ValueGT(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldGT(FieldValue, v))
}

// ValueGTE applies the GTE predicate on the "value" field.
func ValueGTE(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldGTE(FieldValue, v))
}

// ValueLT applies the LT predicate on the "value" field.
func ValueLT(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldLT(FieldValue, v))
}

// ValueLTE applies the LTE predicate on the "value" field.
func ValueLTE(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldLTE(FieldValue, v))
}

// ValueContains applies the Contains predicate on the "value" field.
func ValueContains(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldContains(FieldValue, v))
}

// ValueHasPrefix applies the HasPrefix predicate on the "value" field.
func ValueHasPrefix(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldHasPrefix(FieldValue, v))
}

// ValueHasSuffix applies the HasSuffix predicate on the "value" field.
func ValueHasSuffix(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldHasSuffix(FieldValue, v))
}

// ValueEqualFold applies the EqualFold predicate on the "value" field.
func ValueEqualFold(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldEqualFold(FieldValue, v))
}

// ValueContainsFold applies the ContainsFold predicate on the "value" field.
func ValueContainsFold(v string) predicate.ScanLabels {
	return predicate.ScanLabels(sql.FieldContainsFold(FieldValue, v))
}

// HasScan applies the HasEdge predicate on the "scan" edge.
func HasScan() predicate.ScanLabels {
	return predicate.ScanLabels(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ScanTable, ScanColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasScanWith applies the HasEdge predicate on the "scan" edge with a given conditions (other predicates).
func HasScanWith(preds ...predicate.Scans) predicate.ScanLabels {
	return predicate.ScanLabels(func(s *sql.Selector) {
		step := newScanStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ScanLabels) predicate.ScanLabels {
	return predicate.ScanLabels(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ScanLabels) predicate.ScanLabels {
	return predicate.ScanLabels(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.ScanLabels) predicate.ScanLabels {
	return predicate.ScanLabels(sql.NotPredicates(p))
}
