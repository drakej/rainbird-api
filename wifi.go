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
