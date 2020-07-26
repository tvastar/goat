// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"sync"

	"github.com/tvastar/goat/tokens/ent/tokens"

	"github.com/facebookincubator/ent"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeTokens = "Tokens"
)

// TokensMutation represents an operation that mutate the TokensSlice
// nodes in the graph.
type TokensMutation struct {
	config
	op            Op
	typ           string
	id            *int
	provider      *string
	user          *string
	token         *string
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*Tokens, error)
}

var _ ent.Mutation = (*TokensMutation)(nil)

// tokensOption allows to manage the mutation configuration using functional options.
type tokensOption func(*TokensMutation)

// newTokensMutation creates new mutation for $n.Name.
func newTokensMutation(c config, op Op, opts ...tokensOption) *TokensMutation {
	m := &TokensMutation{
		config:        c,
		op:            op,
		typ:           TypeTokens,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withTokensID sets the id field of the mutation.
func withTokensID(id int) tokensOption {
	return func(m *TokensMutation) {
		var (
			err   error
			once  sync.Once
			value *Tokens
		)
		m.oldValue = func(ctx context.Context) (*Tokens, error) {
			once.Do(func() {
				if m.done {
					err = fmt.Errorf("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Tokens.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withTokens sets the old Tokens of the mutation.
func withTokens(node *Tokens) tokensOption {
	return func(m *TokensMutation) {
		m.oldValue = func(context.Context) (*Tokens, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m TokensMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m TokensMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *TokensMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetProvider sets the provider field.
func (m *TokensMutation) SetProvider(s string) {
	m.provider = &s
}

// Provider returns the provider value in the mutation.
func (m *TokensMutation) Provider() (r string, exists bool) {
	v := m.provider
	if v == nil {
		return
	}
	return *v, true
}

// OldProvider returns the old provider value of the Tokens.
// If the Tokens object wasn't provided to the builder, the object is fetched
// from the database.
// An error is returned if the mutation operation is not UpdateOne, or database query fails.
func (m *TokensMutation) OldProvider(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldProvider is allowed only on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldProvider requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldProvider: %w", err)
	}
	return oldValue.Provider, nil
}

// ResetProvider reset all changes of the "provider" field.
func (m *TokensMutation) ResetProvider() {
	m.provider = nil
}

// SetUser sets the user field.
func (m *TokensMutation) SetUser(s string) {
	m.user = &s
}

// User returns the user value in the mutation.
func (m *TokensMutation) User() (r string, exists bool) {
	v := m.user
	if v == nil {
		return
	}
	return *v, true
}

// OldUser returns the old user value of the Tokens.
// If the Tokens object wasn't provided to the builder, the object is fetched
// from the database.
// An error is returned if the mutation operation is not UpdateOne, or database query fails.
func (m *TokensMutation) OldUser(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldUser is allowed only on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldUser requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUser: %w", err)
	}
	return oldValue.User, nil
}

// ResetUser reset all changes of the "user" field.
func (m *TokensMutation) ResetUser() {
	m.user = nil
}

// SetToken sets the token field.
func (m *TokensMutation) SetToken(s string) {
	m.token = &s
}

// Token returns the token value in the mutation.
func (m *TokensMutation) Token() (r string, exists bool) {
	v := m.token
	if v == nil {
		return
	}
	return *v, true
}

// OldToken returns the old token value of the Tokens.
// If the Tokens object wasn't provided to the builder, the object is fetched
// from the database.
// An error is returned if the mutation operation is not UpdateOne, or database query fails.
func (m *TokensMutation) OldToken(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldToken is allowed only on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldToken requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldToken: %w", err)
	}
	return oldValue.Token, nil
}

// ResetToken reset all changes of the "token" field.
func (m *TokensMutation) ResetToken() {
	m.token = nil
}

// Op returns the operation name.
func (m *TokensMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Tokens).
func (m *TokensMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *TokensMutation) Fields() []string {
	fields := make([]string, 0, 3)
	if m.provider != nil {
		fields = append(fields, tokens.FieldProvider)
	}
	if m.user != nil {
		fields = append(fields, tokens.FieldUser)
	}
	if m.token != nil {
		fields = append(fields, tokens.FieldToken)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *TokensMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case tokens.FieldProvider:
		return m.Provider()
	case tokens.FieldUser:
		return m.User()
	case tokens.FieldToken:
		return m.Token()
	}
	return nil, false
}

// OldField returns the old value of the field from the database.
// An error is returned if the mutation operation is not UpdateOne,
// or the query to the database was failed.
func (m *TokensMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case tokens.FieldProvider:
		return m.OldProvider(ctx)
	case tokens.FieldUser:
		return m.OldUser(ctx)
	case tokens.FieldToken:
		return m.OldToken(ctx)
	}
	return nil, fmt.Errorf("unknown Tokens field %s", name)
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *TokensMutation) SetField(name string, value ent.Value) error {
	switch name {
	case tokens.FieldProvider:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetProvider(v)
		return nil
	case tokens.FieldUser:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUser(v)
		return nil
	case tokens.FieldToken:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetToken(v)
		return nil
	}
	return fmt.Errorf("unknown Tokens field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *TokensMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *TokensMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *TokensMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Tokens numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *TokensMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *TokensMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *TokensMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Tokens nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *TokensMutation) ResetField(name string) error {
	switch name {
	case tokens.FieldProvider:
		m.ResetProvider()
		return nil
	case tokens.FieldUser:
		m.ResetUser()
		return nil
	case tokens.FieldToken:
		m.ResetToken()
		return nil
	}
	return fmt.Errorf("unknown Tokens field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *TokensMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *TokensMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *TokensMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *TokensMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *TokensMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *TokensMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *TokensMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Tokens unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *TokensMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Tokens edge %s", name)
}
