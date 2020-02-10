package main

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// WifiAPModeHandler responds to /wifi/ap GET requests
func WifiAPModeHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Retrieving Wifi APMode")

	responseData, code := rpcCommand("getApMode", map[string]interface{}{})

	log.Debug(code)
	log.Debug(responseData)

	json.NewEncoder(w).Encode(responseData.Result)
}

// WifiConfigHandler responds to /wifi/config GET requests
func WifiConfigHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Retrieving Wifi Configuration from Local API")

	responseData, code := rpcCommand("getWifiParams", map[string]interface{}{})

	log.Debug(code)
	log.Debug(responseData)

	// Let's not broadcast this as it just weakens the user's security
	delete(responseData.Result, "wifiPassword")

	json.NewEncoder(w).Encode(responseData.Result)
}

// WifiIPHandler responds to /wifi/ipv4 GET requests
func WifiIPHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Retrieving Wifi IPv4 Address")

	wifiIPAddress := getStickIPAddress()

	json.NewEncoder(w).Encode(wifiIPAddress)
}

// WifiNetworkHandler responds to /wifi/network GET requests
func WifiNetworkHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Retrieving Wifi Network Status")

	responseData, code := rpcCommand("getNetworkStatus", map[string]interface{}{})

	log.Debug(code)
	log.Debug(responseData)

	json.NewEncoder(w).Encode(responseData.Result)
}
