package main

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func ZipCodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting Zip Code and Country")

	code, responseData := rpcCommand("getZipCode", map[string]interface{}{})

	log.Debug(code)
	log.Debug(responseData)

	json.NewEncoder(w).Encode(responseData.Result)
}

func RainSensorHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting Rain Sensor State")

	code, responseData := sipCommand("CurrentRainSensorStateRequest")

	log.Debug(code)
	log.Debug(responseData)

	json.NewEncoder(w).Encode(responseData)
}
