// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/yseto/podcaster/ent/feeds"
	"github.com/yseto/podcaster/ent/predicate"
	"github.com/yseto/podcaster/ent/users"
)

// UsersUpdate is the builder for updating Users entities.
type UsersUpdate struct {
	config
	hooks     []Hook
	mutation  *UsersMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the UsersUpdate builder.
func (uu *UsersUpdate) Where(ps ...predicate.Users) *UsersUpdate {
	uu.mutation.Where(ps...)
	return uu
}

// SetName sets the "name" field.
func (uu *UsersUpdate) SetName(s string) *UsersUpdate {
	uu.mutation.SetName(s)
	return uu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (uu *UsersUpdate) SetNillableName(s *string) *UsersUpdate {
	if s != nil {
		uu.SetName(*s)
	}
	return uu
}

// SetPassword sets the "password" field.
func (uu *UsersUpdate) SetPassword(s string) *UsersUpdate {
	uu.mutation.SetPassword(s)
	return uu
}

// SetNillablePassword sets the "password" field if the given value is not nil.
func (uu *UsersUpdate) SetNillablePassword(s *string) *UsersUpdate {
	if s != nil {
		uu.SetPassword(*s)
	}
	return uu
}

// AddFeedIDs adds the "feeds" edge to the Feeds entity by IDs.
func (uu *UsersUpdate) AddFeedIDs(ids ...int) *UsersUpdate {
	uu.mutation.AddFeedIDs(ids...)
	return uu
}

// AddFeeds adds the "feeds" edges to the Feeds entity.
func (uu *UsersUpdate) AddFeeds(f ...*Feeds) *UsersUpdate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return uu.AddFeedIDs(ids...)
}

// Mutation returns the UsersMutation object of the builder.
func (uu *UsersUpdate) Mutation() *UsersMutation {
	return uu.mutation
}

// ClearFeeds clears all "feeds" edges to the Feeds entity.
func (uu *UsersUpdate) ClearFeeds() *UsersUpdate {
	uu.mutation.ClearFeeds()
	return uu
}

// RemoveFeedIDs removes the "feeds" edge to Feeds entities by IDs.
func (uu *UsersUpdate) RemoveFeedIDs(ids ...int) *UsersUpdate {
	uu.mutation.RemoveFeedIDs(ids...)
	return uu
}

// RemoveFeeds removes "feeds" edges to Feeds entities.
func (uu *UsersUpdate) RemoveFeeds(f ...*Feeds) *UsersUpdate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return uu.RemoveFeedIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (uu *UsersUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, uu.sqlSave, uu.mutation, uu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (uu *UsersUpdate) SaveX(ctx context.Context) int {
	affected, err := uu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (uu *UsersUpdate) Exec(ctx context.Context) error {
	_, err := uu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uu *UsersUpdate) ExecX(ctx context.Context) {
	if err := uu.Exec(ctx); err != nil {
		panic(err)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (uu *UsersUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *UsersUpdate {
	uu.modifiers = append(uu.modifiers, modifiers...)
	return uu
}

func (uu *UsersUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(users.Table, users.Columns, sqlgraph.NewFieldSpec(users.FieldID, field.TypeInt))
	if ps := uu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uu.mutation.Name(); ok {
		_spec.SetField(users.FieldName, field.TypeString, value)
	}
	if value, ok := uu.mutation.Password(); ok {
		_spec.SetField(users.FieldPassword, field.TypeString, value)
	}
	if uu.mutation.FeedsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   users.FeedsTable,
			Columns: []string{users.FeedsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(feeds.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.RemovedFeedsIDs(); len(nodes) > 0 && !uu.mutation.FeedsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   users.FeedsTable,
			Columns: []string{users.FeedsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(feeds.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.FeedsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   users.FeedsTable,
			Columns: []string{users.FeedsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(feeds.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(uu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, uu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{users.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	uu.mutation.done = true
	return n, nil
}

// UsersUpdateOne is the builder for updating a single Users entity.
type UsersUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *UsersMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetName sets the "name" field.
func (uuo *UsersUpdateOne) SetName(s string) *UsersUpdateOne {
	uuo.mutation.SetName(s)
	return uuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (uuo *UsersUpdateOne) SetNillableName(s *string) *UsersUpdateOne {
	if s != nil {
		uuo.SetName(*s)
	}
	return uuo
}

// SetPassword sets the "password" field.
func (uuo *UsersUpdateOne) SetPassword(s string) *UsersUpdateOne {
	uuo.mutation.SetPassword(s)
	return uuo
}

// SetNillablePassword sets the "password" field if the given value is not nil.
func (uuo *UsersUpdateOne) SetNillablePassword(s *string) *UsersUpdateOne {
	if s != nil {
		uuo.SetPassword(*s)
	}
	return uuo
}

// AddFeedIDs adds the "feeds" edge to the Feeds entity by IDs.
func (uuo *UsersUpdateOne) AddFeedIDs(ids ...int) *UsersUpdateOne {
	uuo.mutation.AddFeedIDs(ids...)
	return uuo
}

// AddFeeds adds the "feeds" edges to the Feeds entity.
func (uuo *UsersUpdateOne) AddFeeds(f ...*Feeds) *UsersUpdateOne {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return uuo.AddFeedIDs(ids...)
}

// Mutation returns the UsersMutation object of the builder.
func (uuo *UsersUpdateOne) Mutation() *UsersMutation {
	return uuo.mutation
}

// ClearFeeds clears all "feeds" edges to the Feeds entity.
func (uuo *UsersUpdateOne) ClearFeeds() *UsersUpdateOne {
	uuo.mutation.ClearFeeds()
	return uuo
}

// RemoveFeedIDs removes the "feeds" edge to Feeds entities by IDs.
func (uuo *UsersUpdateOne) RemoveFeedIDs(ids ...int) *UsersUpdateOne {
	uuo.mutation.RemoveFeedIDs(ids...)
	return uuo
}

// RemoveFeeds removes "feeds" edges to Feeds entities.
func (uuo *UsersUpdateOne) RemoveFeeds(f ...*Feeds) *UsersUpdateOne {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return uuo.RemoveFeedIDs(ids...)
}

// Where appends a list predicates to the UsersUpdate builder.
func (uuo *UsersUpdateOne) Where(ps ...predicate.Users) *UsersUpdateOne {
	uuo.mutation.Where(ps...)
	return uuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (uuo *UsersUpdateOne) Select(field string, fields ...string) *UsersUpdateOne {
	uuo.fields = append([]string{field}, fields...)
	return uuo
}

// Save executes the query and returns the updated Users entity.
func (uuo *UsersUpdateOne) Save(ctx context.Context) (*Users, error) {
	return withHooks(ctx, uuo.sqlSave, uuo.mutation, uuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (uuo *UsersUpdateOne) SaveX(ctx context.Context) *Users {
	node, err := uuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (uuo *UsersUpdateOne) Exec(ctx context.Context) error {
	_, err := uuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uuo *UsersUpdateOne) ExecX(ctx context.Context) {
	if err := uuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (uuo *UsersUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *UsersUpdateOne {
	uuo.modifiers = append(uuo.modifiers, modifiers...)
	return uuo
}

func (uuo *UsersUpdateOne) sqlSave(ctx context.Context) (_node *Users, err error) {
	_spec := sqlgraph.NewUpdateSpec(users.Table, users.Columns, sqlgraph.NewFieldSpec(users.FieldID, field.TypeInt))
	id, ok := uuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Users.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := uuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, users.FieldID)
		for _, f := range fields {
			if !users.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != users.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := uuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uuo.mutation.Name(); ok {
		_spec.SetField(users.FieldName, field.TypeString, value)
	}
	if value, ok := uuo.mutation.Password(); ok {
		_spec.SetField(users.FieldPassword, field.TypeString, value)
	}
	if uuo.mutation.FeedsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   users.FeedsTable,
			Columns: []string{users.FeedsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(feeds.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.RemovedFeedsIDs(); len(nodes) > 0 && !uuo.mutation.FeedsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   users.FeedsTable,
			Columns: []string{users.FeedsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(feeds.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.FeedsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   users.FeedsTable,
			Columns: []string{users.FeedsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(feeds.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(uuo.modifiers...)
	_node = &Users{config: uuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, uuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{users.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	uuo.mutation.done = true
	return _node, nil
}
