// Package ctxlog is a simple wrapper around `context.Context` for storing
// and retrieving `logrus.Entry`. It allows more consistent logging of per-request
// fields, without needing to pass each of them explicitly every time.
//
// Example usage:
//
//     ctx = ctxlog.WithField(ctx, "user", user)
//     // ...
//     ctxlog.With(ctx).Printf("User did something")  // resulting entry will have field "user"
//
// You can also use it to trace request handling across function boundaries:
//
//     func handleHTTP(w http.ResponseWriter, r *http.Request) {
//          ctx := ctxlog.WithField(context.Background(), "request_id", genRequestID())
//          doStuff(ctx)  // if doStuff uses ctxlog too, all entries will have "request_id"
//          // ...
package ctxlog

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type ctxkey struct{}

// With returns `logrus.Entry` from `ctx`. If `ctx` doesn't have any entry yet,
// it will return a new entry (by calling `logrus.NewEntry(logrus.StandardLogger())`).
func With(ctx context.Context) *log.Entry {
	e, ok := ctx.Value(ctxkey{}).(*log.Entry)
	if !ok {
		return log.NewEntry(log.StandardLogger())
	}
	return e
}

// WithField adds a single field to entry stored in `ctx`.
func WithField(ctx context.Context, name string, value interface{}) context.Context {
	return Set(ctx, With(ctx).WithField(name, value))
}

// WithFields adds `fields` to entry stored in `ctx`.
func WithFields(ctx context.Context, fields log.Fields) context.Context {
	return Set(ctx, With(ctx).WithFields(fields))
}

// Fields returns fields currently stored in `ctx`.
func Fields(ctx context.Context) log.Fields {
	return With(ctx).Data
}

// Clear unsets entry stored in `ctx`.
//
// Disclaimer: due to the way `context` package is implemented, it's not possible
// to actually remove the entry, i.e., it will not be garbage collected.
// But context value returned by this function will cause `With` to return a
// fresh entry.
func Clear(ctx context.Context) context.Context {
	// Using Set would set the value to nil of type *log.Entry instead of untyped nil.
	return context.WithValue(ctx, ctxkey{}, nil)
}

// Set associates `logrus.Entry` with a context. Subsequent calls to `With` will
// return this entry.
//
// Normally just calling `WithField` or `WithFields` is enough, but you can use
// this function in case you need to have some other aspects of stored
// `logrus.Entry` changed.
func Set(ctx context.Context, entry *log.Entry) context.Context {
	return context.WithValue(ctx, ctxkey{}, entry)
}
