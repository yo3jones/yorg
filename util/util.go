package util

import "encoding/json"

func Pretty(a any) string {
	data, _ := json.MarshalIndent(a, "", "  ")
	return string(data)
}
