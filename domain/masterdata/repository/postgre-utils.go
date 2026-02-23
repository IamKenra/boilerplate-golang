package repository

import (
	"log"
	"boilerplate-golang/utils/database"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		log.Fatalf("Failed to create snowflake node: %v", err)
	}
}
func (db *postgresMasterdataRepository) GenerateID() (uint64, error) {
	return uint64(node.Generate().Int64()), nil
}

func (db *postgresMasterdataRepository) Health() error {
	return db.sqlDB.Ping()
}

func (db *postgresMasterdataRepository) GenerateUUID() (string, error) {
	return database.GenerateUUIDPostgre(db.sqlDB)
}

func (db *postgresMasterdataRepository) Exec(queryConfig ...database.QueryConfig) (err error) {
	return database.ExecTransaction(db.sqlDB, queryConfig...)
}
