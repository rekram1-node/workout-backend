package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

func writeResponse(w http.ResponseWriter, statusCode int, response any) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

// At some point might want to add injection prevention
// Or just sanitize the requests...
func validateRequest(s interface{}) error {
	validate := validator.New()
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	return fmt.Errorf("request is missing the following keys: %v", parseKeyStrings(err.Error()))
}

func parseKeyStrings(input string) string {
	var result string
	for _, line := range strings.Split(input, "\n") {
		line = strings.ToLower(line)
		if strings.Contains(line, "'") {
			parts := strings.Split(line, "'")
			split := strings.Split(parts[1], ".")
			if len(split) == 1 {
				result += split[0] + ", "
			} else {
				result += split[1] + ", "
			}
		}
	}
	trimmed := strings.TrimSpace(result)
	trimmed = strings.TrimRight(trimmed, ",")

	return trimmed
}
