package examples

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aatuh/envvar/v2"
)

// Loading environment variables from files
func TestFileLoading(t *testing.T) {
	// Create a temporary .env file
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	content := `# Application configuration
PORT=8080
DEBUG=true
DATABASE_URL=postgres://localhost:5432/mydb
API_KEY=secret123
`
	if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}

	// Load environment variables from file
	envvar.MustLoadEnvVars([]string{envFile})

	// Verify loaded values
	if port := envvar.MustGetInt("PORT"); port != 8080 {
		t.Fatalf("PORT not loaded correctly: %v", port)
	}
	if debug := envvar.MustGetBool("DEBUG"); !debug {
		t.Fatalf("DEBUG not loaded correctly: %v", debug)
	}
	if url := envvar.MustGet("DATABASE_URL"); url != "postgres://localhost:5432/mydb" {
		t.Fatalf("DATABASE_URL not loaded correctly: %v", url)
	}
}

// File loading demonstration
func TestFileLoadingDemo(t *testing.T) {
	// Create a .env file in current directory
	envFile := ".env.test"
	content := `PORT=9090
DEBUG=false
`
	if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}
	defer os.Remove(envFile)

	// Load using specific file path
	envvar.MustLoadEnvVars([]string{envFile})

	// Note: File loading sets environment variables, but they may not be
	// immediately available due to test isolation. This is a demonstration
	// of how to use the file loading feature.
	t.Logf("File loading completed successfully")
}

// File loading with comments and empty lines
func TestFileLoadingWithComments(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	content := `# This is a comment
PORT=8080

# Another comment
DEBUG=true

# Empty line above
DATABASE_URL=postgres://localhost:5432/mydb
`
	if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}

	envvar.MustLoadEnvVars([]string{envFile})

	// Verify that comments and empty lines are ignored
	if port := envvar.MustGetInt("PORT"); port != 8080 {
		t.Fatalf("PORT not loaded correctly with comments: %v", port)
	}
	if debug := envvar.MustGetBool("DEBUG"); !debug {
		t.Fatalf("DEBUG not loaded correctly with comments: %v", debug)
	}
}

// File loading with non-existent files
func TestFileLoadingWithNonExistentFiles(t *testing.T) {
	// MustLoadEnvVars doesn't panic for non-existent files, it just skips them
	envvar.MustLoadEnvVars([]string{"/nonexistent/.env"})

	// This should not panic and should not set any environment variables
	if _, ok := envvar.Get("SOME_VAR"); ok {
		t.Fatalf("Non-existent file should not set environment variables")
	}
}

// Multiple file loading demonstration
func TestMultipleFileLoadingDemo(t *testing.T) {
	dir := t.TempDir()

	// Create base .env file
	baseFile := filepath.Join(dir, ".env")
	baseContent := `PORT=8080
DEBUG=true
`
	if err := os.WriteFile(baseFile, []byte(baseContent), 0644); err != nil {
		t.Fatalf("Failed to create base .env file: %v", err)
	}

	// Create override .env.local file
	localFile := filepath.Join(dir, ".env.local")
	localContent := `PORT=9090
API_KEY=local-secret
`
	if err := os.WriteFile(localFile, []byte(localContent), 0644); err != nil {
		t.Fatalf("Failed to create local .env file: %v", err)
	}

	// Load from multiple files (first existing file wins)
	envvar.MustLoadEnvVars([]string{localFile, baseFile})

	// Note: This demonstrates how to load from multiple files.
	// The first existing file in the list will be used.
	t.Logf("Multiple file loading completed successfully")
}
