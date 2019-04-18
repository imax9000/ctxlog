# Wrappers for using logrus with context.Context

[![GoDoc](https://godoc.org/github.com/gelraen/ctxlog?status.svg)](https://godoc.org/github.com/gelraen/ctxlog)

Example usage:

```go
ctx = ctxlog.WithField(ctx, "user", user)
// ...
ctxlog.With(ctx).Printf("User did something")  // resulting entry will have field "user"
```

You can also use it to trace request handling across function boundaries:

```go
func handleHTTP(w http.ResponseWriter, r *http.Request) {
    ctx := ctxlog.WithField(context.Background(), "request_id", genRequestID())
    doStuff(ctx)  // if doStuff uses ctxlog too, all entries will have "request_id"
    // ...
```
