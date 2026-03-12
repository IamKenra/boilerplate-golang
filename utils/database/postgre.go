package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func SetupPostgreDatabaseConnection() (db *sql.DB, err error) {
	var (
		driver   = strings.TrimSpace(os.Getenv("DB_DRIVERNAME"))
		username = os.Getenv("DB_USERNAME")
		password = os.Getenv("DB_PASSWORD")
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		name     = os.Getenv("DB_NAME")
		location = "TimeZone=Asia/Jakarta"
	)

	if driver == "" {
		driver = "postgres"
	}
	if missing := findMissingEnvValues(map[string]string{
		"DB_HOST":     host,
		"DB_NAME":     name,
		"DB_PORT":     port,
		"DB_USERNAME": username,
	}); len(missing) > 0 {
		return nil, fmt.Errorf("missing required database environment variables: %s", strings.Join(missing, ", "))
	}

	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&%s", username, password, host, port, name, location)

	db, err = sql.Open(driver, connection)
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}
	if db == nil {
		return nil, errors.New("failed to open PostgreSQL connection: nil database handle")
	}

	db.SetMaxOpenConns(envInt("DB_MAX_OPEN_CONNS", 25))
	db.SetMaxIdleConns(envInt("DB_MAX_IDLE_CONNS", 25))
	db.SetConnMaxLifetime(envDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute))
	db.SetConnMaxIdleTime(envDuration("DB_CONN_MAX_IDLE_TIME", 10*time.Minute))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping PostgreSQL connection: %w", err)
	}

	return
}

func GenerateIDPostgre(db *sql.DB, sequenceName string) (id uint64, err error) {
	query := db.QueryRow(fmt.Sprintf("SELECT NEXTVAL('%s')", sequenceName))
	err = query.Scan(&id)
	return
}

func GenerateUUIDPostgre(db *sql.DB) (uuid string, err error) {
	query := db.QueryRow(`SELECT gen_random_uuid()`)
	err = query.Scan(&uuid)
	return
}
