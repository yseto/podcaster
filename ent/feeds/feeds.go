// Code generated by ent, DO NOT EDIT.

package feeds

import (
	"entgo.io/ent/dialect/sql"
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
	// Table holds the table name of the feeds in the database.
	Table = "feeds"
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
