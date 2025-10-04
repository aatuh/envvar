package examples

import (
	"testing"
	"time"

	"github.com/aatuh/envvar/v2"
)

// Basic getters - the simplest way to read environment variables
func TestBasicGetters(t *testing.T) {
	// Set up test environment
	t.Setenv("APP_NAME", "myapp")
	t.Setenv("DEBUG", "true")
	t.Setenv("PORT", "8080")
	t.Setenv("RATE", "3.14")
	t.Setenv("TIMEOUT", "5s")

	// Basic string getter
	if name, ok := envvar.Get("APP_NAME"); !ok || name != "myapp" {
		t.Fatalf("Get failed: %v, %v", name, ok)
	}

	// String with default
	if name := envvar.GetOr("MISSING_VAR", "default"); name != "default" {
		t.Fatalf("GetOr failed: %v", name)
	}

	// Boolean getter
	if debug, err := envvar.GetBool("DEBUG"); err != nil || !debug {
		t.Fatalf("GetBool failed: %v, %v", debug, err)
	}

	// Integer getter
	if port, err := envvar.GetInt("PORT"); err != nil || port != 8080 {
		t.Fatalf("GetInt failed: %v, %v", port, err)
	}

	// Float getter
	if rate, err := envvar.GetFloat64("RATE"); err != nil || rate < 3.139 || rate > 3.141 {
		t.Fatalf("GetFloat64 failed: %v, %v", rate, err)
	}

	// Duration getter
	if timeout, err := envvar.GetDuration("TIMEOUT"); err != nil || timeout != 5*time.Second {
		t.Fatalf("GetDuration failed: %v, %v", timeout, err)
	}
}

// Must getters - panic on missing values (use for required config)
func TestMustGetters(t *testing.T) {
	t.Setenv("REQUIRED_PORT", "9090")

	// MustGet panics if value is missing
	port := envvar.MustGet("REQUIRED_PORT")
	if port != "9090" {
		t.Fatalf("MustGet failed: %v", port)
	}

	// MustGetInt panics if value is missing or invalid
	portInt := envvar.MustGetInt("REQUIRED_PORT")
	if portInt != 9090 {
		t.Fatalf("MustGetInt failed: %v", portInt)
	}

	// Test panic on missing value
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("MustGet should panic on missing value")
		}
	}()
	envvar.MustGet("MISSING_REQUIRED")
}

// GetOr variants - provide defaults for optional values
func TestGetOrVariants(t *testing.T) {
	t.Setenv("SET_PORT", "8080")
	t.Setenv("SET_DEBUG", "true")

	// GetOr for strings
	if port := envvar.GetOr("SET_PORT", "3000"); port != "8080" {
		t.Fatalf("GetOr with existing value failed: %v", port)
	}
	if port := envvar.GetOr("MISSING_PORT", "3000"); port != "3000" {
		t.Fatalf("GetOr with default failed: %v", port)
	}

	// GetIntOr for integers
	if port := envvar.GetIntOr("SET_PORT", 3000); port != 8080 {
		t.Fatalf("GetIntOr with existing value failed: %v", port)
	}
	if port := envvar.GetIntOr("MISSING_PORT", 3000); port != 3000 {
		t.Fatalf("GetIntOr with default failed: %v", port)
	}

	// GetBoolOr for booleans
	if debug := envvar.GetBoolOr("SET_DEBUG", false); !debug {
		t.Fatalf("GetBoolOr with existing value failed: %v", debug)
	}
	if debug := envvar.GetBoolOr("MISSING_DEBUG", true); !debug {
		t.Fatalf("GetBoolOr with default failed: %v", debug)
	}
}

// Error handling for invalid values
func TestErrorHandling(t *testing.T) {
	t.Setenv("BAD_INT", "not-a-number")
	t.Setenv("BAD_FLOAT", "not-a-float")
	t.Setenv("BAD_DURATION", "not-a-duration")

	// GetInt returns error for invalid values
	if _, err := envvar.GetInt("BAD_INT"); err == nil {
		t.Fatalf("GetInt should return error for invalid value")
	}

	// GetFloat64 returns error for invalid values
	if _, err := envvar.GetFloat64("BAD_FLOAT"); err == nil {
		t.Fatalf("GetFloat64 should return error for invalid value")
	}

	// GetDuration returns error for invalid values
	if _, err := envvar.GetDuration("BAD_DURATION"); err == nil {
		t.Fatalf("GetDuration should return error for invalid value")
	}

	// GetOrErr returns error for missing values
	if _, err := envvar.GetOrErr("MISSING_VAR"); err == nil {
		t.Fatalf("GetOrErr should return error for missing value")
	}
}
