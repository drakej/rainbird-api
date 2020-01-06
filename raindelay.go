package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func RainDelayHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Requesting Rain Delay setting")

	code, responseData := sipCommand("RainDelayGetRequest")

	log.Debug(code)
	log.Debug(responseData)

	var delayDays int64 = 0

	delayDays, _ = strconv.ParseInt(responseData["delaySetting"], 16, 16)

	json.NewEncoder(w).Encode(map[string]int64{"RainDelayDays": delayDays})
}

func RainDelaySetHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Setting Rain Delay setting")

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 4))

	var rainDelay int

	if err != nil {
		log.Error("Request had something other than a valid rain delay integer")
	}

	if err := r.Body.Close(); err != nil {
		log.Error("Request had something other than a valid rain delay integer")
	}

	if err := json.Unmarshal(body, &rainDelay); err != nil {
		log.Error("Request had something other than a valid rain delay integer")
	}

	code, responseData := sipCommand("RainDelaySetRequest", string(rainDelay))

	log.Debug(code)
	log.Debug(responseData)

	codeInt := 0

	codeInt, _ = strconv.Atoi(code)

	json.NewEncoder(w).Encode(map[string]int{"responseCode": codeInt})
}
