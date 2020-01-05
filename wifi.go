package main

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var wifiParams map[string]interface{}

func WifiConfigHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Retrieving Wifi Configuration from Local API")

	code, responseData := rpcCommand("getWifiParams", map[string]interface{}{})

	log.Debug(code)
	log.Debug(responseData)

	// Let's not broadcast this as it just weakens the user's security
	delete(responseData.Result, "wifiPassword")

	json.NewEncoder(w).Encode(responseData.Result)
}

func WifiIPHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Retrieving Wifi IPv4 Address")

	wifiIPAddress := getStickIPAddress()

	json.NewEncoder(w).Encode(wifiIPAddress)
}

func WifiNetworkHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Retrieving Wifi Network Status")

	code, responseData := rpcCommand("getNetworkStatus", map[string]interface{}{})

	log.Debug(code)
	log.Debug(responseData)

	json.NewEncoder(w).Encode(responseData.Result)
}
