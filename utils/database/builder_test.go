package database

import (
	"reflect"
	"strings"
	"testing"
)

func TestQueryBuilder_PostgreSelectWhereLimit(t *testing.T) {
	cfg := New("postgre", "masterdata", "select")
	cfg.OnSelect = OnSelect{
		Column: []string{"id", "uuid"},
		Where: map[string]interface{}{
			"AND":   map[string]interface{}{"uuid": "abc"},
			"LIMIT": []interface{}{20, 10}, // offset, limit
		},
	}

	if err := cfg.QueryBuilder(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Builder emits double spaces around WHERE; compare normalized whitespace.
	got := strings.Join(strings.Fields(cfg.Result.Query), " ")
	want := "SELECT id, uuid FROM masterdata WHERE uuid = $1 LIMIT $2 OFFSET $3"
	if got != want {
		t.Fatalf("query\n got: %q\nwant: %q", got, want)
	}

	wantValues := []interface{}{"abc", 10, 20}
	if !reflect.DeepEqual(cfg.Result.Value, wantValues) {
		t.Fatalf("values got %v, want %v", cfg.Result.Value, wantValues)
	}
}

func TestQueryBuilder_UpdateWithoutWhereFails(t *testing.T) {
	cfg := New("postgre", "masterdata", "update")
	cfg.OnUpdate = OnUpdate{
		Data: map[string]interface{}{"user_update": "tester"},
	}

	err := cfg.QueryBuilder()
	if err == nil || err.Error() != "unsafe update with no where" {
		t.Fatalf("expected unsafe update error, got %v (query: %q)", err, cfg.Result.Query)
	}
}

func TestQueryBuilder_PostgreInsertPlaceholders(t *testing.T) {
	cfg := New("postgre", "masterdata", "insert")
	cfg.OnInsert = OnInsert{
		Column: []string{"id", "uuid", "user_input"},
		Data:   []interface{}{[]interface{}{uint64(1), "abc", "tester"}},
	}

	if err := cfg.QueryBuilder(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "INSERT INTO masterdata (id, uuid, user_input) VALUES ($1, $2, $3)"
	if cfg.Result.Query != want {
		t.Fatalf("query\n got: %q\nwant: %q", cfg.Result.Query, want)
	}

	wantValues := []interface{}{uint64(1), "abc", "tester"}
	if !reflect.DeepEqual(cfg.Result.Value, wantValues) {
		t.Fatalf("values got %v, want %v", cfg.Result.Value, wantValues)
	}
}
