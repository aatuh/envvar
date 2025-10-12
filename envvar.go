package envvar

import (
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aatuh/envvar/v2/binders"
	"github.com/aatuh/envvar/v2/getters"
	"github.com/aatuh/envvar/v2/lazy"
	"github.com/aatuh/envvar/v2/loaders"
	"github.com/aatuh/envvar/v2/types"
)

// Hook allows optional observability without adding dependencies.
// Provide your own implementation and register with SetHook.
type Hook = types.Hook

// SetHook installs a global hook. It is safe to call at program init.
//
// Parameters:
//   - h: The hook to install.
func SetHook(h Hook) {
	types.SetHook(h)
}

// MustLoadEnvVars loads variables from the first existing path in paths.
// If paths is nil, it tries ".env" then "/env/.env". It panics on
// read/parse error. Re-entrant calls are no-ops.
//
// Parameters:
//   - paths: The paths to load.
func MustLoadEnvVars(paths []string) {
	if err := loaders.LoadOnce(paths); err != nil {
		panic(err)
	}
}

// Get returns the raw value and a boolean indicating presence.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - string: The raw value.
//   - bool: The boolean indicating presence.
func Get(key string) (string, bool) {
	return getters.Get(key)
}

// GetOr returns the value or a default if not present.
//
// Parameters:
//   - key: The key to get.
//   - def: The default value.
//
// Returns:
//   - string: The value or the default.
func GetOr(key, def string) string {
	return getters.GetOr(key, def)
}

// MustGet returns the value or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - string: The value.
func MustGet(key string) string {
	return getters.MustGet(key)
}

// GetOrErr returns the value or an error if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - string: The value.
//   - error: The error if the value is not present.
func GetOrErr(key string) (string, error) {
	return getters.GetOrErr(key)
}

// GetBool returns the value as a boolean.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - bool: The value.
//   - error: The error if the value is not present.
func GetBool(key string) (bool, error) {
	return getters.GetBool(key)
}

// GetBoolOr returns the value as a boolean or a default if not present.
//
// Parameters:
//   - key: The key to get.
//   - def: The default value.
//
// Returns:
//   - bool: The value or the default.
func GetBoolOr(key string, def bool) bool {
	return getters.GetBoolOr(key, def)
}

// MustGetBool returns the value as a boolean or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - bool: The value.
func MustGetBool(key string) bool {
	return getters.MustGetBool(key)
}

// GetInt returns the value as an integer.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - int: The value.
//   - error: The error if the value is not present.
func GetInt(key string) (int, error) {
	return getters.GetInt(key)
}

// GetIntOr returns the value as an integer or a default if not present.
//
// Parameters:
//   - key: The key to get.
//   - def: The default value.
//
// Returns:
//   - int: The value or the default.
func GetIntOr(key string, def int) int {
	return getters.GetIntOr(key, def)
}

// MustGetInt returns the value as an integer or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - int: The value.
func MustGetInt(key string) int {
	return getters.MustGetInt(key)
}

// GetInt64 returns the value as an int64.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - int64: The value.
//   - error: The error if the value is not present.
func GetInt64(key string) (int64, error) {
	return getters.GetInt64(key)
}

// GetInt64Or returns the value as an int64 or a default if not present.
//
// Parameters:
//   - key: The key to get.
//   - def: The default value.
//
// Returns:
//   - int64: The value or the default.
func GetInt64Or(key string, def int64) int64 {
	return getters.GetInt64Or(key, def)
}

// MustGetInt64 returns the value as an int64 or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - int64: The value.
func MustGetInt64(key string) int64 {
	return getters.MustGetInt64(key)
}

// GetUint returns the value as a uint.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - uint: The value.
//   - error: The error if the value is not present.
func GetUint(key string) (uint, error) {
	return getters.GetUint(key)
}

// GetUintOr returns the value as a uint or a default if not present.
//
// Parameters:
//   - key: The key to get.
//   - def: The default value.
//
// Returns:
//   - uint: The value or the default.
func GetUintOr(key string, def uint) uint {
	return getters.GetUintOr(key, def)
}

// MustGetUint returns the value as a uint or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - uint: The value.
func MustGetUint(key string) uint {
	return getters.MustGetUint(key)
}

// GetUint64 returns the value as a uint64.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - uint64: The value.
//   - error: The error if the value is not present.
func GetUint64(key string) (uint64, error) {
	return getters.GetUint64(key)
}

// GetUint64Or returns the value as a uint64 or a default if not present.
//
// Parameters:
//   - key: The key to get.
//   - def: The default value.
//
// Returns:
//   - uint64: The value or the default.
func GetUint64Or(key string, def uint64) uint64 {
	return getters.GetUint64Or(key, def)
}

// MustGetUint64 returns the value as a uint64 or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - uint64: The value.
func MustGetUint64(key string) uint64 {
	return getters.MustGetUint64(key)
}

// GetFloat64 returns the value as a float64.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - float64: The value.
//   - error: The error if the value is not present.
func GetFloat64(key string) (float64, error) {
	return getters.GetFloat64(key)
}

// GetFloat64Or returns the value as a float64 or a default if not present.
//
// Parameters:
//   - key: The key to get.
//   - def: The default value.
//
// Returns:
//   - float64: The value or the default.
func GetFloat64Or(key string, def float64) float64 {
	return getters.GetFloat64Or(key, def)
}

// MustGetFloat64 returns the value as a float64 or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - float64: The value.
func MustGetFloat64(key string) float64 {
	return getters.MustGetFloat64(key)
}

// GetDuration returns the value as a duration.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - time.Duration: The value.
//   - error: The error if the value is not present.
func GetDuration(key string) (time.Duration, error) {
	return getters.GetDuration(key)
}

// GetDurationOr returns the value as a duration or a default if not present.
//
// Parameters:
//   - key: The key to get.
//   - def: The default value.
//
// Returns:
//   - time.Duration: The value or the default.
func GetDurationOr(key string, def time.Duration) time.Duration {
	return getters.GetDurationOr(key, def)
}

// MustGetDuration returns the value as a duration or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - time.Duration: The value.
func MustGetDuration(key string) time.Duration {
	return getters.MustGetDuration(key)
}

// GetURL returns the value as a URL.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - *url.URL: The value.
//   - error: The error if the value is not present.
func GetURL(key string) (*url.URL, error) {
	return getters.GetURL(key)
}

// MustGetURL returns the value as a URL or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - *url.URL: The value.
func MustGetURL(key string) *url.URL {
	return getters.MustGetURL(key)
}

// GetIP returns the value as an IP.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - net.IP: The value.
//   - error: The error if the value is not present.
func GetIP(key string) (net.IP, error) {
	return getters.GetIP(key)
}

// MustGetIP returns the value as an IP or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - net.IP: The value.
func MustGetIP(key string) net.IP {
	return getters.MustGetIP(key)
}

// GetStringSlice returns the value as a slice of strings.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - []string: The value.
//   - error: The error if the value is not present.
func GetStringSlice(key string) ([]string, error) {
	return getters.GetStringSlice(key)
}

// MustGetStringSlice returns the value as a slice of strings or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - []string: The value.
func MustGetStringSlice(key string) []string {
	return getters.MustGetStringSlice(key)
}

// GetStringSliceSep returns the value as a slice of strings with a custom separator.
//
// Parameters:
//   - key: The key to get.
//   - sep: The separator.
//
// Returns:
//   - []string: The value.
//   - error: The error if the value is not present.
func GetStringSliceSep(key, sep string) ([]string, error) {
	return getters.GetStringSliceSep(key, sep)
}

// GetTyped returns the value as a typed value using a converter.
//
// Parameters:
//   - key: The key to get.
//   - conv: The converter function.
//
// Returns:
//   - T: The value.
//   - error: The error if the value is not present.
func GetTyped[T any](key string, conv func(string) (T, error)) (T, error) {
	return getters.GetTyped(key, conv)
}

// MustGetTyped returns the value as a typed value or panics if not present.
//
// Parameters:
//   - key: The key to get.
//   - conv: The converter function.
//
// Returns:
//   - T: The value.
func MustGetTyped[T any](key string, conv func(string) (T, error)) T {
	return getters.MustGetTyped(key, conv)
}

// Bind populates a struct from the process environment using `env` and
// `validate` tags. See BindWithPrefix for details.
//
// Parameters:
//   - dst: The destination.
//
// Returns:
//   - error: The error if the binding fails.
func Bind(dst any) error {
	return binders.Bind(dst)
}

// BindWithPrefix is like Bind but first tries variables with the given
// prefix. For example with prefix "MYAPP_", field `env:"PORT"` resolves
// "MYAPP_PORT" if present, else falls back to "PORT".
//
// Parameters:
//   - dst: The destination.
//   - prefix: The prefix.
//
// Returns:
//   - error: The error if the binding fails.
func BindWithPrefix(dst any, prefix string) error {
	return binders.BindWithPrefix(dst, prefix)
}

// MustBind panics on binding errors.
//
// Parameters:
//   - dst: The destination.
func MustBind(dst any) {
	binders.MustBind(dst)
}

// MustBindWithPrefix panics on binding errors.
//
// Parameters:
//   - dst: The destination.
//   - prefix: The prefix.
func MustBindWithPrefix(dst any, prefix string) {
	binders.MustBindWithPrefix(dst, prefix)
}

// LazyString returns a function that returns the value of the environment
// variable with the given key.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - func(): The function that returns the value of the environment variable
//     with the given key.
func LazyString(key string) func() string {
	return lazy.LazyString(key)
}

// LazyBool returns a function that returns the value of the environment
// variable with the given key as a boolean.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - func(): The function that returns the value of the environment variable
//     with the given key as a boolean.
func LazyBool(key string) func() bool {
	return lazy.LazyBool(key)
}

// LazyInt returns a function that returns the value of the environment
// variable with the given key as an integer.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - func(): The function that returns the value of the environment variable
//     with the given key as an integer.
func LazyInt(key string) func() int {
	return lazy.LazyInt(key)
}

// LazyInt64 returns a function that returns the value of the environment
// variable with the given key as an int64.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - func(): The function that returns the value of the environment variable
//     with the given key as an int64.
func LazyInt64(key string) func() int64 {
	return lazy.LazyInt64(key)
}

// LazyUint returns a function that returns the value of the environment
// variable with the given key as a uint.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - func(): The function that returns the value of the environment variable
//     with the given key as a uint.
func LazyUint(key string) func() uint {
	return lazy.LazyUint(key)
}

// LazyUint64 returns a function that returns the value of the environment
// variable with the given key as a uint64.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - func(): The function that returns the value of the environment variable
//     with the given key as a uint64.
func LazyUint64(key string) func() uint64 {
	return lazy.LazyUint64(key)
}

// LazyFloat64 returns a function that returns the value of the environment
// variable with the given key as a float64.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - func(): The function that returns the value of the environment variable
//     with the given key as a float64.
func LazyFloat64(key string) func() float64 {
	return lazy.LazyFloat64(key)
}

// LazyDuration returns a function that returns the value of the environment
// variable with the given key as a duration.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - func(): The function that returns the value of the environment variable
//     with the given key as a duration.
func LazyDuration(key string) func() time.Duration {
	return lazy.LazyDuration(key)
}

// LazyStringSlice returns a function that returns the value of the environment
// variable with the given key as a slice of strings.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - func(): The function that returns the value of the environment variable
//     with the given key as a slice of strings.
func LazyStringSlice(key string) func() []string {
	return lazy.LazyStringSlice(key)
}

// LazyTyped returns a function that returns the value of the environment
// variable with the given key as a typed value.
//
// Parameters:
//   - key: The key to get.
//   - conv: The converter function.
//
// Returns:
//   - func(): The function that returns the value of the environment variable
//     with the given key as a typed value.
func LazyTyped[T any](key string, conv func(string) (T, error)) func() T {
	return lazy.LazyTyped(key, conv)
}

// DumpRedacted returns environment as a map with secret-like values
// redacted. Redaction is heuristic: keys containing "SECRET", "TOKEN",
// "KEY", or "PASSWORD" are masked.
//
// Returns:
//   - map[string]string: The environment as a map with secret-like values redacted.
func DumpRedacted() map[string]string {
	env := os.Environ()
	out := make(map[string]string, len(env))
	for _, kv := range env {
		k, v, ok := strings.Cut(kv, "=")
		if !ok {
			continue
		}
		upper := strings.ToUpper(k)
		if strings.Contains(upper, "SECRET") ||
			strings.Contains(upper, "TOKEN") ||
			strings.Contains(upper, "PASSWORD") ||
			strings.HasSuffix(upper, "_KEY") {
			out[k] = "***"
		} else {
			out[k] = v
		}
	}
	return out
}
