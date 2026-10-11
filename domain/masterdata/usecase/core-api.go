package usecase

import (
	"svc-boilerplate-fiber-golang/domain/masterdata"
	"svc-boilerplate-fiber-golang/valueobject"
)

func (masterdataUC masterdataUsecase) GetAll(filter masterdata.MasterdataFilter) (response []valueobject.Masterdata, err error) {
	return masterdataUC.postgreRepository.GetAll(filter)
}

func (masterdataUC masterdataUsecase) GetOne(uuid string) (response valueobject.Masterdata, err error) {
	return masterdataUC.postgreRepository.GetOne(uuid)
}

func (masterdataUC masterdataUsecase) Count(filter masterdata.MasterdataFilter) (response int64, err error) {
	return masterdataUC.postgreRepository.Count(filter)
}

func (masterdataUC masterdataUsecase) Store(payload valueobject.MasterdataPayloadInsert) (valueobject.MasterdataPayloadInsert, error) {
	for i := range payload.Data {
		id, err := masterdataUC.postgreRepository.GenerateID()
		if err != nil {
			return payload, err
		}
		uuid, err := masterdataUC.postgreRepository.GenerateUUID()
		if err != nil {
			return payload, err
		}
		payload.Data[i].ID = id
		payload.Data[i].UUID = uuid
		payload.Data[i].UserInput = payload.User
	}

	return payload, masterdataUC.postgreRepository.Store(payload.Data)
}

func (masterdataUC masterdataUsecase) Update(payload valueobject.MasterdataPayloadUpdate) error {
	for i := range payload.Data {
		payload.Data[i].Body.UserUpdate = payload.User
	}
	return masterdataUC.postgreRepository.Update(payload.Data)
}

func (masterdataUC masterdataUsecase) Delete(payload valueobject.MasterdataPayloadDelete) error {
	uuids := make([]string, 0, len(payload.Param))
	for _, item := range payload.Param {
		uuids = append(uuids, item.UUID)
	}
	return masterdataUC.postgreRepository.Delete(uuids)
}
