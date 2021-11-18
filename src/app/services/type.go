package services

import (
	"fmt"

	"cloud.google.com/go/bigquery"
)

type Validation struct {
	Identifier  ValidationIdentifier
	Schedule    string
	Sql         string
	SqlFile     string `yaml:"sql_file"`
	Params      ParamMap
	Matrix      MatrixMap
	Expect      Expect
	Description string
	AlertFormat string `yaml:"alert_format"`
	Channel     string
}

type Expect struct {
	Expression *string
	RowCount   *int `yaml:"row_count"`
}

type ValidationIdentifier struct {
	FoundAt     ValidationLocation
	MatrixValue ParamMap
}

func (v *ValidationIdentifier) ToString() string {
	return fmt.Sprintf("%s, MatrixValue: %+v", v.FoundAt.ToString(), v.MatrixValue)
}

type ValidationLocation struct {
	FileName string
	Index    int
}

func (v *ValidationLocation) ToString() string {
	return fmt.Sprintf("FileName: %s, Index: %d", v.FileName, v.Index)
}

type ParamMap = map[string]interface{}

type MatrixMap = map[string][]interface{}

type BigQueryRow = map[string]bigquery.Value

type SlackTopBlock struct {
	Blocks []interface{} `json:"blocks"`
}

type SlackParentBlockForSingleLeaf struct {
	Type string         `json:"type"`
	Text SlackLeafBlock `json:"text"`
}

type SlackLeafBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
