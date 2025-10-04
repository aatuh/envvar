package binders

import (
	"net/url"
	"strings"
	"testing"
	"time"
)

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
