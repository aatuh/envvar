package internal

import (
	"os"
	"strings"
)

// ExpandMap expands ${NAME} and ${NAME:-def} in the provided map,
// referencing keys in the same map first, then falling back to the
// process environment. Returns a new map.
//
// Parameters:
//   - in: The map to expand.
//
// Returns:
//   - map[string]string: The expanded map.
func ExpandMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return map[string]string{}
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	// Resolve with bounded iterations to avoid cycles.
	for iter := 0; iter < 10; iter++ {
		stable := true
		for k, v := range out {
			nv := expandWithLookup(v, func(name string) (string, bool) {
				if vv, ok := out[name]; ok {
					return vv, true
				}
				if vv, ok := os.LookupEnv(name); ok {
					return vv, true
				}
				return os.LookupEnv(name)
			})
			if nv != v {
				out[k] = nv
				stable = false
			}
		}
		if stable {
			break
		}
	}
	return out
}

// MustExpandMap is like ExpandMap but intended for init-time usage.
// where expansion failure should abort. Currently expansion cannot
// fail, so this mirrors ExpandMap and exists for API completeness.
//
// Parameters:
//   - in: The map to expand.
//
// Returns:
//   - map[string]string: The expanded map.
func MustExpandMap(in map[string]string) map[string]string {
	return ExpandMap(in)
}

// Expand applies ${NAME} and ${NAME:-def} using process env first.
//
// Parameters:
//   - s: The string to expand.
//
// Returns:
//   - string: The expanded string.
func Expand(s string) string {
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
