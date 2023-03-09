package main

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func ZipCodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting Zip Code and Country")

	response, err := rpcCommand("getZipCode", map[string]interface{}{})

	if err != nil {
		log.Error(err)
	}

	json.NewEncoder(w).Encode(response.Result)
}

func RainSensorHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting Rain Sensor State")

	code, responseData := sipCommand("CurrentRainSensorStateRequest")

	log.Debug(code)
	log.Debug(responseData)

	json.NewEncoder(w).Encode(responseData)
}

func SeasonalAdjustHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting current seasonal adjustment factor")

	code, responseData := sipCommand("ZonesSeasonalAdjustFactorRequest")

	log.Debug(code)
	log.Debug(responseData)

	json.NewEncoder(w).Encode(responseData)
}
