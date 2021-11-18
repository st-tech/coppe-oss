package services

import (
	"testing"
)

func TestYamlUnmarshall(t *testing.T) {
	fileDataMap, _ := ReadFiles("../../../yaml")
	for _, data := range fileDataMap {
		_, err := UnmarshalToValidations(data)
		if err != nil {
			t.Error("ERROR: yaml unmarshal failed. ", err)
		}
	}
}

func TestYamlMarshall(t *testing.T) {
	fileDataMap, _ := ReadFiles("../../../yaml")

	var validations []Validation
	for _, data := range fileDataMap {
		_, err := UnmarshalToValidations(data)
		if err != nil {
			t.Error("ERROR: yaml unmarshal failed. ", err)
		}
	}

	for _, v := range validations {
		_, err := MarshalValidation(v)
		if err != nil {
			t.Error("ERROR: yaml marshal failed. ", err)
		}
	}
}
