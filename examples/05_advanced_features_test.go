package examples

import (
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

// Complex validation rules
func TestComplexValidation(t *testing.T) {
	type Config struct {
		// Numeric validation
		Port    int           `env:"PORT" validate:"min=1,max=65535"`
		Timeout time.Duration `env:"TIMEOUT" validate:"min=1s,max=60s"`

		// String length validation
		Name string `env:"NAME" validate:"min=3,max=50"`

		// Oneof validation for strings
		Mode string `env:"MODE" validate:"oneof=dev|staging|prod"`

		// Oneof validation for slices
		Features []string `env:"FEATURES" validate:"oneof=auth|logging|metrics"`
	}

	// Valid configuration
	t.Setenv("PORT", "8080")
	t.Setenv("TIMEOUT", "30s")
	t.Setenv("NAME", "myapp")
	t.Setenv("MODE", "prod")
	t.Setenv("FEATURES", "auth,logging")

	var cfg Config
	if err := envvar.Bind(&cfg); err != nil {
		t.Fatalf("Valid config should pass: %v", err)
	}

	// Test port validation
	t.Setenv("PORT", "99999")
	if err := envvar.Bind(&cfg); err == nil {
		t.Fatalf("Invalid port should fail validation")
	}

	// Test timeout validation
	t.Setenv("PORT", "8080")
	t.Setenv("TIMEOUT", "120s")
	if err := envvar.Bind(&cfg); err == nil {
		t.Fatalf("Invalid timeout should fail validation")
	}

	// Test name length validation
	t.Setenv("TIMEOUT", "30s")
	t.Setenv("NAME", "ab") // too short
	if err := envvar.Bind(&cfg); err == nil {
		t.Fatalf("Invalid name length should fail validation")
	}

	// Test mode oneof validation
	t.Setenv("NAME", "myapp")
	t.Setenv("MODE", "invalid")
	if err := envvar.Bind(&cfg); err == nil {
		t.Fatalf("Invalid mode should fail validation")
	}

	// Test features oneof validation
	t.Setenv("MODE", "prod")
	t.Setenv("FEATURES", "auth,invalid")
	if err := envvar.Bind(&cfg); err == nil {
		t.Fatalf("Invalid features should fail validation")
	}
}

// Error aggregation - multiple validation errors
func TestErrorAggregation(t *testing.T) {
	type Config struct {
		Port int    `env:"PORT,required" validate:"min=1,max=65535"`
		Mode string `env:"MODE,required" validate:"oneof=dev|prod"`
		Name string `env:"NAME,required" validate:"min=3"`
	}

	// Set up invalid values for all fields
	t.Setenv("PORT", "99999")   // out of range
	t.Setenv("MODE", "invalid") // not in oneof
	t.Setenv("NAME", "ab")      // too short

	var cfg Config
	err := envvar.Bind(&cfg)
	if err == nil {
		t.Fatalf("Should have validation errors")
	}

	// Error should contain information about all validation failures
	errorMsg := err.Error()
	if !contains(errorMsg, "PORT") || !contains(errorMsg, "MODE") || !contains(errorMsg, "NAME") {
		t.Fatalf("Error should mention all failed fields: %v", errorMsg)
	}
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			contains(s[1:], substr)))
}
