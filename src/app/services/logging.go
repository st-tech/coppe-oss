package services

import (
	"fmt"
	"strings"
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

func sendLog(severity string, msg string) {
	str := fmt.Sprintf(`{"message": "%s", "severity": "%s"}`, strings.ReplaceAll(msg, "\n", " "), severity)
	fmt.Println(str)
}
