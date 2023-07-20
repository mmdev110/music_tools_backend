package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"example.com/app/customError"
)

func ResponseJSON(w http.ResponseWriter, payload interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Println("payload:")
	PrintStruct(payload)
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

// add err to customize message
func ErrorJSON(w http.ResponseWriter, customError customError.CustomError, err ...error) {
	status := http.StatusBadRequest
	if len(err) > 0 {
		customError.Message = err[0].Error()
	}
	payload := customError
	log.Println(customError.Error())
	ResponseJSON(w, payload, status)
}
