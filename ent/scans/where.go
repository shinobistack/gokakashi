// Code generated by ent, DO NOT EDIT.

package scans

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.Scans {
	return predicate.Scans(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.Scans {
	return predicate.Scans(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.Scans {
	return predicate.Scans(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.Scans {
	return predicate.Scans(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.Scans {
	return predicate.Scans(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.Scans {
	return predicate.Scans(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.Scans {
	return predicate.Scans(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.Scans {
	return predicate.Scans(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.Scans {
	return predicate.Scans(sql.FieldLTE(FieldID, id))
}

// PolicyID applies equality check predicate on the "policy_id" field. It's identical to PolicyIDEQ.
func PolicyID(v uuid.UUID) predicate.Scans {
	return predicate.Scans(sql.FieldEQ(FieldPolicyID, v))
}

// Status applies equality check predicate on the "status" field. It's identical to StatusEQ.
func Status(v string) predicate.Scans {
	return predicate.Scans(sql.FieldEQ(FieldStatus, v))
}

// PolicyIDEQ applies the EQ predicate on the "policy_id" field.
func PolicyIDEQ(v uuid.UUID) predicate.Scans {
	return predicate.Scans(sql.FieldEQ(FieldPolicyID, v))
}

// PolicyIDNEQ applies the NEQ predicate on the "policy_id" field.
func PolicyIDNEQ(v uuid.UUID) predicate.Scans {
	return predicate.Scans(sql.FieldNEQ(FieldPolicyID, v))
}

// PolicyIDIn applies the In predicate on the "policy_id" field.
func PolicyIDIn(vs ...uuid.UUID) predicate.Scans {
	return predicate.Scans(sql.FieldIn(FieldPolicyID, vs...))
}

// PolicyIDNotIn applies the NotIn predicate on the "policy_id" field.
func PolicyIDNotIn(vs ...uuid.UUID) predicate.Scans {
	return predicate.Scans(sql.FieldNotIn(FieldPolicyID, vs...))
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v string) predicate.Scans {
	return predicate.Scans(sql.FieldEQ(FieldStatus, v))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v string) predicate.Scans {
	return predicate.Scans(sql.FieldNEQ(FieldStatus, v))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...string) predicate.Scans {
	return predicate.Scans(sql.FieldIn(FieldStatus, vs...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...string) predicate.Scans {
	return predicate.Scans(sql.FieldNotIn(FieldStatus, vs...))
}

// StatusGT applies the GT predicate on the "status" field.
func StatusGT(v string) predicate.Scans {
	return predicate.Scans(sql.FieldGT(FieldStatus, v))
}

// StatusGTE applies the GTE predicate on the "status" field.
func StatusGTE(v string) predicate.Scans {
	return predicate.Scans(sql.FieldGTE(FieldStatus, v))
}

// StatusLT applies the LT predicate on the "status" field.
func StatusLT(v string) predicate.Scans {
	return predicate.Scans(sql.FieldLT(FieldStatus, v))
}

// StatusLTE applies the LTE predicate on the "status" field.
func StatusLTE(v string) predicate.Scans {
	return predicate.Scans(sql.FieldLTE(FieldStatus, v))
}

// StatusContains applies the Contains predicate on the "status" field.
func StatusContains(v string) predicate.Scans {
	return predicate.Scans(sql.FieldContains(FieldStatus, v))
}

// StatusHasPrefix applies the HasPrefix predicate on the "status" field.
func StatusHasPrefix(v string) predicate.Scans {
	return predicate.Scans(sql.FieldHasPrefix(FieldStatus, v))
}

// StatusHasSuffix applies the HasSuffix predicate on the "status" field.
func StatusHasSuffix(v string) predicate.Scans {
	return predicate.Scans(sql.FieldHasSuffix(FieldStatus, v))
}

// StatusEqualFold applies the EqualFold predicate on the "status" field.
func StatusEqualFold(v string) predicate.Scans {
	return predicate.Scans(sql.FieldEqualFold(FieldStatus, v))
}

// StatusContainsFold applies the ContainsFold predicate on the "status" field.
func StatusContainsFold(v string) predicate.Scans {
	return predicate.Scans(sql.FieldContainsFold(FieldStatus, v))
}

// CheckIsNil applies the IsNil predicate on the "check" field.
func CheckIsNil() predicate.Scans {
	return predicate.Scans(sql.FieldIsNull(FieldCheck))
}

// CheckNotNil applies the NotNil predicate on the "check" field.
func CheckNotNil() predicate.Scans {
	return predicate.Scans(sql.FieldNotNull(FieldCheck))
}

// ReportIsNil applies the IsNil predicate on the "report" field.
func ReportIsNil() predicate.Scans {
	return predicate.Scans(sql.FieldIsNull(FieldReport))
}

// ReportNotNil applies the NotNil predicate on the "report" field.
func ReportNotNil() predicate.Scans {
	return predicate.Scans(sql.FieldNotNull(FieldReport))
}

// HasPolicy applies the HasEdge predicate on the "policy" edge.
func HasPolicy() predicate.Scans {
	return predicate.Scans(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, PolicyTable, PolicyColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPolicyWith applies the HasEdge predicate on the "policy" edge with a given conditions (other predicates).
func HasPolicyWith(preds ...predicate.Policies) predicate.Scans {
	return predicate.Scans(func(s *sql.Selector) {
		step := newPolicyStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasScanLabels applies the HasEdge predicate on the "scan_labels" edge.
func HasScanLabels() predicate.Scans {
	return predicate.Scans(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ScanLabelsTable, ScanLabelsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasScanLabelsWith applies the HasEdge predicate on the "scan_labels" edge with a given conditions (other predicates).
func HasScanLabelsWith(preds ...predicate.ScanLabels) predicate.Scans {
	return predicate.Scans(func(s *sql.Selector) {
		step := newScanLabelsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Scans) predicate.Scans {
	return predicate.Scans(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Scans) predicate.Scans {
	return predicate.Scans(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Scans) predicate.Scans {
	return predicate.Scans(sql.NotPredicates(p))
}
