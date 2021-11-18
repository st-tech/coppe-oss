package services

import (
	"gopkg.in/yaml.v2"
)

func UnmarshalToValidations(data []byte) ([]Validation, error) {
	var validations []Validation

	if err := yaml.Unmarshal(data, &validations); err != nil {
		return nil, err
	}

	return validations, nil
}

func UnmarshalToValidation(data []byte) (Validation, error) {
	var validation Validation
	if err := yaml.Unmarshal(data, &validation); err != nil {
		return Validation{}, err
	}
	return validation, nil
}

func MarshalValidation(v Validation) ([]byte, error) {
	return yaml.Marshal(v)
}
