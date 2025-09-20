package envvar

import (
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aatuh/envvar/types"
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
	if err := LoadOnce(paths); err != nil {
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
	return GetRaw(key)
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
	if v, ok := Get(key); ok {
		return v
	}
	return def
}

// MustGet returns the value or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - string: The value.
func MustGet(key string) string {
	v, ok := Get(key)
	if !ok {
		panic("envvar: missing required " + key)
	}
	return v
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
	v, ok := Get(key)
	if !ok {
		return "", missingErr(key)
	}
	return v, nil
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
	return parseBool(key)
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
	v, ok := Get(key)
	if !ok {
		return def
	}
	b, err := ParseBoolValue(v)
	if err != nil {
		return def
	}
	return b
}

// MustGetBool returns the value as a boolean or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - bool: The value.
//   - error: The error if the value is not present.
func MustGetBool(key string) bool {
	b, err := GetBool(key)
	if err != nil {
		panic(err)
	}
	return b
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
	v, ok := Get(key)
	if !ok {
		return 0, missingErr(key)
	}
	i64, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	if err != nil {
		return 0, typeErr(key, "int", v)
	}
	return int(i64), nil
}

// GetIntOr returns the value as an integer or a default if not present.
//
// Parameters:
//   - key: The key to get.
//   - def: The default value.
//
// Returns:
//   - int: The value or the default.
//   - error: The error if the value is not present.
func GetIntOr(key string, def int) int {
	v, ok := Get(key)
	if !ok {
		return def
	}
	i64, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	if err != nil {
		return def
	}
	return int(i64)
}

// MustGetInt returns the value as an integer or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - int: The value.
//   - error: The error if the value is not present.
func MustGetInt(key string) int {
	v, err := GetInt(key)
	if err != nil {
		panic(err)
	}
	return v
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
	v, ok := Get(key)
	if !ok {
		return 0, missingErr(key)
	}
	f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
	if err != nil {
		return 0, typeErr(key, "float64", v)
	}
	return f, nil
}

// GetFloat64Or returns the value as a float64 or a default if not present.
//
// Parameters:
//   - key: The key to get.
//   - def: The default value.
//
// Returns:
//   - float64: The value or the default.
//   - error: The error if the value is not present.
func GetFloat64Or(key string, def float64) float64 {
	v, ok := Get(key)
	if !ok {
		return def
	}
	f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
	if err != nil {
		return def
	}
	return f
}

// MustGetFloat64 returns the value as a float64 or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - float64: The value.
//   - error: The error if the value is not present.
func MustGetFloat64(key string) float64 {
	v, err := GetFloat64(key)
	if err != nil {
		panic(err)
	}
	return v
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
	v, ok := Get(key)
	if !ok {
		return 0, missingErr(key)
	}
	d, err := time.ParseDuration(strings.TrimSpace(v))
	if err != nil {
		return 0, typeErr(key, "duration", v)
	}
	return d, nil
}

// GetDurationOr returns the value as a duration or a default if not present.
//
// Parameters:
//   - key: The key to get.
//   - def: The default value.
//
// Returns:
//   - time.Duration: The value or the default.
//   - error: The error if the value is not present.
func GetDurationOr(key string, def time.Duration) time.Duration {
	v, ok := Get(key)
	if !ok {
		return def
	}
	d, err := time.ParseDuration(strings.TrimSpace(v))
	if err != nil {
		return def
	}
	return d
}

// MustGetDuration returns the value as a duration or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - time.Duration: The value.
//   - error: The error if the value is not present.
func MustGetDuration(key string) time.Duration {
	v, err := GetDuration(key)
	if err != nil {
		panic(err)
	}
	return v
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
	v, ok := Get(key)
	if !ok {
		return nil, missingErr(key)
	}
	u, err := url.Parse(strings.TrimSpace(v))
	if err != nil || u.Scheme == "" {
		return nil, typeErr(key, "url", v)
	}
	return u, nil
}

// MustGetURL returns the value as a URL or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - *url.URL: The value.
func MustGetURL(key string) *url.URL {
	u, err := GetURL(key)
	if err != nil {
		panic(err)
	}
	return u
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
	v, ok := Get(key)
	if !ok {
		return nil, missingErr(key)
	}
	ip := net.ParseIP(strings.TrimSpace(v))
	if ip == nil {
		return nil, typeErr(key, "ip", v)
	}
	return ip, nil
}

// MustGetIP returns the value as an IP or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - net.IP: The value.
//   - error: The error if the value is not present.
func MustGetIP(key string) net.IP {
	ip, err := GetIP(key)
	if err != nil {
		panic(err)
	}
	return ip
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
	return GetStringSliceSep(key, ",")
}

// MustGetStringSlice returns the value as a slice of strings or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - []string: The value.
//   - error: The error if the value is not present.
func MustGetStringSlice(key string) []string {
	v, err := GetStringSlice(key)
	if err != nil {
		panic(err)
	}
	return v
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
	v, ok := Get(key)
	if !ok {
		return nil, missingErr(key)
	}
	s := strings.TrimSpace(v)
	if s == "" {
		return []string{}, nil
	}
	parts := SplitAndTrim(s, sep)
	return parts, nil
}

// Generic typed getter using a converter.
//
// Parameters:
//   - key: The key to get.
//   - conv: The converter function.
//
// Returns:
//   - T: The value.
//   - error: The error if the value is not present.
func GetTyped[T any](key string, conv func(string) (T, error)) (T, error) {
	var zero T
	v, ok := Get(key)
	if !ok {
		return zero, missingErr(key)
	}
	return conv(strings.TrimSpace(v))
}

// MustGetTyped returns the value as a typed value or panics if not present.
//
// Parameters:
//   - key: The key to get.
//   - conv: The converter function.
//
// Returns:
//   - T: The value.
//   - error: The error if the value is not present.
func MustGetTyped[T any](key string, conv func(string) (T, error)) T {
	v, err := GetTyped(key, conv)
	if err != nil {
		panic(err)
	}
	return v
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

// missingErr returns a missing error.
func missingErr(key string) error {
	return &KeyError{Key: key, Kind: ErrMissing}
}

// typeErr returns a type error.
func typeErr(key, want, got string) error {
	return &KeyError{
		Key:  key,
		Kind: ErrType,
		Msg:  "want " + want + ", got " + got,
	}
}

// parseBool parses a boolean value.
func parseBool(key string) (bool, error) {
	v, ok := Get(key)
	if !ok {
		return false, missingErr(key)
	}
	return ParseBoolValue(v)
}
