package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/ent/schema/index"
)

// Tokens holds the schema definition for the Tokens entity.
type Tokens struct {
	ent.Schema
}

// Fields of the Tokens.
func (Tokens) Fields() []ent.Field {
	return []ent.Field{
		field.String("provider"),
		field.String("user"),
		field.String("token"),
	}
}

// Edges of the Tokens.
func (Tokens) Edges() []ent.Edge {
	return nil
}

func (Tokens) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("provider", "user").Unique(),
	}
}
