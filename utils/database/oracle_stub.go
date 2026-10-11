//go:build !oracle

package database

import (
	"database/sql"
	"errors"
)

func SetupOracleDatabaseConnection() (*sql.DB, error) {
	return nil, errors.New("oracle support disabled; build with -tags oracle")
}
