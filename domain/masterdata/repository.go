package masterdata

import "svc-boilerplate-fiber-golang/valueobject"

type MasterdataFilter struct {
	UUID    string
	Limit   int
	Offset  int
	OrderBy string
	Desc    bool
}

type PostgreRepository interface {
	Health() error
	GenerateID() (uint64, error)
	GenerateUUID() (string, error)

	GetAll(filter MasterdataFilter) ([]valueobject.Masterdata, error)
	GetOne(uuid string) (valueobject.Masterdata, error)
	Count(filter MasterdataFilter) (int64, error)

	Store(rows []valueobject.Masterdata) error
	Update(rows []valueobject.MasterdataDataUpdate) error
	Delete(uuids []string) error
}
