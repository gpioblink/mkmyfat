package tools

import (
	"encoding/json"
)

func PrettyPrintStruct(s any) string {
	pretty, _ := json.MarshalIndent(s, "", "  ")
	return string(pretty)
}
