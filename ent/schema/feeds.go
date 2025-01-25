package schema

import "entgo.io/ent"

// Feeds holds the schema definition for the Feeds entity.
type Feeds struct {
	ent.Schema
}

// Fields of the Feeds.
func (Feeds) Fields() []ent.Field {
	return nil
}

// Edges of the Feeds.
func (Feeds) Edges() []ent.Edge {
	return nil
}
