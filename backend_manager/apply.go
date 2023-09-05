package main

import "errors"

func (app *App) Apply() (*Response, error) {
	if !app.forceExec {
		can, err := app.canApply()
		if err != nil {
			return nil, err
		}
		if !can {
			return nil, errors.New("backend is not ready to apply")
		}
	}
	payload := WorkflowBody{
		EventType:     app.githubEventType,
		ClientPayload: Payload{"Apply"},
	}
	err := app.requestToGitHub(payload)
	return nil, err
}
func (app *App) canApply() (bool, error) {
	pingSuccess, err := ping(app.backendEndpoint)
	if err != nil {
		return false, err
	}
	dbStatus, err := app.getDBStatus()
	if err != nil {
		return false, err
	}
	//落ちてて、DB止まってるならapplyしてOK
	return !pingSuccess && dbStatus == "stopped", nil
}
