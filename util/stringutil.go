package util

import (
	"strconv"
)

func ToString(param interface{}) string {
	if param == nil {
		return ""
	}

	switch param.(type) {
	case int:
		return strconv.Itoa(param.(int))
	case string:
		return param.(string)
	default:
		return ""
	}

	return ""
}
