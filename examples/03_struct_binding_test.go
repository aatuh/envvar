package examples

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/aatuh/envvar/v2"
)

// Basic struct binding - populate struct from environment
func TestBasicStructBinding(t *testing.T) {
	type Config struct {
		Port    int           `env:"PORT,required"`
		Debug   bool          `env:"DEBUG"`
		Timeout time.Duration `env:"TIMEOUT" envdef:"5s"`
		DSN     *url.URL      `env:"DATABASE_URL,required"`
		Tags    []string      `env:"TAGS" envsep:","`
	}

	// Set up test environment
	t.Setenv("PORT", "8080")
	t.Setenv("DEBUG", "true")
	t.Setenv("DATABASE_URL", "https://db.example.com")
	t.Setenv("TAGS", "web,api,go")

	var cfg Config
	if err := envvar.Bind(&cfg); err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	if cfg.Port != 8080 {
		t.Fatalf("Port binding failed: %v", cfg.Port)
	}
	if !cfg.Debug {
		t.Fatalf("Debug binding failed: %v", cfg.Debug)
	}
	if cfg.Timeout != 5*time.Second {
		t.Fatalf("Timeout default failed: %v", cfg.Timeout)
	}
	if cfg.DSN == nil || cfg.DSN.Host != "db.example.com" {
		t.Fatalf("URL binding failed: %v", cfg.DSN)
	}
	if len(cfg.Tags) != 3 {
		t.Fatalf("Tags binding failed: %v", cfg.Tags)
	}
}

// Struct binding with defaults
func TestStructBindingWithDefaults(t *testing.T) {
	type Config struct {
		Port    int           `env:"PORT" envdef:"8080"`
		Timeout time.Duration `env:"TIMEOUT" envdef:"5s"`
		Mode    string        `env:"MODE" envdef:"development"`
		Tags    []string      `env:"TAGS" envdef:"web,api"`
	}

	// Test with defaults
	var cfg Config
	if err := envvar.Bind(&cfg); err != nil {
		t.Fatalf("Config with defaults should bind successfully: %v", err)
	}

	if cfg.Port != 8080 {
		t.Fatalf("Port default should be 8080, got %v", cfg.Port)
	}
	if cfg.Timeout != 5*time.Second {
		t.Fatalf("Timeout default should be 5s, got %v", cfg.Timeout)
	}
	if cfg.Mode != "development" {
		t.Fatalf("Mode default should be development, got %v", cfg.Mode)
	}
	if len(cfg.Tags) != 2 {
		t.Fatalf("Tags default should have 2 items, got %v", cfg.Tags)
	}
}

// JSON field binding
func TestJSONFieldBinding(t *testing.T) {
	type Config struct {
		Database map[string]interface{} `env:"DB_CONFIG" envjson:"true"`
		Features []string               `env:"FEATURES" envjson:"true"`
	}

	// Set up JSON environment variables
	t.Setenv("DB_CONFIG", `{"host":"localhost","port":5432,"ssl":true}`)
	t.Setenv("FEATURES", `["auth","logging","metrics"]`)

	var cfg Config
	if err := envvar.Bind(&cfg); err != nil {
		t.Fatalf("JSON binding failed: %v", err)
	}

	if cfg.Database["host"] != "localhost" {
		t.Fatalf("Database config binding failed: %v", cfg.Database)
	}
	if len(cfg.Features) != 3 {
		t.Fatalf("Features binding failed: %v", cfg.Features)
	}
}

// Prefix binding - try prefixed variables first
func TestPrefixBinding(t *testing.T) {
	type Config struct {
		Port int    `env:"PORT,required"`
		Mode string `env:"MODE,required"`
	}

	// Set up prefixed environment variables
	t.Setenv("MYAPP_PORT", "9090")
	t.Setenv("MYAPP_MODE", "production")
	// Also set unprefixed (should be ignored due to prefix)
	t.Setenv("PORT", "8080")
	t.Setenv("MODE", "development")

	var cfg Config
	if err := envvar.BindWithPrefix(&cfg, "MYAPP_"); err != nil {
		t.Fatalf("Prefix binding failed: %v", err)
	}

	if cfg.Port != 9090 {
		t.Fatalf("Prefixed port binding failed: %v", cfg.Port)
	}
	if cfg.Mode != "production" {
		t.Fatalf("Prefixed mode binding failed: %v", cfg.Mode)
	}
}

// MustBind - panic on binding errors
func TestMustBind(t *testing.T) {
	type Config struct {
		Port int `env:"PORT,required"`
	}

	// Valid configuration
	t.Setenv("PORT", "8080")
	var cfg Config
	envvar.MustBind(&cfg) // Should not panic
	if cfg.Port != 8080 {
		t.Fatalf("MustBind failed: %v", cfg.Port)
	}

	// Invalid configuration - missing required field
	os.Unsetenv("PORT")
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("MustBind should panic on missing required field")
		}
	}()
	envvar.MustBind(&cfg)
}
