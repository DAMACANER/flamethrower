package views

import (
	"strings"
)

func notEmpty(v interface{}) bool {
	if v == nil {
		return false
	}

	switch v := v.(type) {
	case string:
		trimmed := strings.TrimSpace(v)
		return trimmed != "" && trimmed != "None"
	case *string:
		if v == nil {
			return false
		}
		trimmed := strings.TrimSpace(*v)
		return trimmed != "" && trimmed != "None"
	case *int, *int32, *int64:
		return v != nil
	case *float32, *float64:
		return v != nil
	default:
		return true
	}
}
