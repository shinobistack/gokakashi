// Code generated by ent, DO NOT EDIT.

package policies

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.Policies {
	return predicate.Policies(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.Policies {
	return predicate.Policies(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.Policies {
	return predicate.Policies(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.Policies {
	return predicate.Policies(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.Policies {
	return predicate.Policies(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.Policies {
	return predicate.Policies(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.Policies {
	return predicate.Policies(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.Policies {
	return predicate.Policies(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.Policies {
	return predicate.Policies(sql.FieldLTE(FieldID, id))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Policies {
	return predicate.Policies(sql.FieldEQ(FieldName, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Policies {
	return predicate.Policies(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Policies {
	return predicate.Policies(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Policies {
	return predicate.Policies(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Policies {
	return predicate.Policies(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Policies {
	return predicate.Policies(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Policies {
	return predicate.Policies(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Policies {
	return predicate.Policies(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Policies {
	return predicate.Policies(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Policies {
	return predicate.Policies(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Policies {
	return predicate.Policies(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Policies {
	return predicate.Policies(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Policies {
	return predicate.Policies(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Policies {
	return predicate.Policies(sql.FieldContainsFold(FieldName, v))
}

// LabelsIsNil applies the IsNil predicate on the "labels" field.
func LabelsIsNil() predicate.Policies {
	return predicate.Policies(sql.FieldIsNull(FieldLabels))
}

// LabelsNotNil applies the NotNil predicate on the "labels" field.
func LabelsNotNil() predicate.Policies {
	return predicate.Policies(sql.FieldNotNull(FieldLabels))
}

// TriggerIsNil applies the IsNil predicate on the "trigger" field.
func TriggerIsNil() predicate.Policies {
	return predicate.Policies(sql.FieldIsNull(FieldTrigger))
}

// TriggerNotNil applies the NotNil predicate on the "trigger" field.
func TriggerNotNil() predicate.Policies {
	return predicate.Policies(sql.FieldNotNull(FieldTrigger))
}

// CheckIsNil applies the IsNil predicate on the "check" field.
func CheckIsNil() predicate.Policies {
	return predicate.Policies(sql.FieldIsNull(FieldCheck))
}

// CheckNotNil applies the NotNil predicate on the "check" field.
func CheckNotNil() predicate.Policies {
	return predicate.Policies(sql.FieldNotNull(FieldCheck))
}

// HasPolicyLabels applies the HasEdge predicate on the "policy_labels" edge.
func HasPolicyLabels() predicate.Policies {
	return predicate.Policies(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PolicyLabelsTable, PolicyLabelsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPolicyLabelsWith applies the HasEdge predicate on the "policy_labels" edge with a given conditions (other predicates).
func HasPolicyLabelsWith(preds ...predicate.PolicyLabels) predicate.Policies {
	return predicate.Policies(func(s *sql.Selector) {
		step := newPolicyLabelsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasScans applies the HasEdge predicate on the "scans" edge.
func HasScans() predicate.Policies {
	return predicate.Policies(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ScansTable, ScansColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasScansWith applies the HasEdge predicate on the "scans" edge with a given conditions (other predicates).
func HasScansWith(preds ...predicate.Scans) predicate.Policies {
	return predicate.Policies(func(s *sql.Selector) {
		step := newScansStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Policies) predicate.Policies {
	return predicate.Policies(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Policies) predicate.Policies {
	return predicate.Policies(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Policies) predicate.Policies {
	return predicate.Policies(sql.NotPredicates(p))
}
