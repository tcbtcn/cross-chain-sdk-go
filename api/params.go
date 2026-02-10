package api

import (
	"fmt"
	"net/url"
	"strconv"
)

func ConcatQueryParams(params map[string]interface{}) string {
	if len(params) == 0 {
		return ""
	}

	values := url.Values{}
	for k, v := range params {
		if v == nil {
			continue
		}
		switch val := v.(type) {
		case string:
			values.Add(k, val)
		case int:
			values.Add(k, strconv.Itoa(val))
		case int64:
			values.Add(k, strconv.FormatInt(val, 10))
		case bool:
			values.Add(k, strconv.FormatBool(val))
		case []string:
			for _, s := range val {
				values.Add(k, s)
			}
		default:
			values.Add(k, fmt.Sprintf("%v", val))
		}
	}

	query := values.Encode()
	if query != "" {
		return "?" + query
	}
	return ""
}
