// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/tvastar/goat/tokens/ent/tokens"
)

// Tokens is the model entity for the Tokens schema.
type Tokens struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Provider holds the value of the "provider" field.
	Provider string `json:"provider,omitempty"`
	// User holds the value of the "user" field.
	User string `json:"user,omitempty"`
	// Token holds the value of the "token" field.
	Token string `json:"token,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Tokens) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullString{}, // provider
		&sql.NullString{}, // user
		&sql.NullString{}, // token
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Tokens fields.
func (t *Tokens) assignValues(values ...interface{}) error {
	if m, n := len(values), len(tokens.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	t.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field provider", values[0])
	} else if value.Valid {
		t.Provider = value.String
	}
	if value, ok := values[1].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field user", values[1])
	} else if value.Valid {
		t.User = value.String
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field token", values[2])
	} else if value.Valid {
		t.Token = value.String
	}
	return nil
}

// Update returns a builder for updating this Tokens.
// Note that, you need to call Tokens.Unwrap() before calling this method, if this Tokens
// was returned from a transaction, and the transaction was committed or rolled back.
func (t *Tokens) Update() *TokensUpdateOne {
	return (&TokensClient{config: t.config}).UpdateOne(t)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (t *Tokens) Unwrap() *Tokens {
	tx, ok := t.config.driver.(*txDriver)
	if !ok {
		panic("ent: Tokens is not a transactional entity")
	}
	t.config.driver = tx.drv
	return t
}

// String implements the fmt.Stringer.
func (t *Tokens) String() string {
	var builder strings.Builder
	builder.WriteString("Tokens(")
	builder.WriteString(fmt.Sprintf("id=%v", t.ID))
	builder.WriteString(", provider=")
	builder.WriteString(t.Provider)
	builder.WriteString(", user=")
	builder.WriteString(t.User)
	builder.WriteString(", token=")
	builder.WriteString(t.Token)
	builder.WriteByte(')')
	return builder.String()
}

// TokensSlice is a parsable slice of Tokens.
type TokensSlice []*Tokens

func (t TokensSlice) config(cfg config) {
	for _i := range t {
		t[_i].config = cfg
	}
}