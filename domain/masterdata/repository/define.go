package repository

import (
	"svc-boilerplate-fiber-golang/domain/masterdata"

	"gorm.io/gorm"
)

const POSTGRE_TABLE = "masterdata"

type postgresMasterdataRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) masterdata.PostgreRepository {
	return &postgresMasterdataRepository{db: db}
}
