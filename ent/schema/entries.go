package schema

import "entgo.io/ent"

// Entries holds the schema definition for the Entries entity.
type Entries struct {
	ent.Schema
}

// Fields of the Entries.
func (Entries) Fields() []ent.Field {
	return nil
}

// Edges of the Entries.
func (Entries) Edges() []ent.Edge {
	return nil
}
