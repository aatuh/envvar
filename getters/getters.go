package getters

import (
	"errors"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aatuh/envvar/v2/types"
)

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

// GetInt64 returns the value as an int64.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - int64: The value.
//   - error: The error if the value is not present.
func GetInt64(key string) (int64, error) {
	v, ok := Get(key)
	if !ok {
		return 0, missingErr(key)
	}
	i64, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	if err != nil {
		return 0, typeErr(key, "int64", v)
	}
	return i64, nil
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
	v, ok := Get(key)
	if !ok {
		return def
	}
	i64, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	if err != nil {
		return def
	}
	return i64
}

// MustGetInt64 returns the value as an int64 or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - int64: The value.
func MustGetInt64(key string) int64 {
	v, err := GetInt64(key)
	if err != nil {
		panic(err)
	}
	return v
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
	v, ok := Get(key)
	if !ok {
		return 0, missingErr(key)
	}
	u64, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64)
	if err != nil {
		return 0, typeErr(key, "uint", v)
	}
	return uint(u64), nil
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
	v, ok := Get(key)
	if !ok {
		return def
	}
	u64, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64)
	if err != nil {
		return def
	}
	return uint(u64)
}

// MustGetUint returns the value as a uint or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - uint: The value.
func MustGetUint(key string) uint {
	v, err := GetUint(key)
	if err != nil {
		panic(err)
	}
	return v
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
	v, ok := Get(key)
	if !ok {
		return 0, missingErr(key)
	}
	u64, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64)
	if err != nil {
		return 0, typeErr(key, "uint64", v)
	}
	return u64, nil
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
	v, ok := Get(key)
	if !ok {
		return def
	}
	u64, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64)
	if err != nil {
		return def
	}
	return u64
}

// MustGetUint64 returns the value as a uint64 or panics if not present.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - uint64: The value.
func MustGetUint64(key string) uint64 {
	v, err := GetUint64(key)
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

// GetRaw returns a value with expansion applied. Expansion supports
// "${NAME}" and "${NAME:-default}" using current process env.
//
// Parameters:
//   - key: The key to get.
//
// Returns:
//   - string: The value.
//   - bool: The boolean indicating presence.
func GetRaw(key string) (string, bool) {
	start := time.Now()
	v, ok := os.LookupEnv(key)
	var err error
	if ok {
		v = expand(v)
	}
	types.CallOnGet(key, ok, err, time.Since(start))
	return v, ok
}

// ParseBoolValue parses a boolean value.
//
// Parameters:
//   - v: The value to parse.
//
// Returns:
//   - bool: The boolean value.
//   - error: The error if the parsing fails.
func ParseBoolValue(v string) (bool, error) {
	s := strings.TrimSpace(strings.ToLower(v))
	switch s {
	case "1", "t", "true", "y", "yes", "on":
		return true, nil
	case "0", "f", "false", "n", "no", "off":
		return false, nil
	default:
		return false, errors.New("invalid boolean: " + v)
	}
}

// SplitAndTrim splits a string into a slice of strings and trims each string.
//
// Parameters:
//   - s: The string to split.
//   - sep: The separator to split on.
//
// Returns:
//   - []string: The slice of strings.
func SplitAndTrim(s, sep string) []string {
	raw := strings.Split(s, sep)
	out := make([]string, 0, len(raw))
	for _, p := range raw {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
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

// expand applies ${NAME} and ${NAME:-def} using process env first.
func expand(s string) string {
	// First handle ${NAME} and ${NAME:-def} ourselves to preserve defaults,
	// then allow $NAME and ${NAME} leftovers via os.ExpandEnv.
	s = expandWithLookup(s, os.LookupEnv)
	return os.ExpandEnv(s)
}

// expandWithLookup is a generic expander that resolves ${NAME:-def}
// with a provided lookup function.
func expandWithLookup(s string, look func(string) (string, bool)) string {
	// Handle ${NAME:-default} segments. Keep this non-nesting for
	// clarity and performance.
	for {
		i := strings.Index(s, "${")
		if i < 0 {
			break
		}
		j := strings.Index(s[i:], "}")
		if j < 0 {
			break
		}
		j += i
		inner := s[i+2 : j]
		name, def, hasDef := strings.Cut(inner, ":-")
		if name == "" {
			s = s[:i] + s[j+1:]
			continue
		}
		if v, ok := look(name); ok {
			s = s[:i] + v + s[j+1:]
		} else if hasDef {
			s = s[:i] + def + s[j+1:]
		} else {
			// Missing and no default -> drop to empty.
			s = s[:i] + "" + s[j+1:]
		}
	}
	return s
}
