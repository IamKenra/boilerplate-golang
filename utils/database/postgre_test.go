package database

import (
	"strings"
	"testing"
)

func setMinimumPostgreEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DB_HOST", "127.0.0.1")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_NAME", "app_db")
	t.Setenv("DB_USERNAME", "postgres")
	t.Setenv("DB_PASSWORD", "postgres")
}

func TestSetupPostgreDatabaseConnection_MissingRequiredEnv(t *testing.T) {
	t.Setenv("DB_DRIVERNAME", "postgres")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_USERNAME", "")
	t.Setenv("DB_PASSWORD", "")

	db, err := SetupPostgreDatabaseConnection()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if db != nil {
		t.Fatal("expected nil db when env is missing")
	}
	if !strings.Contains(err.Error(), "missing required database environment variables") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSetupPostgreDatabaseConnection_InvalidDriver(t *testing.T) {
	setMinimumPostgreEnv(t)
	t.Setenv("DB_DRIVERNAME", "unknown-driver")

	db, err := SetupPostgreDatabaseConnection()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if db != nil {
		t.Fatal("expected nil db on invalid driver")
	}
	if !strings.Contains(err.Error(), "failed to open PostgreSQL connection") {
		t.Fatalf("unexpected error: %v", err)
	}
}
