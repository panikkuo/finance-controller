package utils

import (
	"fmt"
)

func GetFieldFromJsonAsString(json map[string]interface{}, field string, required bool) (string, error) {
	value, ok := json[field]
	if !ok {
		if required {
			return "", fmt.Errorf("missing in json field %s", field)
		}
	}
	strValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("can't parse to string field %s", field)
	}

	return strValue, nil
}
