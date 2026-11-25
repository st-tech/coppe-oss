package services

import (
	"encoding/json"
	"fmt"
)

func SendInfoLog(msg string) {
	sendLog("info", msg)
}

func SendDebugLog(msg string) {
	sendLog("debug", msg)
}

func SendWarningLog(msg string) {
	sendLog("warning", msg)
}

func SendErrorLog(msg string) {
	sendLog("error", msg)
}

type logEntry struct {
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

func sendLog(severity string, msg string) {
	b, err := json.Marshal(logEntry{Message: msg, Severity: severity})
	if err != nil {
		fmt.Println(msg)
		return
	}
	fmt.Println(string(b))
}
