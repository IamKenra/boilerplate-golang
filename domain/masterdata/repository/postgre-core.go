package repository

import (
	"errors"
	"svc-boilerplate-fiber-golang/domain/masterdata"
	"svc-boilerplate-fiber-golang/valueobject"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var masterdataOrderColumns = map[string]string{
	"id":   "id",
	"uuid": "uuid",
}

func (db *postgresMasterdataRepository) GetAll(filter masterdata.MasterdataFilter) (response []valueobject.Masterdata, err error) {
	var rows []masterdataRow
	query, err := db.applyFilter(db.db.Model(&masterdataRow{}), filter, true)
	if err != nil {
		return nil, err
	}

	if err = query.Find(&rows).Error; err != nil {
		return nil, err
	}

	response = make([]valueobject.Masterdata, 0, len(rows))
	for _, row := range rows {
		response = append(response, toMasterdata(row))
	}
	return response, nil
}

func (db *postgresMasterdataRepository) GetOne(uuid string) (response valueobject.Masterdata, err error) {
	if uuid == "" {
		return response, errors.New("uuid is required")
	}

	var row masterdataRow
	if err = db.db.Where("uuid = ?", uuid).First(&row).Error; err != nil {
		return response, err
	}
	return toMasterdata(row), nil
}

func (db *postgresMasterdataRepository) Count(filter masterdata.MasterdataFilter) (response int64, err error) {
	query, err := db.applyFilter(db.db.Model(&masterdataRow{}), filter, false)
	if err != nil {
		return 0, err
	}
	return response, query.Count(&response).Error
}

func (db *postgresMasterdataRepository) Store(rows []valueobject.Masterdata) error {
	if len(rows) == 0 {
		return nil
	}

	payload := make([]masterdataRow, 0, len(rows))
	for _, row := range rows {
		payload = append(payload, masterdataRow{
			ID:        row.ID,
			UUID:      row.UUID,
			UserInput: row.UserInput,
		})
	}

	return db.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&payload).Error
	})
}

func (db *postgresMasterdataRepository) Update(rows []valueobject.MasterdataDataUpdate) error {
	if len(rows) == 0 {
		return nil
	}

	return db.db.Transaction(func(tx *gorm.DB) error {
		for _, row := range rows {
			if row.Param.UUID == "" {
				return errors.New("uuid is required")
			}
			if err := tx.Model(&masterdataRow{}).
				Where("uuid = ?", row.Param.UUID).
				Updates(map[string]interface{}{"user_update": row.Body.UserUpdate}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *postgresMasterdataRepository) Delete(uuids []string) error {
	if len(uuids) == 0 {
		return nil
	}

	return db.db.Transaction(func(tx *gorm.DB) error {
		for _, uuid := range uuids {
			if uuid == "" {
				return errors.New("uuid is required")
			}
			if err := tx.Where("uuid = ?", uuid).Delete(&masterdataRow{}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *postgresMasterdataRepository) applyFilter(query *gorm.DB, filter masterdata.MasterdataFilter, withPagination bool) (*gorm.DB, error) {
	if filter.UUID != "" {
		query = query.Where("uuid = ?", filter.UUID)
	}

	if withPagination {
		if filter.OrderBy != "" {
			column, ok := masterdataOrderColumns[filter.OrderBy]
			if !ok {
				return nil, errors.New("invalid order by")
			}
			query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: column}, Desc: filter.Desc})
		}
		if filter.Limit > 0 {
			query = query.Limit(filter.Limit)
		}
		if filter.Offset > 0 {
			query = query.Offset(filter.Offset)
		}
	}

	return query, nil
}
