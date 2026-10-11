package repository

import (
	"database/sql"
	"svc-boilerplate-fiber-golang/domain/masterdata"
)

const POSTGRE = "postgre"

const SELECT_DISTINCT = "select-distinct"
const SELECT = "select"
const INSERT = "insert"
const UPDATE = "update"
const DELETE = "delete"

// example for database configuration
const POSTGRE_TABLE = "masterdata"

type postgresMasterdataRepository struct {
	sqlDB *sql.DB
}

func NewPostgresRepository(db *sql.DB) masterdata.PostgreRepository {
	return &postgresMasterdataRepository{db}
}
