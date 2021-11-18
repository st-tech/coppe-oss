package services

import (
	"bytes"
	"text/template"
)

func parseSqlParams(sql string, params ParamMap) (string, error) {
	if len(params) == 0 {
		return sql, nil
	}
	t, err := template.New("Parse sql parameters").Parse(sql)
	if err != nil {
		return "", err
	}
	return parse(t, params)
}

func parseDescription(description string, rows []BigQueryRow, params ParamMap, matrixValue ParamMap) (string, error) {
	data := map[string]interface{}{}
	if len(params) != 0 {
		data["params"] = params
	}
	if len(rows) != 0 {
		data["query_result"] = rows
	}
	if len(matrixValue) != 0 {
		data["matrix"] = matrixValue
	}
	if len(data) == 0 {
		return description, nil
	}
	t, err := template.New("Parse Description with Query Result and SQL Params").Parse(description)
	if err != nil {
		return "", err
	}
	parsedDescription, err := parse(t, data)
	if err != nil {
		return "", err
	}
	return parsedDescription, nil
}

func parse(t *template.Template, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
