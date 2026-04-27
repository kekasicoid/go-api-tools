// pkg/jsonutil/formatter.go
package jsonutil

import (
	"bytes"
	"encoding/json"
	"errors"
)

type JSONFormatter struct{}

func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

func (f *JSONFormatter) Format(input string) (string, error) {
	if !json.Valid([]byte(input)) {
		return "", errors.New("invalid JSON")
	}

	var out bytes.Buffer
	err := json.Indent(&out, []byte(input), "", "  ")
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
