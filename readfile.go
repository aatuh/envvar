package envvar

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aatuh/envvar/v2/internal"
	"github.com/aatuh/envvar/v2/types"
)

// SetEnvVars sets the provided map into process env. Values overwrite
// existing ones. Keys with empty value are set to "".
//
// Parameters:
//   - m: The map to set.
//
// Returns:
//   - error: The error if the setting fails.
func SetEnvVars(m map[string]string) error {
	for k, v := range m {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	return nil
}

var (
	loadOnceGuard sync.Once
	loadErr       error

	// defaultPaths are tried when caller passes nil.
	defaultPaths = []string{".env", "/env/.env"}
)

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
		v = internal.Expand(v)
	}
	types.CallOnGet(key, ok, err, time.Since(start))
	return v, ok
}

// LoadOnce loads the environment variables from the given paths.
//
// Parameters:
//   - paths: The paths to load.
//
// Returns:
//   - error: The error if the loading fails.
func LoadOnce(paths []string) error {
	loadOnceGuard.Do(func() {
		if len(paths) == 0 {
			paths = defaultPaths
		}
		for _, p := range paths {
			info, err := os.Stat(p)
			if err != nil || info.IsDir() {
				continue
			}
			m, err := ReadFile(p)
			if err != nil {
				loadErr = err
				return
			}
			_ = SetEnvVars(m)
			types.CallOnLoad(p, len(m))
			return
		}
		// Not an error if none exist.
	})
	return loadErr
}

// ReadFile reads the environment variables from the given path.
//
// Parameters:
//   - path: The path to read.
//
// Returns:
//   - map[string]string: The map of key-value pairs.
//   - error: The error if the reading fails.
func ReadFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := make(map[string]string)
	sc := bufio.NewScanner(f)
	ln := 0
	for sc.Scan() {
		ln++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			return nil, errors.New("envvar: invalid line " +
				filepath.Base(path) + ":" + strconvI(ln))
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		m[k] = v
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return m, nil
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

// strconvI converts an integer to a string.
func strconvI(i int) string {
	// Small helper avoiding strconv import here.
	const digits = "0123456789"
	if i == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	n := i
	if n < 0 {
		n = -n
	}
	for n > 0 {
		pos--
		//nolint:gosec // safe index into digits
		buf[pos] = digits[n%10]
		n /= 10
	}
	if i < 0 {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
