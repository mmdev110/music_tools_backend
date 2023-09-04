package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func ping(endpoint string) (bool, error) {
	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	client := http.Client{}
	_, err := client.Do(req)

	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (app *App) requestToGitHub(payload WorkflowBody) error {
	js, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, app.githubEndpoint, strings.NewReader(string(js)))
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+app.githubToken)
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	client := http.Client{}
	if !app.sendRequest { //テスト等の場合は送らない
		return nil
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("request failed with code: %d", res.StatusCode)
	}
	return nil
}
