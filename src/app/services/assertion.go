package services

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/antonmedv/expr"
)

func assertQueryResult(queryResult []BigQueryRow, expect Expect) (bool, error) {
	if expect.Expression != nil && expect.RowCount != nil {
		return false, errors.New("only one of row_count and expression can be under expect")
	}

	if expect.Expression != nil {
		if len(queryResult) == 0 {
			return false, fmt.Errorf("the size of query result was 0 while expecting: %s in expression. If you expect the size of query result to be 0, then you can use `row_count:` instead of expression", *expect.Expression)
		}
		out, err := expr.Eval(*expect.Expression, queryResult[0])
		if err != nil {
			return false, fmt.Errorf("expression %s could not be parsed with the query result %v: %v", *expect.Expression, queryResult[0], err)
		}
		if reflect.TypeOf(out).Kind() == reflect.Bool {
			return out.(bool), nil
		}
	}

	if expect.RowCount != nil {
		return *expect.RowCount == len(queryResult), nil
	}

	return false, errors.New("either of row_count and expression must be under expect")
}
