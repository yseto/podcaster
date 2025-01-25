package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Feeds holds the schema definition for the Feeds entity.
type Feeds struct {
	ent.Schema
}

// Fields of the Feeds.
func (Feeds) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
		field.String("url"),
	}
}

// Edges of the Feeds.
func (Feeds) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("entries", Entries.Type),
	}
}
