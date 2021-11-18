package services

import (
	"os"
	"testing"
)

func TestSlackHookUrlFromEnv(t *testing.T) {
	_, notEmpty := os.LookupEnv("SLACK_HOOK_URL")
	if !notEmpty {
		t.Error("Env var: SLACK_HOOK_URL is empty")
	}
}
