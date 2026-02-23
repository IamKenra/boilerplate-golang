package database

import (
	"errors"
	"log"
	"os"
	"reflect"
	"strconv"
)

type OnSelect struct {
	Column []string
	Where  map[string]interface{}
}

type OnInsert struct {
	Column []string
	Data   []interface{}
}

type OnUpdate struct {
	Where map[string]interface{}
	Data  map[string]interface{}
}

type OnDelete struct {
	Where map[string]interface{}
}

type TableInfo struct {
	TechStack string
	Table     string
	Join      string
	Action    string
}

type Result struct {
	Query string
	Value []interface{}
}

type QueryConfig struct {
	TableInfo
	OnSelect
	OnInsert
	OnUpdate
	OnDelete
	Result
	counter int
}

func (cfg *QueryConfig) QueryBuilder() (err error) {
	cfg.counter = 1
	if cfg.Action == "select" {
		cfg.selectBuilder()
		cfg.joinBuilder()
		cfg.whereBuilder(cfg.OnSelect.Where)
		cfg.groupByBuilder(cfg.OnSelect.Where)
		cfg.orderByBuilder(cfg.OnSelect.Where)
		cfg.limitBuilder(cfg.OnSelect.Where)
	} else if cfg.Action == "select-distinct" {
		cfg.selectDistinctBuilder()
		cfg.joinBuilder()
		cfg.whereBuilder(cfg.OnSelect.Where)
		cfg.groupByBuilder(cfg.OnSelect.Where)
		cfg.orderByBuilder(cfg.OnSelect.Where)
		cfg.limitBuilder(cfg.OnSelect.Where)
	} else if cfg.Action == "insert" {
		if cfg.TechStack == "oracle" && len(cfg.OnInsert.Data) > 1 {
			cfg.insertOracleBatchBuilder()
		} else {
			cfg.insertBuilder()
		}
	} else if cfg.Action == "update" {
		cfg.updateBuilder()
		if found := cfg.whereBuilder(cfg.OnUpdate.Where); !found {
			return errors.New("unsafe update with no where")
		}
	} else if cfg.Action == "delete" {
		cfg.deleteBuilder()
		if found := cfg.whereBuilder(cfg.OnDelete.Where); !found {
			return errors.New("unsafe delete with no where")
		}
	}

	// For debuging logical error
	if os.Getenv("MYSQL_DEBUG") == "true" {
		log.Println(cfg.Result.Query, cfg.Result.Value)
	}
	return
}

func (cfg *QueryConfig) joinBuilder() {
	cfg.Result.Query += ` ` + cfg.TableInfo.Join
}

func (cfg *QueryConfig) selectBuilder() {
	cfg.Result.Query += `SELECT `

	for _, x := range cfg.OnSelect.Column {
		cfg.Result.Query += x + ", "
	}

	cfg.Result.Query = cfg.Result.Query[0 : len(cfg.Result.Query)-2]

	cfg.Result.Query += ` FROM ` + cfg.Table
}

func (cfg *QueryConfig) selectDistinctBuilder() {
	cfg.Result.Query += `SELECT DISTINCT `

	for _, x := range cfg.OnSelect.Column {
		cfg.Result.Query += x + ", "
	}

	cfg.Result.Query = cfg.Result.Query[0 : len(cfg.Result.Query)-2]

	cfg.Result.Query += ` FROM ` + cfg.Table
}

func (cfg *QueryConfig) insertBuilder() {
	cfg.Result.Query += `INSERT INTO ` + cfg.Table + ` (`

	for _, x := range cfg.OnInsert.Column {
		cfg.Result.Query += x + `, `
	}

	cfg.Result.Query = cfg.Result.Query[0 : len(cfg.Result.Query)-2]
	cfg.Result.Query += `) VALUES `

	for _, x := range cfg.OnInsert.Data {
		count := len(x.([]interface{}))

		if count < 0 {
			count = 0
		}

		cfg.Result.Query += `(`
		for i := 0; i < count; i++ {
			cfg.Result.Query += cfg.getQuestionMark() + `, `
		}
		cfg.Result.Query = cfg.Result.Query[0 : len(cfg.Result.Query)-2]
		cfg.Result.Query += `),`

		cfg.Result.Value = append(cfg.Result.Value, x.([]interface{})...)
	}

	cfg.Result.Query = cfg.Result.Query[0 : len(cfg.Result.Query)-1]
}

func (cfg *QueryConfig) insertOracleBatchBuilder() {
	cfg.Result.Query += `INSERT ALL`

	for _, x := range cfg.OnInsert.Data {
		count := len(x.([]interface{}))

		if count < 0 {
			count = 0
		}

		cfg.Result.Query += ` INTO ` + cfg.Table + `(`

		for _, x := range cfg.OnInsert.Column {
			cfg.Result.Query += x + `, `
		}

		cfg.Result.Query = cfg.Result.Query[0 : len(cfg.Result.Query)-2]
		cfg.Result.Query += `) VALUES `

		cfg.Result.Query += `(`
		for i := 0; i < count; i++ {
			cfg.Result.Query += cfg.getQuestionMark() + `, `
		}
		cfg.Result.Query = cfg.Result.Query[0 : len(cfg.Result.Query)-2]
		cfg.Result.Query += `) `

		cfg.Result.Value = append(cfg.Result.Value, x.([]interface{})...)
	}

	cfg.Result.Query += `SELECT * FROM dual`
}

func (cfg *QueryConfig) updateBuilder() {
	cfg.Result.Query += `UPDATE ` + cfg.Table + ` SET `

	for i, x := range cfg.OnUpdate.Data {
		cfg.Result.Query += i + ` = ` + cfg.getQuestionMark() + `, `
		cfg.Result.Value = append(cfg.Result.Value, x)
	}

	cfg.Result.Query = cfg.Result.Query[0 : len(cfg.Result.Query)-2]
}

func (cfg *QueryConfig) deleteBuilder() {
	cfg.Result.Query += `DELETE FROM ` + cfg.Table
}

func (cfg *QueryConfig) whereBuilder(param map[string]interface{}) (found bool) {
	cfg.Result.Query += ` WHERE `

	for i, x := range param {
		if i == "AND" {
			for g, v := range x.(map[string]interface{}) {
				if g == "IN" {
					for o, f := range v.(map[string][]string) {
						r := len(f)
						if r < 1 {
							continue
						}
						nullFound := false
						localFound := false
						localQuery := ``
						for _, w := range f {
							if w == "" {
								continue
							}
							if w == "null" {
								nullFound = true
							}
							localFound = true
							localQuery += cfg.getQuestionMark() + `, `
							cfg.Result.Value = append(cfg.Result.Value, w)
						}
						if localFound {
							if nullFound {
								cfg.Result.Query += `(`
							}
							cfg.Result.Query += o + ` IN (`
							cfg.Result.Query += localQuery
							cfg.Result.Query = cfg.Result.Query[0 : len(cfg.Result.Query)-2]
							cfg.Result.Query += `)`
							if nullFound {
								cfg.Result.Query += ` OR ` + o + ` IS NULL)`
							}
							cfg.Result.Query += ` AND `
							found = true
						}
					}
				} else if g == "OR" {
					for _, f := range v.([]map[string]interface{}) {
						var localValues []interface{}
						localQuery := ``
						counter := 0
						for p, t := range f {
							if t == "" {
								continue
							}
							if p == "IN" {
								for o, f := range t.(map[string][]string) {
									r := len(f)
									if r < 1 {
										continue
									}
									nullFound := false
									localFound := false
									localInQuery := ``
									for _, w := range f {
										if w == "" {
											continue
										}
										if w == "null" {
											nullFound = true
										}
										localFound = true
										localInQuery += cfg.getQuestionMark() + `, `
										localValues = append(localValues, w)
									}
									if localFound {
										if nullFound {
											localQuery += `(`
										}
										localQuery += o + ` IN (`
										localQuery += localInQuery
										localQuery = localQuery[0 : len(localQuery)-2]
										localQuery += `)`
										if nullFound {
											localQuery += ` OR ` + o + ` IS NULL)`
										}
										localQuery += ` OR `
										found = true
									}
									counter++
								}
							} else if t == nil {
								localQuery += p + ` IS NULL OR `
							} else if reflect.TypeOf(t).String() == "[]interface {}" {
								localQuery += p + ` ` + t.([]interface{})[0].(string) + ` ` + cfg.getQuestionMark() + ` OR `
								localValues = append(localValues, t.([]interface{})[1])
							} else {
								localQuery += p + ` = ` + cfg.getQuestionMark() + ` OR `
								localValues = append(localValues, t)
							}
							counter++
						}
						if counter >= 2 {
							cfg.Result.Query += `(`
							cfg.Result.Query += localQuery
							cfg.Result.Query = cfg.Result.Query[0 : len(cfg.Result.Query)-4]
							cfg.Result.Query += `) AND `
							cfg.Result.Value = append(cfg.Result.Value, localValues...)
							found = true
						}
					}
				} else if g == "NOT" {
					for o, f := range v.(map[string]interface{}) {
						if f == "" {
							continue
						}
						if f == nil {
							cfg.Result.Query += o + ` IS NOT NULL AND `
							found = true
						} else if o == "IN" {
							for q, l := range f.(map[string][]string) {
								r := len(l)
								if r < 1 {
									continue
								}
								localFound := false
								localQuery := ``
								for _, w := range l {
									if w == "" {
										continue
									}
									localFound = true
									localQuery += cfg.getQuestionMark() + `, `
									cfg.Result.Value = append(cfg.Result.Value, w)
								}
								if localFound {
									cfg.Result.Query += q + ` NOT IN (`
									cfg.Result.Query += localQuery
									cfg.Result.Query = cfg.Result.Query[0 : len(cfg.Result.Query)-2]
									cfg.Result.Query += `) AND `
									found = true
								}
							}
						} else {
							cfg.Result.Query += `NOT ` + o + ` = ` + cfg.getQuestionMark() + ` AND `
							cfg.Result.Value = append(cfg.Result.Value, f)
							found = true
						}
					}
				} else if g == "LIKE" {
					for o, f := range v.(map[string]interface{}) {
						if f == "" {
							continue
						} else {
							cfg.Result.Query += o + ` LIKE ` + cfg.getQuestionMark() + ` AND `
							cfg.Result.Value = append(cfg.Result.Value, f)
							found = true
						}
					}
				} else if g == "BETWEEN" {
					for o, f := range v.(map[string][]interface{}) {
						if len(f) != 2 {
							continue
						} else if f[0] == nil && f[1] == nil {
							continue
						} else if f[0] == nil && f[1] != nil {
							cfg.Result.Query += o + ` < ` + cfg.getQuestionMark() + ` AND `
							cfg.Result.Value = append(cfg.Result.Value, f[1])
							found = true
						} else if f[0] != nil && f[1] == nil {
							cfg.Result.Query += o + ` > ` + cfg.getQuestionMark() + ` AND `
							cfg.Result.Value = append(cfg.Result.Value, f[0])
							found = true
						} else {
							cfg.Result.Query += o + ` BETWEEN ` + cfg.getQuestionMark() + ` AND ` + cfg.getQuestionMark() + ` AND `
							cfg.Result.Value = append(cfg.Result.Value, f...)
							found = true
						}
					}
				} else {
					if g == "" {
						continue
					}
					if v == "" {
						continue
					}
					if v == nil {
						cfg.Result.Query += g + ` IS NULL AND `
						found = true
					} else if reflect.TypeOf(v).String() == "[]interface {}" {
						cfg.Result.Query += g + ` AND `
						cfg.Result.Value = append(cfg.Result.Value, v.([]interface{})...)
						found = true
					} else {
						cfg.Result.Query += g + ` = ` + cfg.getQuestionMark() + ` AND `
						cfg.Result.Value = append(cfg.Result.Value, v)
						found = true
					}
				}
			}
			if found {
				cfg.Result.Query = cfg.Result.Query[0 : len(cfg.Result.Query)-5]
			}
		}
	}

	if !found {
		cfg.Result.Query = cfg.Result.Query[0 : len(cfg.Result.Query)-7]
	}

	return
}

func (cfg *QueryConfig) groupByBuilder(param map[string]interface{}) (err error) {
	var found bool

	for a, c := range param {
		if a == "GROUP" {
			for i, x := range c.(map[string]interface{}) {
				if i == "GROUP_BY" {
					r := len(x.([]string))
					if r < 1 {
						continue
					}
					if r == 1 {
						if x.([]string)[0] == "" {
							continue
						}
					}
					cfg.Result.Query += ` GROUP BY `
					for _, w := range x.([]string) {
						if w == "" {
							continue
						}
						cfg.Result.Query += w + `, `
					}
					cfg.Result.Query = cfg.Result.Query[0 : len(cfg.Result.Query)-2]
					cfg.Result.Query += ` `
					found = true
				}
			}
		}
	}

	if found {
		cfg.Result.Query = cfg.Result.Query[0 : len(cfg.Result.Query)-1]
	}

	return
}

func (cfg *QueryConfig) orderByBuilder(param map[string]interface{}) (err error) {
	var found bool

	for a, c := range param {
		if a == "ORDER" {
			for i, x := range c.(map[string]interface{}) {
				if i == "ORDER_BY" {
					r := len(x.([]string))
					if r < 1 {
						continue
					}
					cfg.Result.Query += ` ORDER BY `
					for _, w := range x.([]string) {
						if w == "" {
							continue
						}
						cfg.Result.Query += w + `, `
					}
					cfg.Result.Query = cfg.Result.Query[0 : len(cfg.Result.Query)-2]
					cfg.Result.Query += ` `
					found = true
				}
			}
		}
	}

	if found {
		cfg.Result.Query = cfg.Result.Query[0 : len(cfg.Result.Query)-1]
	}

	return
}

func (cfg *QueryConfig) limitBuilder(param map[string]interface{}) (err error) {
	for a, c := range param {
		if a == "LIMIT" {
			for _, v := range c.([]interface{}) {
				if v == "" {
					return
				}
			}
			values := c.([]interface{})
			if len(values) != 2 {
				return
			}

			// values[0] = offset, values[1] = limit
			if cfg.TechStack == "postgre" {
				// PostgreSQL syntax: LIMIT $n OFFSET $m
				// Gunakan urutan placeholder yang konsisten dengan Value.
				cfg.Result.Query += ` LIMIT ` + cfg.getQuestionMark()
				cfg.Result.Value = append(cfg.Result.Value, values[1])
				cfg.Result.Query += ` OFFSET ` + cfg.getQuestionMark()
				cfg.Result.Value = append(cfg.Result.Value, values[0])
			} else {
				// MySQL style: LIMIT offset, count
				cfg.Result.Query += ` LIMIT ` + cfg.getQuestionMark() + `, ` + cfg.getQuestionMark()
				cfg.Result.Value = append(cfg.Result.Value, values...)
			}
		}
	}
	return
}

func (cfg *QueryConfig) getQuestionMark() (questionMark string) {
	switch cfg.TechStack {
	case "oracle":
		questionMark = ":x" + strconv.Itoa(cfg.counter)
		cfg.counter++
	case "mysql":
		questionMark = "?"
	case "postgre":
		questionMark = "$" + strconv.Itoa(cfg.counter)
		cfg.counter++
	}
	return
}
