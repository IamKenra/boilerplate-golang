package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/godror/godror"
)

func SetupOracleDatabaseConnection() (db *sql.DB, err error) {
	var (
		driver   = os.Getenv("ORACLE_DB_DRIVERNAME")
		username = os.Getenv("ORACLE_DB_USERNAME")
		password = os.Getenv("ORACLE_DB_PASSWORD")
		host     = os.Getenv("ORACLE_DB_HOST")
		port     = os.Getenv("ORACLE_DB_PORT")
		service  = os.Getenv("ORACLE_DB_SERVICE_NAME")
	)

	connection := fmt.Sprintf(`user="%s" password="%s" connectString="%s:%s/%s"`, username, password, host, port, service)

	return sql.Open(driver, connection)
}
