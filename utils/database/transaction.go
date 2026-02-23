package database

import (
	"database/sql"
	"fmt"
)

func New(tech string, table string, command string) QueryConfig {
	return QueryConfig{
		TableInfo: TableInfo{
			TechStack: tech,
			Table:     table,
			Action:    command,
		},
	}
}

func ExecTransaction(db *sql.DB, query ...QueryConfig) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()
	listStatement := make(map[string]*sql.Stmt)
	for _, builder := range query {
		if _, found := listStatement[builder.Result.Query]; !found {
			listStatement[builder.Result.Query], err = tx.Prepare(builder.Result.Query)
			if err != nil {
				return fmt.Errorf("query: %s | value: %s | message: %s", builder.Result.Query, builder.Result.Value, err.Error())
			}
			defer listStatement[builder.Result.Query].Close()
		}
		_, err = listStatement[builder.Result.Query].Exec(builder.Result.Value...)
		if err != nil {
			return fmt.Errorf("query: %s | value: %s | message: %s", builder.Result.Query, builder.Result.Value, err.Error())
		}
	}
	tx.Commit()
	return
}
