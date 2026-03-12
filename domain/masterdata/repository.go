package masterdata

import (
	"svc-boilerplate-fiber-golang/utils/database"
	"svc-boilerplate-fiber-golang/valueobject"
)

type PostgreRepository interface {
	Health() error
	Exec(...database.QueryConfig) error
	GenerateID() (uint64, error)
	GenerateUUID() (string, error)

	GetAll(param map[string]interface{}) ([]valueobject.Masterdata, error)
	GetOne(param map[string]interface{}) (valueobject.Masterdata, error)
	Count(param map[string]interface{}) (valueobject.Masterdata, error)

	Store(column []string, data []interface{}) (database.QueryConfig, error)
	Update(param map[string]interface{}, data map[string]interface{}) (database.QueryConfig, error)
	Delete(param map[string]interface{}) (database.QueryConfig, error)
}
