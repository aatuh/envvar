package examples

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aatuh/envvar/v2"
)

// Variable expansion in environment variables
func TestVariableExpansionAdvanced(t *testing.T) {
	// Set up environment variables with expansion patterns
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DATABASE_URL", "postgres://${DB_HOST}:${DB_PORT}/mydb")
	t.Setenv("FALLBACK_URL", "postgres://${MISSING_HOST:-localhost}:${MISSING_PORT:-5432}/db")

	// Test expansion with existing variables
	if url, ok := envvar.Get("DATABASE_URL"); !ok || url != "postgres://localhost:5432/mydb" {
		t.Fatalf("Variable expansion failed: %v", url)
	}

	// Test expansion with fallback defaults
	if url, ok := envvar.Get("FALLBACK_URL"); !ok || url != "postgres://localhost:5432/db" {
		t.Fatalf("Variable expansion with defaults failed: %v", url)
	}
}

// Redacted environment dump - hide sensitive values
func TestRedactedDump(t *testing.T) {
	// Set up environment with sensitive and non-sensitive values
	t.Setenv("APP_NAME", "myapp")
	t.Setenv("PORT", "8080")
	t.Setenv("API_SECRET", "secret123")
	t.Setenv("DATABASE_PASSWORD", "dbpass")
	t.Setenv("JWT_TOKEN", "jwt123")
	t.Setenv("PRIVATE_KEY", "key123")

	// Get redacted environment dump
	redacted := envvar.DumpRedacted()

	// Non-sensitive values should be preserved
	if redacted["APP_NAME"] != "myapp" {
		t.Fatalf("APP_NAME should not be redacted: %v", redacted["APP_NAME"])
	}
	if redacted["PORT"] != "8080" {
		t.Fatalf("PORT should not be redacted: %v", redacted["PORT"])
	}

	// Sensitive values should be redacted
	if redacted["API_SECRET"] != "***" {
		t.Fatalf("API_SECRET should be redacted: %v", redacted["API_SECRET"])
	}
	if redacted["DATABASE_PASSWORD"] != "***" {
		t.Fatalf("DATABASE_PASSWORD should be redacted: %v", redacted["DATABASE_PASSWORD"])
	}
	if redacted["JWT_TOKEN"] != "***" {
		t.Fatalf("JWT_TOKEN should be redacted: %v", redacted["JWT_TOKEN"])
	}
	if redacted["PRIVATE_KEY"] != "***" {
		t.Fatalf("PRIVATE_KEY should be redacted: %v", redacted["PRIVATE_KEY"])
	}
}

// testHook implements Hook for tracking access
type testHook struct {
	loads int
	gets  int
}

func (h *testHook) OnLoad(source string, keys int)                        { h.loads++ }
func (h *testHook) OnGet(key string, ok bool, err error, d time.Duration) { h.gets++ }

// Observability hooks - track environment variable access
func TestObservabilityHooks(t *testing.T) {
	// Create a hook to track access
	hook := &testHook{}
	envvar.SetHook(hook)
	defer envvar.SetHook(nil)

	// Set up test environment
	t.Setenv("TRACKED_VAR", "value")
	t.Setenv("ANOTHER_VAR", "another")

	// Access environment variables (should trigger OnGet)
	envvar.MustGet("TRACKED_VAR")
	envvar.MustGet("ANOTHER_VAR")

	// Verify hook was called
	if hook.gets < 2 {
		t.Fatalf("Hook should track gets: %v", hook.gets)
	}
}

// Complex configuration with all features
func TestComplexConfiguration(t *testing.T) {
	type Config struct {
		// Basic types
		Port     int           `env:"PORT" envdef:"8080"`
		Timeout  time.Duration `env:"TIMEOUT" envdef:"30s"`
		Name     string        `env:"NAME" envdef:"myapp"`
		Mode     string        `env:"MODE" envdef:"development"`
		Features []string      `env:"FEATURES" envdef:"auth,logging"`
		Debug    bool          `env:"DEBUG" envdef:"false"`
	}

	// Clear any existing environment variables that might interfere
	os.Unsetenv("DEBUG")

	// Test with defaults
	var cfg Config
	if err := envvar.Bind(&cfg); err != nil {
		t.Fatalf("Config should bind successfully: %v", err)
	}

	if cfg.Port != 8080 {
		t.Fatalf("Port should be 8080, got %v", cfg.Port)
	}
	if cfg.Timeout != 30*time.Second {
		t.Fatalf("Timeout should be 30s, got %v", cfg.Timeout)
	}
	if cfg.Name != "myapp" {
		t.Fatalf("Name should be myapp, got %v", cfg.Name)
	}
	if cfg.Mode != "development" {
		t.Fatalf("Mode should be development, got %v", cfg.Mode)
	}
	if len(cfg.Features) != 2 {
		t.Fatalf("Features should have 2 items, got %v", cfg.Features)
	}
	if cfg.Debug {
		t.Fatalf("Debug should be false, got %v", cfg.Debug)
	}
}

// Error handling for missing required fields
func TestErrorHandlingAdvanced(t *testing.T) {
	type Config struct {
		Port int    `env:"PORT,required"`
		Mode string `env:"MODE,required"`
		Name string `env:"NAME,required"`
	}

	// Don't set any environment variables
	var cfg Config
	err := envvar.Bind(&cfg)
	if err == nil {
		t.Fatalf("Should have errors for missing required fields")
	}

	// Error should contain information about missing fields
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "missing") {
		t.Fatalf("Error should mention missing fields: %v", errorMsg)
	}
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			contains(s[1:], substr)))
}
