package main

import (
	"errors"
	"time"
)

func (app *App) Destroy() (*Response, error) {
	if !app.forceExec {
		can, err := app.canDestroy()
		if err != nil {
			return nil, err
		}
		if !can {
			return nil, errors.New("recent access found. backend is not ready to destroy")
		}
	}
	payload := WorkflowBody{
		EventType:     app.githubEventType,
		ClientPayload: Payload{"Destroy"},
	}
	err := app.requestToGitHub(payload)
	return nil, err
}
func (app *App) canDestroy() (bool, error) {
	found, err := app.findLogWithinThreshold(time.Duration(app.accessThresholdMin))
	//最近のアクセスがあったらdestroyできない
	return !found, err
}
