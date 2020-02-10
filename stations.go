package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func extractStationList(stationString string) []int {
	// Based on code from rainbirdlib on Android
	var j int64 = 0
	var i int = 0
	var j2 int64 = 0

	for i < len(stationString)-1 {
		i2 := i + 1
		intValue, _ := strconv.ParseInt(string(stationString[i]), 16, 8)
		int2Value, _ := strconv.ParseInt(string(stationString[i2]), 16, 8)

		byteValue := byte(intValue)
		byte2Value := byte(int2Value)

		i += 2
		// FrameManager.CMD_WR_PRG = 15
		j2 |= int64((byte2Value&15)|((byteValue<<4)&240)) << int(j)
		j += 8
	}

	stationsList := make([]int, 0)

	for i3 := 0; i3 < 32; i3++ {
		special := 1 << i3
		if int64(special)&j2 == int64(special) {
			result := int(i3 + 1)
			stationsList = append(stationsList, result)
		}
	}

	log.Debug(stationsList)

	return stationsList
}

// StationsAvailableHandler returns available stations for /stations/available
func StationsAvailableHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Retrieving Available Stations")

	code, responseData := sipCommand("AvailableStationsRequest", "00")

	log.Debug(code)
	log.Debug(responseData)

	stationsList := extractStationList(responseData["setStations"])

	json.NewEncoder(w).Encode(stationsList)
}

// StationsActiveHandler returns the active stations
func StationsActiveHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Retrieving Available Stations")

	code, responseData := sipCommand("CurrentStationsActiveRequest", "00")

	log.Debug(code)
	log.Debug(responseData)

	stationsList := extractStationList(responseData["activeStations"])

	json.NewEncoder(w).Encode(stationsList)
}
