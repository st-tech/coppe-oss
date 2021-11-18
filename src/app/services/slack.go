package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func PostSlackMessage(message string, channel string) error {
	hookUrl, err := getSlackAlertUrl(channel)
	if err != nil {
		return err
	}

	json, err := json.Marshal(
		SlackTopBlock{
			Blocks: []interface{}{
				SlackParentBlockForSingleLeaf{
					Type: "section",
					Text: SlackLeafBlock{
						Type: "mrkdwn",
						Text: message,
					},
				},
			},
		},
	)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, hookUrl, bytes.NewBuffer(json))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	if res.StatusCode != 200 || buf.String() != "ok" {
		return fmt.Errorf("ERROR: non-ok response returned from slack - %s", buf.String())
	}
	return nil
}

func getSlackAlertUrl(channel string) (string, error) {
	defaultUrlEnv := "SLACK_HOOK_URL"
	if channel != "" {
		specifiedUrlEnv := fmt.Sprintf("%s_%s", defaultUrlEnv, strings.ToUpper(channel))

		specifiedUrl, specifiedFound := os.LookupEnv(specifiedUrlEnv)
		if !specifiedFound {
			defaultUrl, defaultFound := os.LookupEnv(defaultUrlEnv)
			if !specifiedFound && !defaultFound {
				return "", fmt.Errorf("no value set for SLACK_HOOK_URL and %s. Check your environment variable", specifiedUrlEnv)
			}
			return defaultUrl, nil
		}
		return specifiedUrl, nil
	} else {
		defaultUrl, found := os.LookupEnv("SLACK_HOOK_URL")
		if !found {
			return "", errors.New("no value set for SLACK_HOOK_URL. Check your environment variable")
		}
		return defaultUrl, nil
	}
}
