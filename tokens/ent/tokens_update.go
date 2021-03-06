// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/tvastar/goat/tokens/ent/predicate"
	"github.com/tvastar/goat/tokens/ent/tokens"
)

// TokensUpdate is the builder for updating Tokens entities.
type TokensUpdate struct {
	config
	hooks      []Hook
	mutation   *TokensMutation
	predicates []predicate.Tokens
}

// Where adds a new predicate for the builder.
func (tu *TokensUpdate) Where(ps ...predicate.Tokens) *TokensUpdate {
	tu.predicates = append(tu.predicates, ps...)
	return tu
}

// SetProvider sets the provider field.
func (tu *TokensUpdate) SetProvider(s string) *TokensUpdate {
	tu.mutation.SetProvider(s)
	return tu
}

// SetUser sets the user field.
func (tu *TokensUpdate) SetUser(s string) *TokensUpdate {
	tu.mutation.SetUser(s)
	return tu
}

// SetToken sets the token field.
func (tu *TokensUpdate) SetToken(s string) *TokensUpdate {
	tu.mutation.SetToken(s)
	return tu
}

// Mutation returns the TokensMutation object of the builder.
func (tu *TokensUpdate) Mutation() *TokensMutation {
	return tu.mutation
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (tu *TokensUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(tu.hooks) == 0 {
		affected, err = tu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*TokensMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			tu.mutation = mutation
			affected, err = tu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(tu.hooks) - 1; i >= 0; i-- {
			mut = tu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, tu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (tu *TokensUpdate) SaveX(ctx context.Context) int {
	affected, err := tu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tu *TokensUpdate) Exec(ctx context.Context) error {
	_, err := tu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tu *TokensUpdate) ExecX(ctx context.Context) {
	if err := tu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tu *TokensUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   tokens.Table,
			Columns: tokens.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: tokens.FieldID,
			},
		},
	}
	if ps := tu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tu.mutation.Provider(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: tokens.FieldProvider,
		})
	}
	if value, ok := tu.mutation.User(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: tokens.FieldUser,
		})
	}
	if value, ok := tu.mutation.Token(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: tokens.FieldToken,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, tu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{tokens.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// TokensUpdateOne is the builder for updating a single Tokens entity.
type TokensUpdateOne struct {
	config
	hooks    []Hook
	mutation *TokensMutation
}

// SetProvider sets the provider field.
func (tuo *TokensUpdateOne) SetProvider(s string) *TokensUpdateOne {
	tuo.mutation.SetProvider(s)
	return tuo
}

// SetUser sets the user field.
func (tuo *TokensUpdateOne) SetUser(s string) *TokensUpdateOne {
	tuo.mutation.SetUser(s)
	return tuo
}

// SetToken sets the token field.
func (tuo *TokensUpdateOne) SetToken(s string) *TokensUpdateOne {
	tuo.mutation.SetToken(s)
	return tuo
}

// Mutation returns the TokensMutation object of the builder.
func (tuo *TokensUpdateOne) Mutation() *TokensMutation {
	return tuo.mutation
}

// Save executes the query and returns the updated entity.
func (tuo *TokensUpdateOne) Save(ctx context.Context) (*Tokens, error) {
	var (
		err  error
		node *Tokens
	)
	if len(tuo.hooks) == 0 {
		node, err = tuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*TokensMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			tuo.mutation = mutation
			node, err = tuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(tuo.hooks) - 1; i >= 0; i-- {
			mut = tuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, tuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (tuo *TokensUpdateOne) SaveX(ctx context.Context) *Tokens {
	t, err := tuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return t
}

// Exec executes the query on the entity.
func (tuo *TokensUpdateOne) Exec(ctx context.Context) error {
	_, err := tuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tuo *TokensUpdateOne) ExecX(ctx context.Context) {
	if err := tuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tuo *TokensUpdateOne) sqlSave(ctx context.Context) (t *Tokens, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   tokens.Table,
			Columns: tokens.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: tokens.FieldID,
			},
		},
	}
	id, ok := tuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing Tokens.ID for update")}
	}
	_spec.Node.ID.Value = id
	if value, ok := tuo.mutation.Provider(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: tokens.FieldProvider,
		})
	}
	if value, ok := tuo.mutation.User(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: tokens.FieldUser,
		})
	}
	if value, ok := tuo.mutation.Token(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: tokens.FieldToken,
		})
	}
	t = &Tokens{config: tuo.config}
	_spec.Assign = t.assignValues
	_spec.ScanValues = t.scanValues()
	if err = sqlgraph.UpdateNode(ctx, tuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{tokens.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return t, nil
}
