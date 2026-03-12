package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func SetupMysqlDatabaseConnection() (db *sql.DB, err error) {
	var (
		driver   = strings.TrimSpace(os.Getenv("DB_DRIVERNAME"))
		username = os.Getenv("DB_USERNAME")
		password = os.Getenv("DB_PASSWORD")
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		name     = os.Getenv("DB_NAME")
		location = "loc=Asia%2FJakarta"
	)

	if driver == "" {
		driver = "mysql"
	}
	if missing := findMissingEnvValues(map[string]string{
		"DB_HOST":     host,
		"DB_NAME":     name,
		"DB_PORT":     port,
		"DB_USERNAME": username,
	}); len(missing) > 0 {
		return nil, fmt.Errorf("missing required database environment variables: %s", strings.Join(missing, ", "))
	}

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&%s", username, password, host, port, name, location)

	db, err = sql.Open(driver, connection)
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL connection: %w", err)
	}
	if db == nil {
		return nil, errors.New("failed to open MySQL connection: nil database handle")
	}

	db.SetMaxOpenConns(envInt("DB_MAX_OPEN_CONNS", 25))
	db.SetMaxIdleConns(envInt("DB_MAX_IDLE_CONNS", 25))
	db.SetConnMaxLifetime(envDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute))
	db.SetConnMaxIdleTime(envDuration("DB_CONN_MAX_IDLE_TIME", 10*time.Minute))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping MySQL connection: %w", err)
	}

	return
}

func GenerateID(db *sql.DB) (id uint64, err error) {
	query := db.QueryRow(`SELECT UUID_SHORT()`)
	err = query.Scan(&id)
	return
}

func GenerateUUID(db *sql.DB) (uuid string, err error) {
	query := db.QueryRow(`SELECT UUID()`)
	err = query.Scan(&uuid)
	return
}
