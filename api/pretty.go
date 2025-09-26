package api

import (
	"encoding/json"
	"fmt"
)

func PrettyPrintJSON(data []byte) error {
	var obj interface{}

	if err := json.Unmarshal(data, &obj); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	formatted, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format JSON: %w", err)
	}

	fmt.Println(string(formatted))
	return nil
}
