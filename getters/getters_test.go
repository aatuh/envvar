package getters

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestGetAndExpansion(t *testing.T) {
	t.Setenv("EV_A", "1")
	t.Setenv("EV_B", "${EV_A}")
	t.Setenv("EV_C", "${MISSING_X:-def}")

	if v := MustGet("EV_B"); v != "1" {
		t.Fatalf("EV_B: want 1, got %q", v)
	}
	if v := MustGet("EV_C"); v != "def" {
		t.Fatalf("EV_C: want def, got %q", v)
	}
}

func TestParseBoolAndGetBoolOr(t *testing.T) {
	t.Setenv("FLAG_YES", "yes")
	t.Setenv("FLAG_BAD", "zzz")

	b, err := GetBool("FLAG_YES")
	if err != nil || !b {
		t.Fatalf("GetBool yes failed: %v %v", b, err)
	}
	if got := GetBoolOr("FLAG_BAD", true); got != true {
		t.Fatalf("GetBoolOr parse error should fallback to default")
	}
	if got := GetBoolOr("MISSING_FLAG", false); got != false {
		t.Fatalf("GetBoolOr missing should fallback to default")
	}
}

func TestNumericDurationURLStringSlice(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("RATE", "3.14")
	t.Setenv("TTL", "250ms")
	t.Setenv("DSN", "https://example.com/x")
	t.Setenv("TAGS", " a, b , ,c ")

	if v, err := GetInt("PORT"); err != nil || v != 8080 {
		t.Fatalf("GetInt: %v %v", v, err)
	}
	if v := GetIntOr("BADINT", 42); v != 42 {
		t.Fatalf("GetIntOr fallback failed: %v", v)
	}

	if v, err := GetFloat64("RATE"); err != nil || v < 3.139 || v > 3.141 {
		t.Fatalf("GetFloat64: %v %v", v, err)
	}

	if d, err := GetDuration("TTL"); err != nil || d != 250*time.Millisecond {
		t.Fatalf("GetDuration: %v %v", d, err)
	}

	if u, err := GetURL("DSN"); err != nil || u.Host != "example.com" {
		t.Fatalf("GetURL: %v %v", u, err)
	}

	sl, err := GetStringSlice("TAGS")
	if err != nil {
		t.Fatalf("GetStringSlice err: %v", err)
	}
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(sl, want) {
		t.Fatalf("GetStringSlice: want %v, got %v", want, sl)
	}
}

func TestGetTypedAndMustPanics(t *testing.T) {
	t.Setenv("UPPER", " go ")
	v, err := GetTyped("UPPER", func(s string) (string, error) { return strings.ToUpper(s), nil })
	if err != nil || v != "GO" {
		t.Fatalf("GetTyped: %q %v", v, err)
	}

	// MustGetInt panics on missing
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("MustGetInt should panic on missing key")
		}
	}()
	_ = MustGetInt("NO_SUCH_INT")
}

func TestErrorsKinds(t *testing.T) {
	t.Setenv("BAD_FLOAT", "nan?")
	if _, err := GetFloat64("BAD_FLOAT"); err == nil {
		t.Fatalf("expected type error for BAD_FLOAT")
	} else {
		var ke *KeyError
		if !errors.As(err, &ke) || ke.Kind != ErrType {
			t.Fatalf("want KeyError ErrType, got %T %v", err, err)
		}
	}
	if _, err := GetOrErr("REALLY_MISSING"); err == nil {
		t.Fatalf("expected missing error")
	} else {
		var ke *KeyError
		if !errors.As(err, &ke) || ke.Kind != ErrMissing {
			t.Fatalf("want KeyError ErrMissing, got %T %v", err, err)
		}
	}
}

func TestSplitAndTrimAndParseBoolValue(t *testing.T) {
	parts := SplitAndTrim(" a , , b, c ", ",")
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(parts, want) {
		t.Fatalf("SplitAndTrim: want %v, got %v", want, parts)
	}
	cases := map[string]bool{
		"1": true, "t": true, "true": true, "y": true, "yes": true, "on": true,
		"0": false, "f": false, "false": false, "n": false, "no": false, "off": false,
	}
	for in, want := range cases {
		got, err := ParseBoolValue(in)
		if err != nil || got != want {
			t.Fatalf("ParseBoolValue(%q)=%v,%v", in, got, err)
		}
	}
	if _, err := ParseBoolValue("maybe"); err == nil {
		t.Fatalf("expected error for invalid boolean")
	}
}
