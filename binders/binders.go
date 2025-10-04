package binders

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/aatuh/envvar/v2/expand"
	"github.com/aatuh/envvar/v2/validate"
)

// Bind populates a struct from the process environment using `env` and
// `validate` tags. See BindWithPrefix for details.
//
// Parameters:
//   - dst: The destination.
//
// Returns:
//   - error: The error if the binding fails.
func Bind(dst any) error {
	return bindWithOptions(dst, "")
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
	return bindWithOptions(dst, prefix)
}

// MustBind panics on binding errors.
//
// Parameters:
//   - dst: The destination.
func MustBind(dst any) {
	if err := Bind(dst); err != nil {
		panic(err)
	}
}

// MustBindWithPrefix panics on binding errors.
//
// Parameters:
//   - dst: The destination.
//   - prefix: The prefix.
func MustBindWithPrefix(dst any, prefix string) {
	if err := BindWithPrefix(dst, prefix); err != nil {
		panic(err)
	}
}

// bindWithOptions binds the options.
func bindWithOptions(dst any, prefix string) error {
	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("envvar: Bind expects pointer to struct")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("envvar: Bind expects pointer to struct")
	}

	var errs MultiError

	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if f.PkgPath != "" { // unexported
			continue
		}
		ev, ok := f.Tag.Lookup("env")
		if !ok {
			continue
		}
		name, req := parseEnvTag(ev)
		def := f.Tag.Get("envdef")
		sep := f.Tag.Get("envsep")
		if sep == "" {
			sep = ","
		}
		jsonMode := strings.EqualFold(f.Tag.Get("envjson"), "true")

		raw, exists := lookupPrefixed(prefix, name)
		if !exists && def != "" {
			raw = def
			exists = true
		}
		if !exists && req {
			errs = append(errs, missingErr(name))
			continue
		}
		if !exists {
			continue
		}
		raw = expand.Expand(raw)

		fv := rv.Field(i)
		if !fv.CanSet() {
			continue
		}
		if err := setField(fv, raw, sep, jsonMode); err != nil {
			errs = append(errs, fmt.Errorf("envvar: %s: %w", name, err))
			continue
		}
		// Validation
		if vtag, ok := f.Tag.Lookup("validate"); ok && vtag != "" {
			if err := validate.ValidateField(fv, vtag, sep); err != nil {
				errs = append(errs, fmt.Errorf("envvar: %s: %w",
					name, err))
			}
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// lookupPrefixed looks up the prefixed name.
func lookupPrefixed(prefix, name string) (string, bool) {
	if prefix != "" {
		if v, ok := os.LookupEnv(prefix + name); ok {
			return v, true
		}
	}
	return os.LookupEnv(name)
}

// parseEnvTag parses the env tag.
func parseEnvTag(tag string) (name string, required bool) {
	name = tag
	if i := strings.Index(tag, ","); i >= 0 {
		name = tag[:i]
		for _, part := range strings.Split(tag[i+1:], ",") {
			if strings.TrimSpace(part) == "required" {
				required = true
			}
		}
	}
	name = strings.TrimSpace(name)
	return
}

// setField sets the field.
func setField(v reflect.Value, raw, sep string, jsonMode bool) error {
	// If JSON mode is enabled, unmarshal into the field type.
	if jsonMode {
		return setFieldJSON(v, raw)
	}

	t := v.Type()
	kind := t.Kind()

	// Pointers
	if kind == reflect.Ptr {
		// Special-case *url.URL
		if t.Elem().PkgPath() == "net/url" && t.Elem().Name() == "URL" {
			u, err := url.Parse(raw)
			if err != nil || u.Scheme == "" {
				return fmt.Errorf("invalid url: %s", raw)
			}
			// Set the pointed value directly
			elem := reflect.New(t.Elem())
			elem.Elem().Set(reflect.ValueOf(*u))
			v.Set(elem)
			return nil
		}
		elem := reflect.New(t.Elem())
		if err := setField(elem.Elem(), raw, sep, false); err != nil {
			return err
		}
		v.Set(elem)
		return nil
	}

	switch kind {
	case reflect.String:
		v.SetString(raw)
		return nil
	case reflect.Bool:
		b, err := ParseBoolValue(raw)
		if err != nil {
			return err
		}
		v.SetBool(b)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		// Special case time.Duration
		if t.PkgPath() == "time" && t.Name() == "Duration" {
			d, err := time.ParseDuration(raw)
			if err != nil {
				return fmt.Errorf("invalid duration: %s", raw)
			}
			v.SetInt(int64(d))
			return nil
		}
		i, err := strconv.ParseInt(raw, 10, t.Bits())
		if err != nil {
			return fmt.Errorf("invalid int: %s", raw)
		}
		v.SetInt(i)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(raw, 10, t.Bits())
		if err != nil {
			return fmt.Errorf("invalid uint: %s", raw)
		}
		v.SetUint(u)
		return nil
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(raw, t.Bits())
		if err != nil {
			return fmt.Errorf("invalid float: %s", raw)
		}
		v.SetFloat(f)
		return nil
	case reflect.Slice:
		if t.Elem().Kind() != reflect.String {
			return fmt.Errorf("only []string slices supported")
		}
		parts := SplitAndTrim(raw, sep)
		sv := reflect.MakeSlice(t, len(parts), len(parts))
		for i := range parts {
			sv.Index(i).SetString(parts[i])
		}
		v.Set(sv)
		return nil
	case reflect.Struct:
		// url.URL supported via pointer. Direct struct is awkward;
		// keep a helpful error for clarity.
		if t.PkgPath() == "net/url" && t.Name() == "URL" {
			return fmt.Errorf("use *url.URL in struct, not url.URL")
		}
		return fmt.Errorf("unsupported struct type %s", t.String())
	default:
		return fmt.Errorf("unsupported kind %s", kind)
	}
}

// setFieldJSON sets the field as JSON.
func setFieldJSON(v reflect.Value, raw string) error {
	t := v.Type()
	kind := t.Kind()

	// If pointer, allocate and unmarshal into element.
	if kind == reflect.Ptr {
		elem := reflect.New(t.Elem())
		if err := json.Unmarshal([]byte(raw), elem.Interface()); err != nil {
			return fmt.Errorf("invalid json: %v", err)
		}
		v.Set(elem)
		return nil
	}

	// Non-pointer: unmarshal into a new value of the same type, then set.
	tmp := reflect.New(t).Interface()
	if err := json.Unmarshal([]byte(raw), tmp); err != nil {
		return fmt.Errorf("invalid json: %v", err)
	}
	v.Set(reflect.ValueOf(tmp).Elem())
	return nil
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
		return false, fmt.Errorf("invalid boolean: %s", v)
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
