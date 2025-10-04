# envvar

Tiny-but-mighty environment reader and struct binder for Go.

* Zero deps, stdlib only.
* Typed getters: bool, int, float64, duration, URL, IP, slices.
* Struct binding with tags, defaults, and JSON decode.
* `${VAR}` and `${VAR:-def}` expansion.
* Pluggable sources (env + maps + composites).
* Lazy getters and safe redacted dumps.
* Optional hooks for metrics/tracing without adding deps.

## Install

```bash
go get github.com/aatuh/envvar/v2
```

## Quick start

Check the examples from the [examples package]

Run examples with:

```bash
go test -v -count 1 ./examples

# Run specific example.
go test -v -count 1 ./examples -run TestBasicGetters
```

## Features

### Typed getters

* `Get`, `GetOr`, `MustGet`
* `GetBool`, `GetInt`, `GetFloat64`, `GetDuration`
* `GetURL`, `GetIP`, `GetStringSlice` (+ `GetStringSliceSep`)
* Generic: `GetTyped[T](key, conv)`
* All have `Must*` and `Or` variants where it makes sense.

### Expansion

* `${NAME}` and `${NAME:-default}` are expanded in values read from
  env and when using `ExpandMap`.

### Lazy getters

Cache-on-first-use helpers, e.g. `LazyBool("DEBUG")()`.

### Struct binding

Populate a struct from environment with tags:

* `env:"NAME[,required]"` choose the env var name and requiredness.
* `envdef:"value"` default used if missing.
* `envsep:","` separator for `[]string` (default ",").
* `envjson:"true"` JSON decode into field type (maps, slices, structs).

Pointer fields are allocated automatically.

#### Prefix binding

Try a prefixed variable first, then fall back to the base name:

```go
// Tries MYAPP_PORT first, then PORT.
envvar.MustBindWithPrefix(&cfg, "MYAPP_")
```

### Environment variable sources

By default, getters and `Bind` read from process environment variables.
You can also load from files using `LoadEnvVars`:

```go
// Load from specific files
if err := envvar.LoadEnvVars([]string{"./.env.local", "./.env"}); err != nil {
  // handle error
}

// Or use the default paths (.env, then /env/.env)
envvar.MustLoadEnvVars(nil)

port := envvar.MustGetInt("PORT") // reads from loaded env vars
```

### Map expansion helper

Expand `${VAR}` and `${VAR:-def}` inside a map, using map values first,
then process environment variables:

```go
in := map[string]string{
  "HOST": "db.local",
  "DSN":  "postgres://${HOST}:${PORT:-5432}/app",
}
out := envvar.MustExpandMap(in)
// out["DSN"] == "postgres://db.local:5432/app"
```

### .env loading

```go
// Try custom locations in order:
if err := envvar.LoadEnvVars([]string{"./.env.local", "./.env"}); err != nil {
  // handle error
}
```

### Redacted dump

```go
log.Printf("env: %#v", envvar.DumpRedacted())
```

Heuristics redact keys containing `SECRET`, `TOKEN`, `PASSWORD`, or
suffix `_KEY`.

### Error handling

* All binding errors are aggregated and returned as a `MultiError`.
* Missing required fields are reported clearly.

### Observability hooks (optional)

```go
type myHook struct{}
func (myHook) OnLoad(src string, keys int) {}
func (myHook) OnGet(k string, ok bool, err error, d time.Duration) {}

envvar.SetHook(myHook{})
```
