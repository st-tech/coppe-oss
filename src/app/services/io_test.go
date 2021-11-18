package services

import (
	"testing"
)

func TestReadFiles(t *testing.T) {
	_, err := ReadFiles("../../../yaml")
	if err != nil {
		t.Error("ERROR: file retrieval failed.", err)
	}
}
