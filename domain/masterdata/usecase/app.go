package usecase

import (
	"svc-boilerplate-fiber-golang/domain/masterdata"
)

/*
*
the struct of usecases
*/
type masterdataUsecase struct {
	postgreRepository masterdata.PostgreRepository
}

/*
*
the initiator function for usecase
*/

func NewMasterdataUsecase(postgres masterdata.PostgreRepository) masterdata.Usecase {
	return &masterdataUsecase{
		postgreRepository: postgres,
	}
}
