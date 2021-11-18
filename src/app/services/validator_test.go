package services

import (
	"testing"

	"github.com/adhocore/gronx"
)

func TestSchedule(t *testing.T) {
	fileDataMap, _ := ReadFiles("../../../yaml")
	var validations []Validation
	for _, data := range fileDataMap {
		validationsInFile, err := UnmarshalToValidations(data)
		if err != nil {
			t.Error("ERROR: yaml unmarshal failed. ", err)
		}
		validations = append(validations, validationsInFile...)
	}

	gron := gronx.New()
	for _, v := range validations {
		_, err := gron.IsDue(v.Schedule)
		if err != nil {
			t.Errorf("ERROR: schedule: %s is not valid format", v.Schedule)
		}
	}
}

func TestGenerateAlertMsg(t *testing.T) {
	sql := "select count(*) as cnt from foo"
	v := Validation{
		Identifier: ValidationIdentifier{
			FoundAt: ValidationLocation{
				FileName: "streaming-datatransfer.yaml",
				Index:    0,
			},
			MatrixValue: ParamMap{
				"env": "prd",
			},
		},
		Schedule:    "* * * * *",
		Sql:         sql,
		Description: "Error message here",
	}
	queryResult := []BigQueryRow{
		{"cnt": 0},
	}

	expectedMsg := "*File:*\nstreaming-datatransfer.yaml\n*Index:*\n0\n*Matrix Value:*\n- env: prd\n*Description:*\nError message here\n"

	msg, err := generateAlertMessage(v, sql, queryResult)
	if err != nil {
		t.Error(err)
	}
	if msg != expectedMsg {
		t.Errorf("alert message is not correctly generated.\nexpect: %s\n\nactual: %s\n", expectedMsg, msg)
	}
}
