package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Entries holds the schema definition for the Entries entity.
type Entries struct {
	ent.Schema
}

// Fields of the Entries.
func (Entries) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
		field.String("description"),
		field.String("url"),
	}
}

// Edges of the Entries.
func (Entries) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("feeds", Feeds.Type).
			Ref("entries").
			Unique(),
	}
}

func (Entries) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("url").
			Edges("feeds").
			Unique(),
	}
}
