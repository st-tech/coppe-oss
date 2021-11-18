package services

import (
	"os"
	"testing"
	"time"
)

func TestTimezone(t *testing.T) {
	timezone, found := os.LookupEnv("TIMEZONE")
	if !found {
		t.Error("env for TIMEZONE not found")
	}
	_, err := time.LoadLocation(timezone)
	if err != nil {
		t.Error(err)
	}
}
