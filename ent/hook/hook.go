// Code generated by ent, DO NOT EDIT.

package hook

import (
	"context"
	"fmt"

	"github.com/shinobistack/gokakashi/ent"
)

// The AgentLabelsFunc type is an adapter to allow the use of ordinary
// function as AgentLabels mutator.
type AgentLabelsFunc func(context.Context, *ent.AgentLabelsMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f AgentLabelsFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.AgentLabelsMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.AgentLabelsMutation", m)
}

// The AgentTasksFunc type is an adapter to allow the use of ordinary
// function as AgentTasks mutator.
type AgentTasksFunc func(context.Context, *ent.AgentTasksMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f AgentTasksFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.AgentTasksMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.AgentTasksMutation", m)
}

// The AgentsFunc type is an adapter to allow the use of ordinary
// function as Agents mutator.
type AgentsFunc func(context.Context, *ent.AgentsMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f AgentsFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.AgentsMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.AgentsMutation", m)
}

// The IntegrationTypeFunc type is an adapter to allow the use of ordinary
// function as IntegrationType mutator.
type IntegrationTypeFunc func(context.Context, *ent.IntegrationTypeMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f IntegrationTypeFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.IntegrationTypeMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.IntegrationTypeMutation", m)
}

// The IntegrationsFunc type is an adapter to allow the use of ordinary
// function as Integrations mutator.
type IntegrationsFunc func(context.Context, *ent.IntegrationsMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f IntegrationsFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.IntegrationsMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.IntegrationsMutation", m)
}

// The PoliciesFunc type is an adapter to allow the use of ordinary
// function as Policies mutator.
type PoliciesFunc func(context.Context, *ent.PoliciesMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f PoliciesFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.PoliciesMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.PoliciesMutation", m)
}

// The PolicyLabelsFunc type is an adapter to allow the use of ordinary
// function as PolicyLabels mutator.
type PolicyLabelsFunc func(context.Context, *ent.PolicyLabelsMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f PolicyLabelsFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.PolicyLabelsMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.PolicyLabelsMutation", m)
}

// The ScanLabelsFunc type is an adapter to allow the use of ordinary
// function as ScanLabels mutator.
type ScanLabelsFunc func(context.Context, *ent.ScanLabelsMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ScanLabelsFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.ScanLabelsMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ScanLabelsMutation", m)
}

// The ScanNotifyFunc type is an adapter to allow the use of ordinary
// function as ScanNotify mutator.
type ScanNotifyFunc func(context.Context, *ent.ScanNotifyMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ScanNotifyFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.ScanNotifyMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ScanNotifyMutation", m)
}

// The ScansFunc type is an adapter to allow the use of ordinary
// function as Scans mutator.
type ScansFunc func(context.Context, *ent.ScansMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ScansFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.ScansMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ScansMutation", m)
}

// The V2AgentsFunc type is an adapter to allow the use of ordinary
// function as V2Agents mutator.
type V2AgentsFunc func(context.Context, *ent.V2AgentsMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f V2AgentsFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.V2AgentsMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.V2AgentsMutation", m)
}

// Condition is a hook condition function.
type Condition func(context.Context, ent.Mutation) bool

// And groups conditions with the AND operator.
func And(first, second Condition, rest ...Condition) Condition {
	return func(ctx context.Context, m ent.Mutation) bool {
		if !first(ctx, m) || !second(ctx, m) {
			return false
		}
		for _, cond := range rest {
			if !cond(ctx, m) {
				return false
			}
		}
		return true
	}
}

// Or groups conditions with the OR operator.
func Or(first, second Condition, rest ...Condition) Condition {
	return func(ctx context.Context, m ent.Mutation) bool {
		if first(ctx, m) || second(ctx, m) {
			return true
		}
		for _, cond := range rest {
			if cond(ctx, m) {
				return true
			}
		}
		return false
	}
}

// Not negates a given condition.
func Not(cond Condition) Condition {
	return func(ctx context.Context, m ent.Mutation) bool {
		return !cond(ctx, m)
	}
}

// HasOp is a condition testing mutation operation.
func HasOp(op ent.Op) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		return m.Op().Is(op)
	}
}

// HasAddedFields is a condition validating `.AddedField` on fields.
func HasAddedFields(field string, fields ...string) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		if _, exists := m.AddedField(field); !exists {
			return false
		}
		for _, field := range fields {
			if _, exists := m.AddedField(field); !exists {
				return false
			}
		}
		return true
	}
}

// HasClearedFields is a condition validating `.FieldCleared` on fields.
func HasClearedFields(field string, fields ...string) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		if exists := m.FieldCleared(field); !exists {
			return false
		}
		for _, field := range fields {
			if exists := m.FieldCleared(field); !exists {
				return false
			}
		}
		return true
	}
}

// HasFields is a condition validating `.Field` on fields.
func HasFields(field string, fields ...string) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		if _, exists := m.Field(field); !exists {
			return false
		}
		for _, field := range fields {
			if _, exists := m.Field(field); !exists {
				return false
			}
		}
		return true
	}
}

// If executes the given hook under condition.
//
//	hook.If(ComputeAverage, And(HasFields(...), HasAddedFields(...)))
func If(hk ent.Hook, cond Condition) ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if cond(ctx, m) {
				return hk(next).Mutate(ctx, m)
			}
			return next.Mutate(ctx, m)
		})
	}
}

// On executes the given hook only for the given operation.
//
//	hook.On(Log, ent.Delete|ent.Create)
func On(hk ent.Hook, op ent.Op) ent.Hook {
	return If(hk, HasOp(op))
}

// Unless skips the given hook only for the given operation.
//
//	hook.Unless(Log, ent.Update|ent.UpdateOne)
func Unless(hk ent.Hook, op ent.Op) ent.Hook {
	return If(hk, Not(HasOp(op)))
}

// FixedError is a hook returning a fixed error.
func FixedError(err error) ent.Hook {
	return func(ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(context.Context, ent.Mutation) (ent.Value, error) {
			return nil, err
		})
	}
}

// Reject returns a hook that rejects all operations that match op.
//
//	func (T) Hooks() []ent.Hook {
//		return []ent.Hook{
//			Reject(ent.Delete|ent.Update),
//		}
//	}
func Reject(op ent.Op) ent.Hook {
	hk := FixedError(fmt.Errorf("%s operation is not allowed", op))
	return On(hk, op)
}

// Chain acts as a list of hooks and is effectively immutable.
// Once created, it will always hold the same set of hooks in the same order.
type Chain struct {
	hooks []ent.Hook
}

// NewChain creates a new chain of hooks.
func NewChain(hooks ...ent.Hook) Chain {
	return Chain{append([]ent.Hook(nil), hooks...)}
}

// Hook chains the list of hooks and returns the final hook.
func (c Chain) Hook() ent.Hook {
	return func(mutator ent.Mutator) ent.Mutator {
		for i := len(c.hooks) - 1; i >= 0; i-- {
			mutator = c.hooks[i](mutator)
		}
		return mutator
	}
}

// Append extends a chain, adding the specified hook
// as the last ones in the mutation flow.
func (c Chain) Append(hooks ...ent.Hook) Chain {
	newHooks := make([]ent.Hook, 0, len(c.hooks)+len(hooks))
	newHooks = append(newHooks, c.hooks...)
	newHooks = append(newHooks, hooks...)
	return Chain{newHooks}
}

// Extend extends a chain, adding the specified chain
// as the last ones in the mutation flow.
func (c Chain) Extend(chain Chain) Chain {
	return c.Append(chain.hooks...)
}
