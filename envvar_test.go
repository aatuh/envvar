package envvar

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	itypes "github.com/aatuh/envvar/v2/types"
)

// testHook implements types.Hook for integration tests.
type testHook struct {
	loads int
	gets  int
}

func (h *testHook) OnLoad(source string, keys int)                        { h.loads++ }
func (h *testHook) OnGet(key string, ok bool, err error, d time.Duration) { h.gets++ }

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

func TestDumpRedacted(t *testing.T) {
	t.Setenv("PUBLIC_VAL", "ok")
	t.Setenv("MY_SECRET", "shh")
	t.Setenv("API_KEY", "abc")
	t.Setenv("PASSWORD", "p")

	m := DumpRedacted()
	if m["PUBLIC_VAL"] != "ok" {
		t.Fatalf("PUBLIC_VAL not preserved: %q", m["PUBLIC_VAL"])
	}
	if m["MY_SECRET"] != "***" || m["API_KEY"] != "***" || m["PASSWORD"] != "***" {
		t.Fatalf("secret redaction failed: %#v", m)
	}
}

func TestBindBasicAndValidation(t *testing.T) {
	type C struct {
		Port    int           `env:"PORT,required" validate:"min=1,max=65535"`
		Debug   bool          `env:"DEBUG"`
		Timeout time.Duration `env:"TIMEOUT" envdef:"5s" validate:"min=1s,max=10s"`
		DSN     *url.URL      `env:"DATABASE_URL,required"`
		Tags    []string      `env:"TAGS" envsep:"|" validate:"oneof=a|b|c"`
		Raw     string        `env:"RAW"`
	}

	t.Setenv("PORT", "9090")
	t.Setenv("DEBUG", "true")
	t.Setenv("DATABASE_URL", "https://db.local")
	t.Setenv("TAGS", "a| c")
	t.Setenv("RAW", "${PORT}")

	var c C
	if err := Bind(&c); err != nil {
		t.Fatalf("Bind err: %v", err)
	}
	if c.Port != 9090 || !c.Debug || c.Timeout != 5*time.Second {
		t.Fatalf("Bind values wrong: %+v", c)
	}
	if c.DSN == nil || c.DSN.Host != "db.local" {
		t.Fatalf("URL pointer not set correctly: %+v", c.DSN)
	}
	if got := strings.Join(c.Tags, "|"); got != "a|c" {
		t.Fatalf("Tags wrong: %v", c.Tags)
	}
	if c.Raw != "9090" { // expansion applied via GetRaw
		t.Fatalf("Raw expansion failed: %q", c.Raw)
	}
}

func TestBindJSONAndPointerErrors(t *testing.T) {
	type C struct {
		M map[string]int `env:"MAP" envjson:"true"`
		U url.URL        `env:"URL_DIRECT"` // should cause error; require *url.URL
	}
	t.Setenv("MAP", `{"a":1,"b":2}`)
	t.Setenv("URL_DIRECT", "https://example.com")
	var c C
	err := Bind(&c)
	if err == nil {
		t.Fatalf("Bind expected error due to url.URL field")
	}
	if _, ok := err.(MultiError); !ok {
		t.Fatalf("want MultiError, got %T", err)
	}
	if c.M["a"] != 1 || c.M["b"] != 2 {
		t.Fatalf("JSON binding failed: %#v", c.M)
	}
	if !strings.Contains(err.Error(), "use *url.URL") {
		t.Fatalf("error should mention pointer URL: %v", err)
	}
}

func TestBindPrefixAndAggregation(t *testing.T) {
	type C struct {
		Port int    `env:"PORT,required" validate:"min=2"`
		Mode string `env:"MODE,required" validate:"oneof=a|b"`
	}
	t.Setenv("MYAPP_PORT", "1") // violates min=2
	// Missing MODE to produce missing error
	var c C
	err := BindWithPrefix(&c, "MYAPP_")
	if err == nil {
		t.Fatalf("expected MultiError with multiple issues")
	}
	me, ok := err.(MultiError)
	if !ok || len(me) < 2 {
		t.Fatalf("expected aggregated errors, got: %T %v", err, err)
	}
}

func TestReadFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env.test")
	content := "# comment\nKV1=aaa\nKV2=bbb\n"
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	m, err := ReadFile(p)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if m["KV1"] != "aaa" || m["KV2"] != "bbb" {
		t.Fatalf("map mismatch: %#v", m)
	}
}

func TestHooksIntegration(t *testing.T) {
	// Separate test to verify hooks without interfering with other tests
	h := &testHook{}
	itypes.SetHook(h)
	defer itypes.SetHook(nil)

	// Trigger OnGet via Get
	t.Setenv("HX_KEY", "val")
	if _, ok := Get("HX_KEY"); !ok {
		t.Fatal("Get should find HX_KEY")
	}

	// Trigger OnLoad via LoadOnce on a temp file
	dir := t.TempDir()
	p := filepath.Join(dir, ".env.hook")
	if err := os.WriteFile(p, []byte("HLOAD=1\n"), 0644); err != nil {
		t.Fatal(err)
	}
	// Use a unique key to avoid any prior LoadOnce effects being invisible
	_ = os.Unsetenv("HLOAD")
	_ = LoadOnce([]string{p})

	if h.gets < 1 {
		t.Fatalf("expected OnGet calls, got %d", h.gets)
	}
	// loads may be 0 if LoadOnce already ran; ensure code path safe.
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
