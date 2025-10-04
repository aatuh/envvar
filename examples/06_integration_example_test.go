package examples

import (
	"strings"
	"testing"
	"time"

	"github.com/aatuh/envvar/v2"
)

// Complete application configuration - real-world usage
func TestCompleteApplicationConfig(t *testing.T) {
	// This example shows how to configure a typical web application
	// using environment variables with all the features of envvar
	// Note: For nested structs, you need to bind them separately or use a flat structure

	type AppConfig struct {
		// Application settings
		AppName     string `env:"APP_NAME,required" validate:"min=1,max=50"`
		Environment string `env:"ENV,required" validate:"oneof=development|staging|production"`
		Debug       bool   `env:"DEBUG" envdef:"false"`

		// Database settings (flattened for easier binding)
		DBHost string `env:"DB_HOST,required"`
		DBPort int    `env:"DB_PORT" envdef:"5432" validate:"min=1,max=65535"`
		DBName string `env:"DB_NAME,required"`
		DBUser string `env:"DB_USER,required"`
		DBPass string `env:"DB_PASS,required"`
		DBSSL  bool   `env:"DB_SSL" envdef:"true"`

		// API settings
		APIPort         int           `env:"API_PORT" envdef:"8080" validate:"min=1,max=65535"`
		APITimeout      time.Duration `env:"API_TIMEOUT" envdef:"30s" validate:"min=1s,max=300s"`
		APIRateLimit    int           `env:"API_RATE_LIMIT" envdef:"100" validate:"min=1,max=10000"`
		APIAllowedHosts []string      `env:"API_ALLOWED_HOSTS" envsep:"," envdef:"localhost,127.0.0.1"`
		APIKey          string        `env:"API_KEY,required"`
		APISecret       string        `env:"API_SECRET,required"`

		// Logging settings
		LogLevel  string `env:"LOG_LEVEL" envdef:"info" validate:"oneof=debug|info|warn|error"`
		LogFormat string `env:"LOG_FORMAT" envdef:"json" validate:"oneof=json|text"`
		LogOutput string `env:"LOG_OUTPUT" envdef:"stdout" validate:"oneof=stdout|stderr|file"`

		// Features and metadata
		Features []string               `env:"FEATURES" envjson:"true" envdef:"[]"`
		Metadata map[string]interface{} `env:"METADATA" envjson:"true" envdef:"{}"`
	}

	// Set up comprehensive environment
	t.Setenv("APP_NAME", "my-awesome-app")
	t.Setenv("ENV", "production")
	t.Setenv("DEBUG", "false")

	// Database configuration
	t.Setenv("DB_HOST", "db.production.com")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_NAME", "myapp_prod")
	t.Setenv("DB_USER", "app_user")
	t.Setenv("DB_PASS", "secure_password")
	t.Setenv("DB_SSL", "true")

	// API configuration
	t.Setenv("API_PORT", "8080")
	t.Setenv("API_TIMEOUT", "30s")
	t.Setenv("API_RATE_LIMIT", "1000")
	t.Setenv("API_ALLOWED_HOSTS", "api.example.com,app.example.com")
	t.Setenv("API_KEY", "ak_live_123456789")
	t.Setenv("API_SECRET", "sk_live_987654321")

	// Logging configuration
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("LOG_FORMAT", "json")
	t.Setenv("LOG_OUTPUT", "stdout")

	// Features as JSON
	t.Setenv("FEATURES", `["auth","payments","notifications","analytics"]`)

	// Metadata as JSON
	t.Setenv("METADATA", `{"version":"1.0.0","region":"us-east-1","deployment":"2024-01-15"}`)

	// Bind configuration
	var cfg AppConfig
	if err := envvar.Bind(&cfg); err != nil {
		t.Fatalf("Configuration binding failed: %v", err)
	}

	// Verify application settings
	if cfg.AppName != "my-awesome-app" {
		t.Fatalf("App name binding failed: %v", cfg.AppName)
	}
	if cfg.Environment != "production" {
		t.Fatalf("Environment binding failed: %v", cfg.Environment)
	}
	if cfg.Debug {
		t.Fatalf("Debug should be false: %v", cfg.Debug)
	}

	// Verify database configuration
	if cfg.DBHost != "db.production.com" {
		t.Fatalf("Database host binding failed: got %q, want %q", cfg.DBHost, "db.production.com")
	}
	if cfg.DBPort != 5432 {
		t.Fatalf("Database port binding failed: %v", cfg.DBPort)
	}
	if !cfg.DBSSL {
		t.Fatalf("Database SSL should be true: %v", cfg.DBSSL)
	}

	// Verify API configuration
	if cfg.APIPort != 8080 {
		t.Fatalf("API port binding failed: %v", cfg.APIPort)
	}
	if cfg.APITimeout != 30*time.Second {
		t.Fatalf("API timeout binding failed: %v", cfg.APITimeout)
	}
	if len(cfg.APIAllowedHosts) != 2 {
		t.Fatalf("API allowed hosts binding failed: %v", cfg.APIAllowedHosts)
	}

	// Verify logging configuration
	if cfg.LogLevel != "info" {
		t.Fatalf("Log level binding failed: %v", cfg.LogLevel)
	}
	if cfg.LogFormat != "json" {
		t.Fatalf("Log format binding failed: %v", cfg.LogFormat)
	}

	// Verify features JSON binding
	if len(cfg.Features) != 4 {
		t.Fatalf("Features binding failed: %v", cfg.Features)
	}

	// Verify metadata JSON binding
	if cfg.Metadata["version"] != "1.0.0" {
		t.Fatalf("Metadata binding failed: %v", cfg.Metadata)
	}
}

// Multi-environment configuration with prefix binding
func TestMultiEnvironmentConfig(t *testing.T) {
	// This example shows how to handle multiple environments
	// using prefix binding for environment-specific overrides

	type Config struct {
		Port     int    `env:"PORT,required" validate:"min=1,max=65535"`
		Database string `env:"DATABASE_URL,required"`
		Debug    bool   `env:"DEBUG"`
		LogLevel string `env:"LOG_LEVEL" envdef:"info" validate:"oneof=debug|info|warn|error"`
	}

	// Set up base configuration
	t.Setenv("PORT", "8080")
	t.Setenv("DATABASE_URL", "postgres://localhost:5432/app")
	t.Setenv("DEBUG", "false")
	t.Setenv("LOG_LEVEL", "info")

	// Set up environment-specific overrides
	t.Setenv("PROD_PORT", "80")
	t.Setenv("PROD_DATABASE_URL", "postgres://prod-db:5432/app")
	t.Setenv("PROD_DEBUG", "false")
	t.Setenv("PROD_LOG_LEVEL", "warn")

	// Test base configuration
	var baseConfig Config
	if err := envvar.Bind(&baseConfig); err != nil {
		t.Fatalf("Base config binding failed: %v", err)
	}
	if baseConfig.Port != 8080 {
		t.Fatalf("Base port should be 8080: %v", baseConfig.Port)
	}

	// Test production configuration with prefix
	var prodConfig Config
	if err := envvar.BindWithPrefix(&prodConfig, "PROD_"); err != nil {
		t.Fatalf("Production config binding failed: %v", err)
	}
	if prodConfig.Port != 80 {
		t.Fatalf("Production port should be 80: %v", prodConfig.Port)
	}
	if prodConfig.Database != "postgres://prod-db:5432/app" {
		t.Fatalf("Production database should be overridden: %v", prodConfig.Database)
	}
}

// Configuration validation with comprehensive error reporting
func TestConfigurationValidation(t *testing.T) {
	type Config struct {
		Port     int           `env:"PORT,required" validate:"min=1,max=65535"`
		Host     string        `env:"HOST,required" validate:"min=1,max=255"`
		Database string        `env:"DATABASE,required" validate:"min=1"`
		Mode     string        `env:"MODE,required" validate:"oneof=dev|staging|prod"`
		Timeout  time.Duration `env:"TIMEOUT,required" validate:"min=1s,max=300s"`
	}

	// Set up invalid configuration
	t.Setenv("PORT", "99999")   // out of range
	t.Setenv("HOST", "")        // too short
	t.Setenv("DATABASE", "")    // too short
	t.Setenv("MODE", "invalid") // not in oneof
	t.Setenv("TIMEOUT", "600s") // too long

	var cfg Config
	err := envvar.Bind(&cfg)
	if err == nil {
		t.Fatalf("Should have validation errors")
	}

	// Error should be a MultiError containing all validation failures
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "multiple errors") {
		t.Fatalf("Should be a MultiError: %v", err)
	}

	// Each field should be mentioned in the error
	fields := []string{"PORT", "HOST", "DATABASE", "MODE", "TIMEOUT"}
	for _, field := range fields {
		if !strings.Contains(errorMsg, field) {
			t.Fatalf("Error should mention %s: %v", field, errorMsg)
		}
	}
}
