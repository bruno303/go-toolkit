package config

import (
	"embed"
	"os"
	"testing"
	"time"

	"github.com/bruno303/go-toolkit/pkg/log"
)

//go:embed testdata/*
var testFS embed.FS

type TestConfig struct {
	AppName     string        `yaml:"app_name" env:"APP_NAME"`
	Port        int           `yaml:"port" env:"PORT"`
	Debug       bool          `yaml:"debug" env:"DEBUG"`
	Timeout     time.Duration `yaml:"timeout" env:"TIMEOUT"`
	DatabaseURL string        `yaml:"database_url" env:"DATABASE_URL"`
	Features    []string      `yaml:"features" env:"FEATURES"`
}

type NestedTestConfig struct {
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
}

type ServerConfig struct {
	Host string `yaml:"host" env:"SERVER_HOST"`
	Port int    `yaml:"port" env:"SERVER_PORT"`
}

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     int    `yaml:"port" env:"DB_PORT"`
	Username string `yaml:"username" env:"DB_USERNAME"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
}

func TestLoadConfig_Success(t *testing.T) {
	// Reset log state and configure for testing
	resetLogStateForConfig()

	// Clear any existing environment variables
	clearTestEnvVars()

	// Set explicit config file since we don't have config.yaml at root
	os.Setenv("CONFIG_FILE", "testdata/config.yaml")
	defer os.Unsetenv("CONFIG_FILE")

	var cfg TestConfig
	LoadConfig(&cfg, testFS)

	// Verify the values from the YAML file
	if cfg.AppName != "test-app" {
		t.Errorf("expected app_name to be 'test-app', got %s", cfg.AppName)
	}
	if cfg.Port != 8080 {
		t.Errorf("expected port to be 8080, got %d", cfg.Port)
	}
	if cfg.Debug != true {
		t.Errorf("expected debug to be true, got %v", cfg.Debug)
	}
	if cfg.Timeout != 30*time.Second {
		t.Errorf("expected timeout to be 30s, got %v", cfg.Timeout)
	}
	if cfg.DatabaseURL != "postgres://localhost/testdb" {
		t.Errorf("expected database_url to be 'postgres://localhost/testdb', got %s", cfg.DatabaseURL)
	}
	if len(cfg.Features) != 2 || cfg.Features[0] != "feature1" || cfg.Features[1] != "feature2" {
		t.Errorf("expected features to be ['feature1', 'feature2'], got %v", cfg.Features)
	}
}

func TestLoadConfig_WithEnvironmentOverrides(t *testing.T) {
	// Reset log state and configure for testing
	resetLogStateForConfig()

	// Set explicit config file and environment variables to override YAML values
	os.Setenv("CONFIG_FILE", "testdata/config.yaml")
	os.Setenv("APP_NAME", "env-app")
	os.Setenv("PORT", "9090")
	os.Setenv("DEBUG", "false")
	os.Setenv("TIMEOUT", "60s")
	os.Setenv("DATABASE_URL", "postgres://env-host/envdb")
	defer clearTestEnvVars()

	var cfg TestConfig
	LoadConfig(&cfg, testFS)

	// Verify environment variables took precedence
	if cfg.AppName != "env-app" {
		t.Errorf("expected app_name to be overridden to 'env-app', got %s", cfg.AppName)
	}
	if cfg.Port != 9090 {
		t.Errorf("expected port to be overridden to 9090, got %d", cfg.Port)
	}
	if cfg.Debug != false {
		t.Errorf("expected debug to be overridden to false, got %v", cfg.Debug)
	}
	if cfg.Timeout != 60*time.Second {
		t.Errorf("expected timeout to be overridden to 60s, got %v", cfg.Timeout)
	}
	if cfg.DatabaseURL != "postgres://env-host/envdb" {
		t.Errorf("expected database_url to be overridden to 'postgres://env-host/envdb', got %s", cfg.DatabaseURL)
	}
}

func TestLoadConfig_WithCustomConfigFile(t *testing.T) {
	// Reset log state and configure for testing
	resetLogStateForConfig()

	// Set custom config file
	os.Setenv("CONFIG_FILE", "testdata/custom-config.yaml")
	defer func() {
		os.Unsetenv("CONFIG_FILE")
		clearTestEnvVars()
	}()

	var cfg TestConfig
	LoadConfig(&cfg, testFS)

	// Verify values from custom config file
	if cfg.AppName != "custom-app" {
		t.Errorf("expected app_name to be 'custom-app', got %s", cfg.AppName)
	}
	if cfg.Port != 3000 {
		t.Errorf("expected port to be 3000, got %d", cfg.Port)
	}
}

func TestLoadConfig_NestedStructure(t *testing.T) {
	// Reset log state and configure for testing
	resetLogStateForConfig()

	// Clear environment variables
	clearTestEnvVars()

	// Set custom config file for nested structure
	os.Setenv("CONFIG_FILE", "testdata/nested-config.yaml")
	defer func() {
		os.Unsetenv("CONFIG_FILE")
	}()

	var cfg NestedTestConfig
	LoadConfig(&cfg, testFS)

	// Verify nested structure values
	if cfg.Server.Host != "localhost" {
		t.Errorf("expected server.host to be 'localhost', got %s", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected server.port to be 8080, got %d", cfg.Server.Port)
	}
	if cfg.DB.Host != "db-host" {
		t.Errorf("expected db.host to be 'db-host', got %s", cfg.DB.Host)
	}
	if cfg.DB.Port != 5432 {
		t.Errorf("expected db.port to be 5432, got %d", cfg.DB.Port)
	}
	if cfg.DB.Username != "testuser" {
		t.Errorf("expected db.username to be 'testuser', got %s", cfg.DB.Username)
	}
	if cfg.DB.Password != "testpass" {
		t.Errorf("expected db.password to be 'testpass', got %s", cfg.DB.Password)
	}
}

func TestLoadConfig_NestedWithEnvironmentOverrides(t *testing.T) {
	// Reset log state and configure for testing
	resetLogStateForConfig()

	// Set environment variables for nested config
	os.Setenv("CONFIG_FILE", "testdata/nested-config.yaml")
	os.Setenv("SERVER_HOST", "env-server")
	os.Setenv("SERVER_PORT", "9000")
	os.Setenv("DB_HOST", "env-db")
	os.Setenv("DB_PASSWORD", "env-password")
	defer func() {
		os.Unsetenv("CONFIG_FILE")
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PASSWORD")
	}()

	var cfg NestedTestConfig
	LoadConfig(&cfg, testFS)

	// Verify environment overrides
	if cfg.Server.Host != "env-server" {
		t.Errorf("expected server.host to be overridden to 'env-server', got %s", cfg.Server.Host)
	}
	if cfg.Server.Port != 9000 {
		t.Errorf("expected server.port to be overridden to 9000, got %d", cfg.Server.Port)
	}
	if cfg.DB.Host != "env-db" {
		t.Errorf("expected db.host to be overridden to 'env-db', got %s", cfg.DB.Host)
	}
	if cfg.DB.Password != "env-password" {
		t.Errorf("expected db.password to be overridden to 'env-password', got %s", cfg.DB.Password)
	}
	// Values not overridden should remain from YAML
	if cfg.DB.Username != "testuser" {
		t.Errorf("expected db.username to remain 'testuser', got %s", cfg.DB.Username)
	}
}

func TestLoadConfig_FileNotFound_Panics(t *testing.T) {
	// Reset log state and configure for testing
	resetLogStateForConfig()

	// Set non-existent config file
	os.Setenv("CONFIG_FILE", "non-existent.yaml")
	defer os.Unsetenv("CONFIG_FILE")

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected LoadConfig to panic when config file is not found")
		}
	}()

	var cfg TestConfig
	LoadConfig(&cfg, testFS)
}

func TestLoadConfig_InvalidYAML_Panics(t *testing.T) {
	// Reset log state and configure for testing
	resetLogStateForConfig()

	// Set config file with invalid YAML
	os.Setenv("CONFIG_FILE", "testdata/invalid.yaml")
	defer os.Unsetenv("CONFIG_FILE")

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected LoadConfig to panic when YAML is invalid")
		}
	}()

	var cfg TestConfig
	LoadConfig(&cfg, testFS)
}

func TestLoadConfig_DefaultConfigFileName(t *testing.T) {
	// Reset log state and configure for testing
	resetLogStateForConfig()

	// Ensure CONFIG_FILE is not set - this should try to open "config.yaml" which doesn't exist at root
	os.Unsetenv("CONFIG_FILE")
	clearTestEnvVars()

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected LoadConfig to panic when default config.yaml is not found at root")
		}
	}()

	// This should use the default "config.yaml" file (which doesn't exist in our test FS)
	var cfg TestConfig
	LoadConfig(&cfg, testFS)
}

// Helper functions

func clearTestEnvVars() {
	os.Unsetenv("CONFIG_FILE")
	os.Unsetenv("APP_NAME")
	os.Unsetenv("PORT")
	os.Unsetenv("DEBUG")
	os.Unsetenv("TIMEOUT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("FEATURES")
	os.Unsetenv("SERVER_HOST")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USERNAME")
	os.Unsetenv("DB_PASSWORD")
}

func TestLoadConfig_MinimalConfigFile(t *testing.T) {
	// Reset log state and configure for testing
	resetLogStateForConfig()

	// Set config file with minimal content (empty object)
	os.Setenv("CONFIG_FILE", "testdata/minimal.yaml")
	defer os.Unsetenv("CONFIG_FILE")

	var cfg TestConfig
	LoadConfig(&cfg, testFS)

	// Verify default values (zero values for the struct)
	if cfg.AppName != "" {
		t.Errorf("expected app_name to be empty, got %s", cfg.AppName)
	}
	if cfg.Port != 0 {
		t.Errorf("expected port to be 0, got %d", cfg.Port)
	}
	if cfg.Debug != false {
		t.Errorf("expected debug to be false, got %v", cfg.Debug)
	}
	if cfg.Timeout != 0 {
		t.Errorf("expected timeout to be 0, got %v", cfg.Timeout)
	}
}

func TestLoadConfig_EmptyConfigFile_Panics(t *testing.T) {
	// Reset log state and configure for testing
	resetLogStateForConfig()

	// Set config file with truly empty content (YAML decoder treats EOF as error)
	os.Setenv("CONFIG_FILE", "testdata/empty.yaml")
	defer os.Unsetenv("CONFIG_FILE")

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected LoadConfig to panic when config file is completely empty (EOF)")
		}
	}()

	var cfg TestConfig
	LoadConfig(&cfg, testFS)
}

func TestLoadConfig_PartialEnvironmentOverride(t *testing.T) {
	// Reset log state and configure for testing
	resetLogStateForConfig()

	// Set explicit config file and only some environment variables
	os.Setenv("CONFIG_FILE", "testdata/config.yaml")
	os.Setenv("APP_NAME", "partial-env")
	os.Setenv("PORT", "5000")
	// Don't set DEBUG, TIMEOUT, DATABASE_URL - they should use YAML values
	defer func() {
		os.Unsetenv("CONFIG_FILE")
		os.Unsetenv("APP_NAME")
		os.Unsetenv("PORT")
	}()

	var cfg TestConfig
	LoadConfig(&cfg, testFS)

	// Verify partial overrides
	if cfg.AppName != "partial-env" {
		t.Errorf("expected app_name to be overridden to 'partial-env', got %s", cfg.AppName)
	}
	if cfg.Port != 5000 {
		t.Errorf("expected port to be overridden to 5000, got %d", cfg.Port)
	}
	// These should remain from YAML
	if cfg.Debug != true {
		t.Errorf("expected debug to remain true from YAML, got %v", cfg.Debug)
	}
	if cfg.Timeout != 30*time.Second {
		t.Errorf("expected timeout to remain 30s from YAML, got %v", cfg.Timeout)
	}
	if cfg.DatabaseURL != "postgres://localhost/testdb" {
		t.Errorf("expected database_url to remain from YAML, got %s", cfg.DatabaseURL)
	}
}

func TestLoadConfig_TypeConversionErrors(t *testing.T) {
	// Reset log state and configure for testing
	resetLogStateForConfig()

	// Set explicit config file and environment variables with invalid types that should cause envconfig to fail
	os.Setenv("CONFIG_FILE", "testdata/config.yaml")
	os.Setenv("PORT", "not-a-number")
	os.Setenv("TIMEOUT", "invalid-duration")
	defer func() {
		os.Unsetenv("CONFIG_FILE")
		os.Unsetenv("PORT")
		os.Unsetenv("TIMEOUT")
	}()

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected LoadConfig to panic when environment variables have invalid types")
		}
	}()

	var cfg TestConfig
	LoadConfig(&cfg, testFS)
}

func TestLoadConfig_ComplexEnvironmentVariables(t *testing.T) {
	// Reset log state and configure for testing
	resetLogStateForConfig()

	// Set explicit config file and environment variables with complex values
	os.Setenv("CONFIG_FILE", "testdata/config.yaml")
	os.Setenv("FEATURES", "env-feature1,env-feature2,env-feature3")
	defer func() {
		os.Unsetenv("CONFIG_FILE")
		os.Unsetenv("FEATURES")
	}()

	var cfg TestConfig
	LoadConfig(&cfg, testFS)

	// Verify array parsing from environment
	expectedFeatures := []string{"env-feature1", "env-feature2", "env-feature3"}
	if len(cfg.Features) != len(expectedFeatures) {
		t.Errorf("expected %d features, got %d", len(expectedFeatures), len(cfg.Features))
	}
	for i, expected := range expectedFeatures {
		if i >= len(cfg.Features) || cfg.Features[i] != expected {
			t.Errorf("expected feature[%d] to be %s, got %v", i, expected, cfg.Features)
		}
	}
}

func resetLogStateForConfig() {
	// Configure a simple logger for testing
	log.SetLogger(log.NewSlogAdapter(log.SlogAdapterOpts{
		Level:      log.LevelDebug,
		FormatJson: false,
		Name:       "config-test",
	}))
}
