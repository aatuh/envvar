package examples

import (
	"strconv"
	"strings"
	"testing"

	"github.com/aatuh/envvar/v2"
)

// Advanced getters - URL, IP, and string slices
func TestAdvancedGetters(t *testing.T) {
	t.Setenv("DATABASE_URL", "https://user:pass@db.example.com:5432/mydb")
	t.Setenv("API_HOST", "192.168.1.100")
	t.Setenv("ALLOWED_HOSTS", "localhost,127.0.0.1,example.com")
	t.Setenv("TAGS", "go,web,api")

	// URL getter
	if u, err := envvar.GetURL("DATABASE_URL"); err != nil {
		t.Fatalf("GetURL failed: %v", err)
	} else if u.Host != "db.example.com:5432" {
		t.Fatalf("GetURL host wrong: %v", u.Host)
	}

	// IP getter
	if ip, err := envvar.GetIP("API_HOST"); err != nil {
		t.Fatalf("GetIP failed: %v", err)
	} else if ip.String() != "192.168.1.100" {
		t.Fatalf("GetIP wrong: %v", ip)
	}

	// String slice getter (comma-separated by default)
	if hosts, err := envvar.GetStringSlice("ALLOWED_HOSTS"); err != nil {
		t.Fatalf("GetStringSlice failed: %v", err)
	} else if len(hosts) != 3 {
		t.Fatalf("GetStringSlice length wrong: %v", hosts)
	}

	// String slice with custom separator
	if tags, err := envvar.GetStringSliceSep("TAGS", ","); err != nil {
		t.Fatalf("GetStringSliceSep failed: %v", err)
	} else if len(tags) != 3 {
		t.Fatalf("GetStringSliceSep length wrong: %v", tags)
	}
}

// Generic typed getters - custom conversion functions
func TestTypedGetters(t *testing.T) {
	t.Setenv("UPPER_TEXT", "hello world")
	t.Setenv("JSON_SIZE", "1024")

	// Custom string transformation
	upper, err := envvar.GetTyped("UPPER_TEXT", func(s string) (string, error) {
		return strings.ToUpper(s), nil
	})
	if err != nil || upper != "HELLO WORLD" {
		t.Fatalf("GetTyped string transformation failed: %v, %v", upper, err)
	}

	// Custom integer parsing with validation
	size, err := envvar.GetTyped("JSON_SIZE", func(s string) (int, error) {
		// Custom validation logic here
		return strconv.Atoi(s)
	})
	if err != nil || size != 1024 {
		t.Fatalf("GetTyped integer parsing failed: %v, %v", size, err)
	}

	// MustGetTyped panics on error
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("MustGetTyped should panic on error")
		}
	}()
	envvar.MustGetTyped("INVALID_JSON_SIZE", func(s string) (int, error) {
		return strconv.Atoi(s)
	})
}

// Variable expansion - ${VAR} and ${VAR:-default}
func TestVariableExpansion(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_NAME", "mydb")
	t.Setenv("DATABASE_URL", "postgres://${DB_HOST}:${DB_PORT:-5432}/${DB_NAME}")
	t.Setenv("FALLBACK_URL", "postgres://${MISSING_HOST:-localhost}:${MISSING_PORT:-5432}/db")

	// Expansion with existing variables
	if url, ok := envvar.Get("DATABASE_URL"); !ok || url != "postgres://localhost:5432/mydb" {
		t.Fatalf("Variable expansion failed: %v", url)
	}

	// Expansion with fallback defaults
	if url, ok := envvar.Get("FALLBACK_URL"); !ok || url != "postgres://localhost:5432/db" {
		t.Fatalf("Variable expansion with defaults failed: %v", url)
	}
}

// Lazy getters - cache values on first access
func TestLazyGetters(t *testing.T) {
	t.Setenv("EXPENSIVE_CONFIG", "complex-value")
	t.Setenv("DEBUG_MODE", "true")

	// Lazy getters return functions that cache the value
	lazyConfig := envvar.LazyString("EXPENSIVE_CONFIG")
	lazyDebug := envvar.LazyBool("DEBUG_MODE")

	// First call evaluates and caches
	config1 := lazyConfig()
	debug1 := lazyDebug()

	// Second call returns cached value (no re-evaluation)
	config2 := lazyConfig()
	debug2 := lazyDebug()

	if config1 != config2 || config1 != "complex-value" {
		t.Fatalf("LazyString caching failed: %v, %v", config1, config2)
	}

	if debug1 != debug2 || !debug1 {
		t.Fatalf("LazyBool caching failed: %v, %v", debug1, debug2)
	}
}
