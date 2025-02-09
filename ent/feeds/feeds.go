// Code generated by ent, DO NOT EDIT.

package feeds

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the feeds type in the database.
	Label = "feeds"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldTitle holds the string denoting the title field in the database.
	FieldTitle = "title"
	// FieldURL holds the string denoting the url field in the database.
	FieldURL = "url"
	// EdgeEntries holds the string denoting the entries edge name in mutations.
	EdgeEntries = "entries"
	// Table holds the table name of the feeds in the database.
	Table = "feeds"
	// EntriesTable is the table that holds the entries relation/edge.
	EntriesTable = "entries"
	// EntriesInverseTable is the table name for the Entries entity.
	// It exists in this package in order to avoid circular dependency with the "entries" package.
	EntriesInverseTable = "entries"
	// EntriesColumn is the table column denoting the entries relation/edge.
	EntriesColumn = "feeds_entries"
)

// Columns holds all SQL columns for feeds fields.
var Columns = []string{
	FieldID,
	FieldTitle,
	FieldURL,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "feeds"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"users_feeds",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the Feeds queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByTitle orders the results by the title field.
func ByTitle(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTitle, opts...).ToFunc()
}

// ByURL orders the results by the url field.
func ByURL(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldURL, opts...).ToFunc()
}

// ByEntriesCount orders the results by entries count.
func ByEntriesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newEntriesStep(), opts...)
	}
}

// ByEntries orders the results by entries terms.
func ByEntries(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newEntriesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newEntriesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(EntriesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, EntriesTable, EntriesColumn),
	)
}
