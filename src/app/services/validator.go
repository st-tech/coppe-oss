package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"time"
)

func Validate(ctx context.Context, v Validation) error {
	SendInfoLog(fmt.Sprintf("Validation id - %s", v.Identifier.ToString()))
	SendDebugLog(fmt.Sprintf("Validating: %+v", v))

	sql, err := parseSqlParams(v.Sql, mergeParamsAndMatrix(v.Params, v.Identifier.MatrixValue))
	if err != nil {
		return fmt.Errorf("error occurred while parseing SQL with parameters: %v", err)
	}
	SendDebugLog(fmt.Sprintf("SQL-after-parse: %s", sql))

	tctx, cancel := context.WithTimeout(ctx, 500*time.Second) // since Cloud Functions can only run for ~540 seconds
	defer cancel()

	queryResult, err := queryToBigQuery(tctx, sql, false)
	if err != nil {
		return fmt.Errorf("error occurred while query to BigQuery: %v", err)
	}

	if len(queryResult) > 0 {
		SendInfoLog(fmt.Sprintf("%+v", queryResult[0]))
	}
	assertSucceeded, err := assertQueryResult(queryResult, v.Expect)
	if err != nil {
		return fmt.Errorf("error occurred while assertion: %v", err)
	}

	if !assertSucceeded {
		msg, err := generateAlertMessage(v, sql, queryResult)
		if err != nil {
			return err
		}
		SendWarningLog(msg)
		PostSlackMessage(msg, v.Channel)
		return nil
	}

	return nil
}

func GetSql(sql string, sqlFile string) (string, error) {
	if sql != "" && sqlFile != "" {
		return "", errors.New("only one of sql: and sql_file: must exist in one validation")
	}
	if sql != "" {
		return sql, nil
	}
	if sqlFile != "" {
		bytes, err := readFile(filepath.Join("sql", sqlFile))
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
	return "", errors.New("no SQL found, sql: or sql_file: is required for validation")
}

func mergeParamsAndMatrix(paramMap ParamMap, matrixValue ParamMap) ParamMap {
	params := ParamMap{}
	for k, v := range paramMap {
		params[k] = v
	}
	for k, v := range matrixValue {
		params[k] = v
	}
	return params
}

func generateAlertMessage(v Validation, parsedSql string, queryResult []BigQueryRow) (string, error) {
	parsedDescription, err := parseDescription(v.Description, queryResult, v.Params, v.Identifier.MatrixValue)
	if err != nil {
		return "", fmt.Errorf("error occurred while parsing description with query result and params: ERROR: %v, \ndescription: %s\nquery result (first row) %v\nSQL params %v", err, v.Description, queryResult[0], v.Params)
	}

	resultJson, err := json.Marshal(queryResult[0])
	if err != nil {
		return "", fmt.Errorf("error occurred while marshaling query result to json: %v", queryResult[0])
	}

	matrixValueStr := ""
	for k, v := range v.Identifier.MatrixValue {
		matrixValueStr += fmt.Sprintf("- %s: %v\n", k, v)
	}
	validationIdMsg := fmt.Sprintf("*File:*\n%s\n*Index:*\n%d\n", v.Identifier.FoundAt.FileName, v.Identifier.FoundAt.Index)
	if matrixValueStr != "" {
		validationIdMsg = fmt.Sprintf("%s%s\n%s", validationIdMsg, "*Matrix Value:*", matrixValueStr)
	}
	switch v.AlertFormat {
	case "verbose":
		return validationIdMsg + fmt.Sprintf("*Description:*\n%s\n*SQL*\n```%s```\n*Expected*\n%v\n*Query Result (first row)*\n%s", parsedDescription, parsedSql, v.Expect, string(resultJson)), nil
	case "simple", "":
		return validationIdMsg + fmt.Sprintf("*Description:*\n%s\n", parsedDescription), nil
	}
	return "", errors.New("alert_format only accepts simple or verobse currently (if empty, simple format is applied)")
}
