# envvar

Tiny-but-mighty environment reader and struct binder for Go.

* Zero deps, stdlib only.
* Typed getters: bool, int, float64, duration, URL, IP, slices.
* Struct binding with tags, defaults, validation, and JSON decode.
* `${VAR}` and `${VAR:-def}` expansion.
* Pluggable sources (env + maps + composites).
* Lazy getters and safe redacted dumps.
* Optional hooks for metrics/tracing without adding deps.

## Install

```bash
go get github.com/aatuh/envvar/v2
```

## Quick start

```go
package main

import (
  "log"
  "net/url"
  "time"

  "github.com/aatuh/envvar/v2"
)

type Config struct {
  Port    int           `env:"PORT,required" validate:"min=1,max=65535"`
  Debug   bool          `env:"DEBUG"`
  Timeout time.Duration `env:"TIMEOUT" envdef:"5s"`
  DSN     *url.URL      `env:"DATABASE_URL,required"`
  Tags    []string      `env:"TAGS" envsep:"," validate:"oneof=a|b|c"`
}

func main() {
  // Load .env if present (tries .env, then /env/.env)
  envvar.MustLoadEnvVars(nil)

  // Getters
  port := envvar.GetIntOr("PORT", 8080)
  ttl := envvar.GetDurationOr("CACHE_TTL", 5*time.Second)
  log.Println("port", port, "ttl", ttl)

  // Bind into struct (from active source, defaults to process env)
  var cfg Config
  envvar.MustBind(&cfg)
  log.Printf("cfg: %+v", cfg)
}
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
* `validate:"min=..,max=.."` numeric or string length checks.
* `validate:"oneof=a|b|c"` restrict to allowed values.
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

### Validation rules

* `min` and `max` work for ints, uints, durations, and string length.

  * For `time.Duration`, specify values like `"250ms"`, `"5s"`, etc.
* `oneof` works for strings and `[]string`. The allowed set is always
  pipe-separated: `oneof=a|b|c`.
* All binding errors are aggregated and returned as a `MultiError`.

### Observability hooks (optional)

```go
type myHook struct{}
func (myHook) OnLoad(src string, keys int) {}
func (myHook) OnGet(k string, ok bool, err error, d time.Duration) {}

envvar.SetHook(myHook{})
```

## Examples

See the `/examples` directory for comprehensive, progressive examples:

* `01_basic_getters_test.go` - Basic getter functions
* `02_advanced_getters_test.go` - Advanced getters and lazy evaluation
* `03_struct_binding_test.go` - Struct binding and validation
* `04_file_loading_test.go` - Loading from files
* `05_advanced_features_test.go` - Redaction, hooks, and complex validation
* `06_integration_example_test.go` - Complete application configuration

Run examples with:
```bash
go test ./examples/...
```

## Design notes

* URL fields should be `*url.URL` in structs; direct `url.URL` is
  rejected to avoid copying pitfalls.
* `Must*` helpers panic on missing/invalid values; prefer non-`Must`
  forms for user input.
* The package uses stdlib-only dependencies while providing a clean API
  for environment variable management and struct binding.
* For nested structs, use a flattened structure or bind them separately.