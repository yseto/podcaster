package server

import (
	"context"

	myContext "github.com/yseto/podcaster/server/context"
)

type myContextKey struct{}

const (
	myUserKey = "userid"
)

func emptyUserContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, myContextKey{}, myContext.New())
}

func newUserContext(ctx context.Context, userid uint64) {
	ctx.Value(myContextKey{}).(*myContext.Context).Set(myUserKey, userid)
}

func userFromContext(ctx context.Context) uint64 {
	return ctx.Value(myContextKey{}).(*myContext.Context).Get(myUserKey).(uint64)
}
