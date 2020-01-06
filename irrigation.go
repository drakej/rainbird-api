package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

var irrigationState bool

func IrrigationStateHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Retrieving Irrigation State from Local API")

	code, responseData := sipCommand("CurrentIrrigationStateRequest")

	log.Debug(code)
	log.Debug(responseData)

	irrigationState = false

	if responseData["irrigationState"] != "00" {
		irrigationState = true
	}

	json.NewEncoder(w).Encode(irrigationState)
}

func IrrigationStopHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Sending StopIrrigationRequest")

	code, responseData := sipCommand("StopIrrigationRequest")

	log.Debug(code)
	log.Debug(responseData)

	codeInt := -1

	codeInt, _ = strconv.Atoi(code)

	json.NewEncoder(w).Encode(map[string]int{"resultCode": codeInt})
}
