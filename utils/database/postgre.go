package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func SetupPostgreDatabaseConnection() (db *sql.DB, err error) {
	var (
		driver   = os.Getenv("DB_DRIVERNAME")
		username = os.Getenv("DB_USERNAME")
		password = os.Getenv("DB_PASSWORD")
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		name     = os.Getenv("DB_NAME")
		location = "TimeZone=Asia/Jakarta"
	)

	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&%s", username, password, host, port, name, location)

	db, err = sql.Open(driver, connection)

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	db.SetConnMaxLifetime(100 * time.Millisecond)

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
