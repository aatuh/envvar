package loaders

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

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
