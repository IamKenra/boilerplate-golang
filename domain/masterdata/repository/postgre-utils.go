package repository

import (
	"log"
	"os"
	"strconv"
	"svc-boilerplate-fiber-golang/utils/database"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

var node *snowflake.Node

func init() {
	var err error
	nodeID := int64(1)
	if raw := os.Getenv("NODE_ID"); raw != "" {
		nodeID, err = strconv.ParseInt(raw, 10, 64)
		if err != nil || nodeID < 0 || nodeID > 1023 {
			log.Fatalf("invalid NODE_ID: %q", raw)
		}
	}
	node, err = snowflake.NewNode(nodeID)
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
	return uuid.NewString(), nil
}

func (db *postgresMasterdataRepository) Exec(queryConfig ...database.QueryConfig) (err error) {
	return database.ExecTransaction(db.sqlDB, queryConfig...)
}
