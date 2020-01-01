package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ControllerInfo struct {
	StationNames      map[int]string
	AvailableStations []int
	Name              string
	ProgramNames      map[int]string
}

type SerialNumber struct {
	SerialNumber string
}

type ModelVersion struct {
	Model                 int
	ProtocolRevisionMajor int
	ProtocolRevisionMinor int
}

type FirmwareVersion struct {
	FirmwareRevisionMajor int64
	FirmwareRevisionMinor int64
}

var controllerInfo ControllerInfo

var serialNumber SerialNumber

var modelVersion ModelVersion

var firmwareVersion FirmwareVersion

func ControllerInfoHandler(w http.ResponseWriter, r *http.Request) {
	if controllerInfo.Name == "" {
		log.Info("Retrieving Controller Info from RainBird Cloud")
		requestData := RPCRequest{
			Id:     int(time.Now().Unix()),
			Method: "requestWeatherAndStatus",
			Params: map[string]interface{}{
				"StickId": viper.GetString("controller.mac"),
				"ZipCode": viper.GetString("controller.zipcode"),
				"Country": "US",
			},
			JsonRPC: "2.0",
		}

		jsonData, err := json.Marshal(requestData)

		if err != nil {
			log.Error("Could not marshal json for cloud status request")
		}

		reader := bytes.NewReader(jsonData)

		resp, err := http.Post(fmt.Sprintf("http://%s/phone-api", viper.GetString("rainbirdcloud.host")), "application/json", reader)

		if err != nil {
			log.Error(err)
		}

		var RPCResponse CloudRPCResponse

		respBody, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Error(err)
		}

		log.Debug(string(respBody))

		err = json.Unmarshal(respBody, &RPCResponse)

		if err != nil {
			log.Error(err)
		}

		controllerInfo = ControllerInfo{
			StationNames:      RPCResponse.Result.Controller.StationNames,
			AvailableStations: RPCResponse.Result.Controller.AvailableStations,
			Name:              RPCResponse.Result.Controller.Name,
			ProgramNames:      RPCResponse.Result.Controller.ProgramNames,
		}
	}

	json.NewEncoder(w).Encode(controllerInfo)
}

func ControllerSerialNumberHandler(w http.ResponseWriter, r *http.Request) {
	if serialNumber.SerialNumber == "" {
		log.Info("Retrieving Controller Serial Number from Local API")

		code, responseData := sipCommand("SerialNumberRequest")

		log.Debug(code)
		log.Debug(responseData)

		serialNumber = SerialNumber{
			SerialNumber: responseData["serialNumber"],
		}
	}

	json.NewEncoder(w).Encode(serialNumber)
}

func ControllerModelVersionHandler(w http.ResponseWriter, r *http.Request) {
	if modelVersion.Model == 0 {
		log.Info("Retrieving Controller Model & Protocol Version from Local API")

		code, responseData := sipCommand("ModelAndVersionRequest")

		log.Debug(code)
		log.Debug(responseData)

		model, err := strconv.Atoi(responseData["modelID"])

		if err != nil {
			model = 0
		}

		protocolRevisionMajor, err := strconv.Atoi(responseData["protocolRevisionMajor"])

		if err != nil {
			protocolRevisionMajor = 0
		}

		protocolRevisionMinor, err := strconv.Atoi(responseData["protocolRevisionMinor"])

		if err != nil {
			protocolRevisionMinor = 0
		}

		modelVersion = ModelVersion{
			Model: model,
			ProtocolRevisionMajor: protocolRevisionMajor,
			ProtocolRevisionMinor: protocolRevisionMinor,
		}
	}

	json.NewEncoder(w).Encode(modelVersion)
}

func ControllerFWVersionHandler(w http.ResponseWriter, r *http.Request) {
	if firmwareVersion.FirmwareRevisionMajor == 0 {
		log.Info("Retrieving Controller Firmware Version from Local API")

		code, responseData := sipCommand("ControllerFirmwareVersionRequest")

		log.Debug(code)
		log.Debug(responseData)

		firmwareRevisionMajor, err := strconv.ParseInt(responseData["firmwareRevisionMajor"], 16, 64)

		if err != nil {
			log.Warning(err)
			firmwareRevisionMajor = 0
		}

		firmwareRevisionMinor, err := strconv.ParseInt(responseData["firmwareRevisionMinor"], 16, 64)

		if err != nil {
			log.Warning(err)
			firmwareRevisionMinor = 0
		}

		firmwareVersion = FirmwareVersion{
			FirmwareRevisionMajor: firmwareRevisionMajor,
			FirmwareRevisionMinor: firmwareRevisionMinor,
		}
	}

	json.NewEncoder(w).Encode(firmwareVersion)
}
