package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"example.com/app/customError"
)

func ResponseJSON(w http.ResponseWriter, payload interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	//fmt.Println("payload:")
	//PrintStruct(payload)
	js, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		fmt.Println(err)
	}
	w.Write(js)
}

//func ErrorJSON(w http.ResponseWriter, err error, statuses ...int) {
//	status := http.StatusBadRequest
//	if len(statuses) > 0 {
//		status = statuses[0]
//	}
//	payload := struct {
//		Message string `json:"message"`
//	}{err.Error()}
//	log.Println(err.Error())
//	ResponseJSON(w, payload, status)
//}

// add err to overwrite message
func ErrorJSON(w http.ResponseWriter, customError customError.CustomError, err error) {
	status := http.StatusBadRequest
	if err != nil {
		customError.Message = err.Error()
	}
	payload := customError
	log.Println(customError.Error())
	ResponseJSON(w, payload, status)
}

func BodyToString(body io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	bytes := buf.String()
	return string(bytes)
}
func BodyToStruct(body io.Reader, payload any) error {
	if err := json.NewDecoder(body).Decode(payload); err != nil {
		return err
	}
	return nil
}

func ToJSON(payload interface{}) (string, error) {
	b, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
