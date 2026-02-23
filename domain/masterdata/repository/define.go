package repository

import (
	"boilerplate-golang/domain/masterdata"
	"database/sql"
)

const POSTGRE = "postgre"

const SELECT_DISTINCT = "select-distinct"
const SELECT = "select"
const INSERT = "insert"
const UPDATE = "update"
const DELETE = "delete"

// exaple for database configuration
const POSTGRE_TABLE = "masterdata"

const POSTGRE_JENIS_TABLE = "masterdata.jenis"
const POSTGRE_STATUS_TABLE = "masterdata.status"
const POSTGRE_LOKASI_TABLE = "masterdata.lokasi"
const POSTGRE_VENDOR_TABLE = "masterdata.vendor"
const POSTGRE_DEPARTEMEN_TABLE = "masterdata.departemen"

type postgresMasterdataRepository struct {
	sqlDB *sql.DB
}

func NewPostgresRepository(db *sql.DB) masterdata.PostgreRepository {
	return &postgresMasterdataRepository{db}
}
