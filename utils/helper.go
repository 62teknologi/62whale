package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/iancoleman/strcase"
)

type response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func ResponseData(status string, message string, data any) *response {
	return &response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func JsonFileParser(fileDir string) (map[string]any, error) {
	jsonFile, err := os.Open(fileDir)
	if err != nil {
		return nil, fmt.Errorf("error while parsing input: %w", err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var input map[string]any
	json.Unmarshal(byteValue, &input)
	return input, nil
}

func MapValuesShifter(dest map[string]any, origin map[string]any) map[string]any {
	val := reflect.ValueOf(dest)
	for _, inp := range val.MapKeys() {
		key := inp.String()
		if origin[key] == nil {
			if origin[strcase.ToCamel(key)] != nil {
				dest[key] = origin[strcase.ToCamel(key)]
				continue
			}
			continue
		}
		dest[key] = origin[key]
	}
	return dest
}

func MapNullValuesRemover(m map[string]any) {
	val := reflect.ValueOf(m)
	for _, e := range val.MapKeys() {
		v := val.MapIndex(e)
		if v.IsNil() || len(fmt.Sprintf("%v", v)) == 0 {
			delete(m, e.String())
			continue
		}
		switch t := v.Interface().(type) {
		// If key is a JSON object (Go Map), use recursion to go deeper
		case map[string]any:
			MapNullValuesRemover(t)
		}
	}
}

func LogJson(data interface{}) {
	bytes, err := json.Marshal(data)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(bytes))
}
