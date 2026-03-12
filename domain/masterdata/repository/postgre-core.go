package repository

import (
	"svc-boilerplate-fiber-golang/utils/database"
	"svc-boilerplate-fiber-golang/valueobject"
)

func (db *postgresMasterdataRepository) GetAll(param map[string]interface{}) (response []valueobject.Masterdata, err error) {
	var result valueobject.Masterdata

	builder := database.New(POSTGRE, POSTGRE_TABLE, SELECT)

	builder.OnSelect = database.OnSelect{
		Column: []string{"id", "uuid"},
		Where:  param,
	}

	err = builder.QueryBuilder()

	if err != nil {
		return
	}

	query, err := db.sqlDB.Query(builder.Result.Query, builder.Result.Value...)

	if err != nil {
		return
	}

	defer query.Close()

	for query.Next() {
		err = query.Scan(
			&result.ID,
			&result.UUID,
		)

		if err != nil {
			return
		}

		response = append(response, result)
	}

	return
}

func (db *postgresMasterdataRepository) GetOne(param map[string]interface{}) (response valueobject.Masterdata, err error) {
	builder := database.New(POSTGRE, POSTGRE_TABLE, SELECT)

	builder.OnSelect = database.OnSelect{
		Column: []string{"id", "uuid"},
		Where:  param,
	}

	err = builder.QueryBuilder()

	if err != nil {
		return
	}

	query := db.sqlDB.QueryRow(builder.Result.Query, builder.Result.Value...)

	err = query.Scan(
		&response.ID,
		&response.UUID,
	)

	return
}

func (db *postgresMasterdataRepository) Count(param map[string]interface{}) (response valueobject.Masterdata, err error) {
	builder := database.New(POSTGRE, POSTGRE_TABLE, SELECT)

	delete(param, "LIMIT")

	builder.OnSelect = database.OnSelect{
		Column: []string{"count(*)"},
		Where:  param,
	}

	err = builder.QueryBuilder()

	if err != nil {
		return
	}

	query := db.sqlDB.QueryRow(builder.Result.Query, builder.Result.Value...)

	err = query.Scan(
		&response.Count,
	)

	return
}

func (db *postgresMasterdataRepository) Store(column []string, data []interface{}) (builder database.QueryConfig, err error) {
	builder = database.New(POSTGRE, POSTGRE_TABLE, INSERT)

	builder.OnInsert = database.OnInsert{
		Column: column,
		Data:   data,
	}

	err = builder.QueryBuilder()
	return
}

func (db *postgresMasterdataRepository) Update(param map[string]interface{}, data map[string]interface{}) (builder database.QueryConfig, err error) {
	builder = database.New(POSTGRE, POSTGRE_TABLE, UPDATE)

	builder.OnUpdate = database.OnUpdate{
		Where: param,
		Data:  data,
	}

	err = builder.QueryBuilder()
	return
}

func (db *postgresMasterdataRepository) Delete(param map[string]interface{}) (builder database.QueryConfig, err error) {
	builder = database.New(POSTGRE, POSTGRE_TABLE, DELETE)

	builder.OnDelete = database.OnDelete{
		Where: param,
	}

	err = builder.QueryBuilder()
	return
}
