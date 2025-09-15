# envvar

Tiny helpers for loading and reading environment variables with defaults,
typed access, and lazy evaluation.

## Install

```go
import "github.com/aatuh/envvar"
```

## Quick start

`.env`:

```
PORT=8080
DEBUG=true
TAGS=alpha,beta
```

Go:

```go
// Optional: getters auto-load on first use.
envvar.MustLoadEnvVars(nil)

port := envvar.GetOr("PORT", "8000")
debug := envvar.GetBoolOr("DEBUG", false)
tags  := envvar.MustGetStringSlice("TAGS")
```

## Typed getters

```go
p := envvar.MustGetInt("PORT")
host, err := envvar.GetOrErr("HOST")
b, err := envvar.GetBool("DEBUG")
list, err := envvar.GetStringSlice("TAGS")
```

Generic conversion:

```go
toUint16 := func(s string) (uint16, error) {
  v, err := strconv.ParseUint(s, 10, 16)
  return uint16(v), err
}

u16, err := envvar.GetTyped("PORT", toUint16)
```

## Lazy getters

Evaluate once on first call, then reuse the cached value.

```go
getPort  := envvar.LazyInt("PORT")
getDebug := envvar.LazyBool("DEBUG")
getTags  := envvar.LazyStringSlice("TAGS")

toDur := func(s string) (time.Duration, error) { return time.ParseDuration(s) }
getTimeout := envvar.LazyTyped("TIMEOUT", toDur)
```

## Loading behavior

- Looks for the first existing file in: `.env`, `/env/.env`.
- Use custom paths with precedence:

```go
envvar.MustLoadEnvVars([]string{".env.local", ".env"})
```

- File parsing is simple `KEY=VALUE` per line; values are not quoted.
- All loads happen once per process; subsequent calls are no-ops.

## Utilities

```go
m, err := envvar.ReadEnvVarFile(".env")
_ = envvar.SetEnvVars(m)
```

## Notes

- Functions prefixed with `Must*` panic on missing or invalid values.
- `Get*` variants return errors or allow defaults (e.g. `GetOr`).
- CSV lists are split on `,` for `*StringSlice` helpers.
