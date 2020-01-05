package main

import (
	"encoding/json"
	"net/http"

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
