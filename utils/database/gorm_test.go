package database

import (
	"strings"
	"testing"
)

func TestSetupDatabaseConnection_MissingRequiredEnv(t *testing.T) {
	for _, dialect := range []string{"postgres", "mysql"} {
		t.Run(dialect, func(t *testing.T) {
			t.Setenv("DB_DIALECT", dialect)
			for _, key := range []string{"DB_HOST", "DB_PORT", "DB_NAME", "DB_USERNAME", "DB_PASSWORD"} {
				t.Setenv(key, "")
			}

			db, err := SetupDatabaseConnection()
			if db != nil {
				t.Fatal("expected nil db when env is missing")
			}
			want := "missing required database environment variables: DB_HOST, DB_NAME, DB_PORT, DB_USERNAME"
			if err == nil || err.Error() != want {
				t.Fatalf("got %v, want %q", err, want)
			}
		})
	}
}

func TestSetupDatabaseConnection_UnsupportedDialect(t *testing.T) {
	t.Setenv("DB_DIALECT", "sqlite")

	db, err := SetupDatabaseConnection()
	if db != nil {
		t.Fatal("expected nil db on unsupported dialect")
	}
	if err == nil || !strings.Contains(err.Error(), "unsupported DB_DIALECT: sqlite") {
		t.Fatalf("unexpected error: %v", err)
	}
}
