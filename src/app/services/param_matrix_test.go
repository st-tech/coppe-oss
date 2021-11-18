package services

import (
	"testing"
)

func TestGenerateParamCombinations_MatrixOnly(t *testing.T) {
	matrixMap := MatrixMap{
		"env":      []interface{}{"dev", "stg", "prd"},
		"platform": []interface{}{"browser", "android", "server", "ios"},
	}

	params, err := GenerateMatrixCombinations(matrixMap)
	if err != nil {
		t.Error(err)
	}

	expectedPatternSize := 12
	if len(params) != expectedPatternSize {
		t.Errorf("Should generate %d patterns, but got %d", expectedPatternSize, len(params))
	}
}

func TestGenerateParamCombinations_PatternMatch(t *testing.T) {
	matrixMap := MatrixMap{
		"env":      []interface{}{"stg", "prd"},
		"platform": []interface{}{"browser", "android"},
	}

	params, err := GenerateMatrixCombinations(matrixMap)
	if err != nil {
		t.Error(err)
	}

	stgBrowser := 0
	stgAndroid := 0
	prdBrowser := 0
	prdAndroid := 0

	for _, param := range params {
		env := param["env"]
		platform := param["platform"]

		if env == "stg" && platform == "browser" {
			stgBrowser++
		}
		if env == "stg" && platform == "android" {
			stgAndroid++
		}
		if env == "prd" && platform == "browser" {
			prdBrowser++
		}
		if env == "prd" && platform == "android" {
			prdAndroid++
		}
	}
	if stgBrowser != 1 || stgAndroid != 1 || prdBrowser != 1 || prdAndroid != 1 {
		t.Errorf("param combination is not generated correctly. output: %v", params)
	}
}
