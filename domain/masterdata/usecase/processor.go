package usecase

import (
	"svc-boilerplate-fiber-golang/utils/database"
	"svc-boilerplate-fiber-golang/valueobject"
)

func (masterdata masterdataUsecase) ProcessStore(payload valueobject.MasterdataPayloadInsert) (queryConfig []database.QueryConfig, err error) {
	for _, x := range payload.Data {
		/**
		add data you wanted to insert on this interface{}...
		*/
		data := []interface{}{
			[]interface{}{
				x.ID,
				x.UUID,
				x.UserInput,
			},
		}

		/**
		column on data and this line should have same order
		*/
		column := []string{
			"id", "uuid", "user_input",
		}

		queryInsert, err := masterdata.postgreRepository.Store(column, data)

		if err != nil {
			return queryConfig, err
		}

		queryConfig = append(queryConfig, queryInsert)
	}
	return
}

func (masterdata masterdataUsecase) ProcessUpdate(payload valueobject.MasterdataPayloadUpdate) (queryConfig []database.QueryConfig, err error) {
	for _, x := range payload.Data {
		var param = map[string]interface{}{
			"AND": map[string]interface{}{
				"uuid": x.Param.UUID, // add the parameter to update the row
			},
		}
		var data = map[string]interface{}{
			"user_update": x.Body.UserUpdate, // add the data to update the row
		}

		queryUpdate, err := masterdata.postgreRepository.Update(param, data)

		if err != nil {
			return queryConfig, err
		}

		queryConfig = append(queryConfig, queryUpdate)
	}
	return
}

func (masterdata masterdataUsecase) ProcessDelete(payload valueobject.MasterdataPayloadDelete) (queryConfig []database.QueryConfig, err error) {
	for _, x := range payload.Param {
		var param = map[string]interface{}{
			"AND": map[string]interface{}{
				"uuid": x.UUID, // add the parameter to delete the row
			},
		}

		queryDelete, err := masterdata.postgreRepository.Delete(param)

		if err != nil {
			return queryConfig, err
		}

		queryConfig = append(queryConfig, queryDelete)
	}
	return
}
