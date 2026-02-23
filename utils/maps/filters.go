package maps

func EnsureAndClause(param map[string]interface{}) map[string]interface{} {
	if param == nil {
		return nil
	}

	if _, ok := param["AND"]; ok {
		return param
	}

	if len(param) == 0 {
		return nil
	}

	return map[string]interface{}{
		"AND": param,
	}
}
