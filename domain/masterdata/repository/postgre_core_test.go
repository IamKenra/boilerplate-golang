package repository

import (
	"strings"
	"testing"

	"svc-boilerplate-fiber-golang/domain/masterdata"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newDryRunRepository(t *testing.T, dialect string) *postgresMasterdataRepository {
	t.Helper()

	var dialector gorm.Dialector
	switch dialect {
	case "postgres":
		dialector = postgres.New(postgres.Config{
			DSN:                  "host=localhost user=test password=test dbname=test port=5432 sslmode=disable TimeZone=Asia/Jakarta",
			PreferSimpleProtocol: true,
		})
	case "mysql":
		dialector = mysql.New(mysql.Config{
			DSN:                       "test:test@tcp(localhost:3306)/test?parseTime=true&loc=Asia%2FJakarta",
			SkipInitializeWithVersion: true,
		})
	default:
		t.Fatalf("unsupported dialect: %s", dialect)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		DryRun:                 true,
		DisableAutomaticPing:   true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("gorm.Open failed: %v", err)
	}
	return &postgresMasterdataRepository{db: db}
}

func TestMasterdataRepositoryDryRun_GetAll_PostgreSQL(t *testing.T) {
	repo := newDryRunRepository(t, "postgres")

	query, err := repo.applyFilter(repo.db.Model(&masterdataRow{}), masterdata.MasterdataFilter{
		UUID:    "abc",
		OrderBy: "uuid",
		Desc:    true,
		Limit:   10,
		Offset:  20,
	}, true)
	if err != nil {
		t.Fatalf("applyFilter failed: %v", err)
	}

	var rows []masterdataRow
	tx := query.Find(&rows)
	if tx.Error != nil {
		t.Fatalf("Find failed: %v", tx.Error)
	}

	got := tx.Statement.SQL.String()
	for _, want := range []string{
		`SELECT * FROM "masterdata"`,
		`WHERE uuid = $1`,
		`ORDER BY "uuid" DESC`,
		`LIMIT $2 OFFSET $3`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("query missing %q\nquery: %s", want, got)
		}
	}

	wantVars := []interface{}{"abc", 10, 20}
	if len(tx.Statement.Vars) != len(wantVars) {
		t.Fatalf("vars got %v, want %v", tx.Statement.Vars, wantVars)
	}
	for i := range wantVars {
		if tx.Statement.Vars[i] != wantVars[i] {
			t.Fatalf("vars got %v, want %v", tx.Statement.Vars, wantVars)
		}
	}
}

func TestMasterdataRepositoryDryRun_GetAll_MySQL(t *testing.T) {
	repo := newDryRunRepository(t, "mysql")

	query, err := repo.applyFilter(repo.db.Model(&masterdataRow{}), masterdata.MasterdataFilter{
		UUID:    "abc",
		OrderBy: "id",
		Limit:   10,
		Offset:  20,
	}, true)
	if err != nil {
		t.Fatalf("applyFilter failed: %v", err)
	}

	var rows []masterdataRow
	tx := query.Find(&rows)
	if tx.Error != nil {
		t.Fatalf("Find failed: %v", tx.Error)
	}

	got := tx.Statement.SQL.String()
	for _, want := range []string{
		"SELECT * FROM `masterdata`",
		"WHERE uuid = ?",
		"ORDER BY `id`",
		"LIMIT ? OFFSET ?",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("query missing %q\nquery: %s", want, got)
		}
	}

	wantVars := []interface{}{"abc", 10, 20}
	if len(tx.Statement.Vars) != len(wantVars) {
		t.Fatalf("vars got %v, want %v", tx.Statement.Vars, wantVars)
	}
	for i := range wantVars {
		if tx.Statement.Vars[i] != wantVars[i] {
			t.Fatalf("vars got %v, want %v", tx.Statement.Vars, wantVars)
		}
	}
}

func TestMasterdataRepositoryDryRun_RejectsUnsafeOrderBy(t *testing.T) {
	repo := newDryRunRepository(t, "postgres")

	_, err := repo.applyFilter(repo.db.Model(&masterdataRow{}), masterdata.MasterdataFilter{
		OrderBy: "id; DROP TABLE users",
	}, true)
	if err == nil || err.Error() != "invalid order by" {
		t.Fatalf("expected invalid order by, got %v", err)
	}
}
