package maps

func Merge(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			if existingVal, ok := result[k]; ok {
				if existingMap, ok1 := existingVal.(map[string]interface{}); ok1 {
					if newMap, ok2 := v.(map[string]interface{}); ok2 {
						result[k] = Merge(existingMap, newMap)
					} else {
						result[k] = v
					}
				} else {
					result[k] = v
				}
			} else {
				result[k] = v
			}
		}
	}
	return result
}
