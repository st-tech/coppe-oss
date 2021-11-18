package services

import (
	"testing"
)

func TestAssertion_SimpleExpression_OneIntColumn(t *testing.T) {
	queryResult := []BigQueryRow{{"foo": 0}}
	expression := "foo == 0"
	expect := Expect{
		Expression: &expression,
	}

	result, err := assertQueryResult(queryResult, expect)
	if err != nil {
		t.Error(err)
	}
	if !result {
		t.Errorf("Assertion Failed while expected to be true")
	}
}

func TestAssertion_SimpleExpression_OneStringColumn(t *testing.T) {
	queryResult := []BigQueryRow{{"foo": "f"}}
	expression := "foo == \"f\""
	expect := Expect{
		Expression: &expression,
	}

	result, err := assertQueryResult(queryResult, expect)
	if err != nil {
		t.Error(err)
	}
	if !result {
		t.Errorf("Assertion Failed while expected to be true")
	}
}

func TestAssertion_SimpleExpression_MultipleColumns(t *testing.T) {
	queryResult := []BigQueryRow{{"foo": "f", "bar": "b"}}
	expression := "foo == \"f\""
	expect := Expect{
		Expression: &expression,
	}

	result, err := assertQueryResult(queryResult, expect)
	if err != nil {
		t.Error(err)
	}
	if !result {
		t.Errorf("Assertion Failed while expected to be true")
	}
}

func TestAssertion_SimpleExpression_With_MultipleColumns_OfSameValue(t *testing.T) {
	queryResult := []BigQueryRow{{"foo": "f", "bar": "f"}}
	expression := "foo == bar"
	expect := Expect{
		Expression: &expression,
	}

	result, err := assertQueryResult(queryResult, expect)
	if err != nil {
		t.Error(err)
	}
	if !result {
		t.Errorf("Assertion Failed while expected to be true")
	}
}

func TestAssertion_ComplexExpression(t *testing.T) {
	queryResult := []BigQueryRow{{"foo": "f", "bar": "b"}}
	expression := "foo == \"f\" && bar == \"b\""
	expect := Expect{
		Expression: &expression,
	}

	result, err := assertQueryResult(queryResult, expect)
	if err != nil {
		t.Error(err)
	}
	if !result {
		t.Errorf("Assertion Failed while expected to be true")
	}
}

func TestAssertion_RowCount0(t *testing.T) {
	queryResult := []BigQueryRow{}
	rowCount := 0
	expect := Expect{
		RowCount: &rowCount,
	}

	result, err := assertQueryResult(queryResult, expect)
	if err != nil {
		t.Error(err)
	}
	if !result {
		t.Errorf("Assertion Failed while expected to be true")
	}
}

func TestAssertion_RowCountNot0(t *testing.T) {
	queryResult := []BigQueryRow{{"foo": "f"}}
	rowCount := 1
	expect := Expect{
		RowCount: &rowCount,
	}

	result, err := assertQueryResult(queryResult, expect)
	if err != nil {
		t.Error(err)
	}
	if !result {
		t.Errorf("Assertion Failed while expected to be true")
	}
}

func TestAssertion_Expression_and_RowCountNone0(t *testing.T) {
	queryResult := []BigQueryRow{{"foo": "f"}}
	expression := "foo == \"f\""
	rowCount := 1

	expect := Expect{
		Expression: &expression,
		RowCount:   &rowCount,
	}

	_, err := assertQueryResult(queryResult, expect)
	if err == nil {
		t.Error("Error must occur when row_count and expression are both in expect")
	}
}

func TestAssertion_Expression_and_RowCount0(t *testing.T) {
	queryResult := []BigQueryRow{{"foo": "f"}}
	expression := "foo == \"f\""
	rowCount := 0
	expect := Expect{
		Expression: &expression,
		RowCount:   &rowCount,
	}

	_, err := assertQueryResult(queryResult, expect)
	if err == nil {
		t.Error("Error must occur when row_count and expression are both in expect")
	}
}
