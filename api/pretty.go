package api

import (
	"encoding/json"
	"fmt"
)

func PrettyPrintJSON(data []byte) string {
	var obj interface{}

	if err := json.Unmarshal(data, &obj); err != nil {
		return fmt.Sprintf("invalid JSON: %v", err)
	}

	formatted, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Sprintf("failed to format JSON: %v", err)
	}

	return string(formatted)
}
