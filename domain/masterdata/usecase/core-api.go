package usecase

import (
	"boilerplate-golang/valueobject"
)

func (masterdata masterdataUsecase) GetAll(param map[string]interface{}) (response []valueobject.Masterdata, err error) {
	response, err = masterdata.postgreRepository.GetAll(param)
	return
}

func (masterdata masterdataUsecase) GetOne(param map[string]interface{}) (response valueobject.Masterdata, err error) {
	response, err = masterdata.postgreRepository.GetOne(param)
	return
}

func (masterdata masterdataUsecase) Count(param map[string]interface{}) (response valueobject.Masterdata, err error) {
	response, err = masterdata.postgreRepository.Count(param)
	return
}

func (masterdata masterdataUsecase) Store(payload valueobject.MasterdataPayloadInsert) (valueobject.MasterdataPayloadInsert, error) {
	for i := range payload.Data {
		payload.Data[i].ID, _ = masterdata.postgreRepository.GenerateID()
		payload.Data[i].UUID, _ = masterdata.postgreRepository.GenerateUUID()
		payload.Data[i].UserInput = payload.User
	}

	queryConfig, err := masterdata.ProcessStore(payload)

	if err != nil {
		return payload, err
	}

	return payload, masterdata.postgreRepository.Exec(queryConfig...)
}

func (masterdata masterdataUsecase) Update(payload valueobject.MasterdataPayloadUpdate) (err error) {
	for i := range payload.Data {
		payload.Data[i].Body.UserUpdate = payload.User
	}

	queryConfig, err := masterdata.ProcessUpdate(payload)

	if err != nil {
		return
	}

	return masterdata.postgreRepository.Exec(queryConfig...)
}

func (masterdata masterdataUsecase) Delete(payload valueobject.MasterdataPayloadDelete) (err error) {
	queryConfig, err := masterdata.ProcessDelete(payload)

	if err != nil {
		return
	}

	return masterdata.postgreRepository.Exec(queryConfig...)
}
