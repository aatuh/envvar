package validate

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ValidateField validates a field based on the tag and separator.
//
// Parameters:
//   - v: The value to validate.
//   - tag: The tag to validate.
//   - sep: The separator to validate.
//
// Returns:
//   - error: The error if the validation fails.
func ValidateField(v reflect.Value, tag, sep string) error {
	rules := parseRules(tag)
	for k, val := range rules {
		switch k {
		case "min":
			if err := checkMin(v, val); err != nil {
				return err
			}
		case "max":
			if err := checkMax(v, val); err != nil {
				return err
			}
		case "oneof":
			if err := checkOneOf(v, val, sep); err != nil {
				return err
			}
		}
	}
	return nil
}

// parseRules parses the rules from the tag.
func parseRules(tag string) map[string]string {
	out := make(map[string]string)
	for _, part := range strings.Split(tag, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		k, v, ok := strings.Cut(part, "=")
		if ok {
			out[strings.TrimSpace(k)] = strings.TrimSpace(v)
		}
	}
	return out
}

// checkMin checks if the value is >= the minimum. Supports durations
// with human strings like "250ms", "5s", etc.
func checkMin(v reflect.Value, s string) error {
	// Handle time.Duration before generic ints (Duration is Int64).
	if v.Type().PkgPath() == "time" && v.Type().Name() == "Duration" {
		d, err := time.ParseDuration(s)
		if err != nil {
			return err
		}
		if time.Duration(v.Int()) < d {
			return fmt.Errorf("duration < %s", d)
		}
		return nil
	}

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		min, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		if v.Int() < min {
			return fmt.Errorf("must be >= %d", min)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		min, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		if v.Uint() < min {
			return fmt.Errorf("must be >= %d", min)
		}
	case reflect.String:
		min, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		if len(v.String()) < min {
			return fmt.Errorf("length must be >= %d", min)
		}
	}
	return nil
}

// checkMax checks if the value is <= the maximum. Supports durations
// with human strings like "250ms", "5s", etc.
func checkMax(v reflect.Value, s string) error {
	// Handle time.Duration before generic ints (Duration is Int64).
	if v.Type().PkgPath() == "time" && v.Type().Name() == "Duration" {
		d, err := time.ParseDuration(s)
		if err != nil {
			return err
		}
		if time.Duration(v.Int()) > d {
			return fmt.Errorf("duration > %s", d)
		}
		return nil
	}

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		max, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		if v.Int() > max {
			return fmt.Errorf("must be <= %d", max)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		max, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		if v.Uint() > max {
			return fmt.Errorf("must be <= %d", max)
		}
	case reflect.String:
		max, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		if len(v.String()) > max {
			return fmt.Errorf("length must be <= %d", max)
		}
	}
	return nil
}

// checkOneOf ensures the value is in the allowed set. Allowed values
// are always pipe-separated (e.g. "a|b|c"), regardless of envsep.
func checkOneOf(v reflect.Value, vals string, sep string) error {
	// Allowed values are pipe-separated regardless of envsep.
	allowed := strings.Split(vals, "|")

	switch v.Kind() {
	case reflect.String:
		s := v.String()
		for _, a := range allowed {
			if s == a {
				return nil
			}
		}
		return fmt.Errorf("must be one of %v", allowed)

	case reflect.Slice:
		// Only []string supported by binder; envsep controls how the
		// field value was split, not the allowed set.
		if v.Type().Elem().Kind() != reflect.String {
			return fmt.Errorf("oneof on string fields only")
		}
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).String()
			ok := false
			for _, a := range allowed {
				if item == a {
					ok = true
					break
				}
			}
			if !ok {
				return fmt.Errorf("value %q not in %v",
					item, allowed)
			}
		}
		return nil

	default:
		return fmt.Errorf("oneof supports string or []string")
	}
}
