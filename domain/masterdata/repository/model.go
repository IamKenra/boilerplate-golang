package repository

import "svc-boilerplate-fiber-golang/valueobject"

type masterdataRow struct {
	ID         uint64 `gorm:"column:id;primaryKey"`
	UUID       string `gorm:"column:uuid"`
	UserInput  string `gorm:"column:user_input"`
	UserUpdate string `gorm:"column:user_update"`
}

func (masterdataRow) TableName() string {
	return POSTGRE_TABLE
}

func toMasterdata(row masterdataRow) valueobject.Masterdata {
	var response valueobject.Masterdata
	response.ID = row.ID
	response.UUID = row.UUID
	response.UserInput = row.UserInput
	response.UserUpdate = row.UserUpdate
	return response
}
