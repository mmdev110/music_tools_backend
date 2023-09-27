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
	//最近のアクセスがあったらdestroy不可
	found, _ := app.findLogWithinThreshold(time.Duration(app.accessThresholdMin))
	if found {
		return false, errors.New("recent access found. backend is not ready to destroy")
	}
	pingSuccess, _ := ping(app.backendEndpoint)
	//ping通らなかったらdestroy不可
	if !pingSuccess {
		return false, errors.New("backend is not active and not ready to destroy")
	}
	dbStatus, _ := app.getDBStatus()
	//DB起動してなかったらdestroy不可
	if dbStatus != DB_STATUS_AVAILABLE {
		return false, errors.New("DB is not active. backend is not ready to destroy")
	}
	return true, nil

}
