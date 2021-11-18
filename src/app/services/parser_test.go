package services

import "testing"

func TestParseSqlParams_Sample(t *testing.T) {
	input := "SELECT count(*) FROM `{{.table}}` WHERE state = \"{{.st}}\""
	params := ParamMap{"table": "bigquery-public-data.usa_names.usa_1910_2013", "st": "TX"}
	result, err := parseSqlParams(input, params)
	if err != nil {
		t.Error(err)
	}
	expected := "SELECT count(*) FROM `bigquery-public-data.usa_names.usa_1910_2013` WHERE state = \"TX\""
	if result != expected {
		t.Errorf("SQL Parse failed.\nExpected: %s\nActual: %s", expected, result)
	}
}

func TestParseSqlParams_Env(t *testing.T) {
	input := "SELECT count(*) FROM `streaming-datatransfer-{{.env}}.streaming_datatransfer.streaming_changetracktransfer_T*`"
	params := ParamMap{"env": "stg"}
	result, err := parseSqlParams(input, params)
	if err != nil {
		t.Error(err)
	}
	expected := "SELECT count(*) FROM `streaming-datatransfer-stg.streaming_datatransfer.streaming_changetracktransfer_T*`"
	if result != expected {
		t.Errorf("SQL Parse failed.\nExpected: %s\nActual: %s", expected, result)
	}
}

func TestParseSqlParams_Float(t *testing.T) {
	input := "SELECT * FROM event_log_count_comparison_from_yesterday WHERE {{.min_ratio }} > comp_ratio OR comp_ratio > {{.max_ratio }}"
	params := ParamMap{"min_ratio": 0.95, "max_ratio": 1.05}
	result, err := parseSqlParams(input, params)
	if err != nil {
		t.Error(err)
	}
	expected := "SELECT * FROM event_log_count_comparison_from_yesterday WHERE 0.95 > comp_ratio OR comp_ratio > 1.05"
	if result != expected {
		t.Errorf("SQL Parse failed.\nExpected: %s\nActual: %s", expected, result)
	}
}

func TestParseSimple(t *testing.T) {
	input := "Hello, {{.name}}"
	params := ParamMap{"name": "Sunagimo"}
	result, err := parseSqlParams(input, params)
	if err != nil {
		t.Error(err)
	}
	expected := "Hello, Sunagimo"
	if result != expected {
		t.Errorf("SQL Parse failed.\nExpected: %s\nActual: %s", expected, result)
	}
}

func TestParseDescription_QueryResultOnly(t *testing.T) {
	description := `以下のエラーが発生しています。
	{{ range .query_result -}}
	{{ .error }}
	{{ end -}}
	`
	rows := []BigQueryRow{
		{"error": "cat"},
		{"error": "dog"},
	}

	expected := `以下のエラーが発生しています。
	cat
	dog
	`

	result, err := parseDescription(description, rows, ParamMap{}, ParamMap{})
	if err != nil {
		t.Error(err)
		return
	}
	if result != expected {
		t.Errorf("assertion failed.\nresult: %s\nexpected: %s\n", result, expected)
	}
}

func TestParseDescription_ParamsOnly(t *testing.T) {
	description := `以下のテーブルでエラーが発生しています。
	project-{{.params.env}}.source.error_log_{{.params.platform}}
	`
	params := ParamMap{"env": "stg", "platform": "android"}

	expected := `以下のテーブルでエラーが発生しています。
	project-stg.source.error_log_android
	`

	result, err := parseDescription(description, []BigQueryRow{}, params, ParamMap{})
	if err != nil {
		t.Error(err)
		return
	}
	if result != expected {
		t.Errorf("assertion failed.\nresult: %s\nexpected: %s\n", result, expected)
	}
}

func TestParseDescription_SpecialCharacters(t *testing.T) {
	description := "```\"hello\" {{ .params.hi }} ```"
	params := ParamMap{"hi": "'hi'"}

	expected := "```\"hello\" 'hi' ```"

	result, err := parseDescription(description, []BigQueryRow{}, params, ParamMap{})
	if err != nil {
		t.Error(err)
		return
	}
	if result != expected {
		t.Errorf("assertion failed.\nresult: %s\nexpected: %s\n", result, expected)
	}
}
