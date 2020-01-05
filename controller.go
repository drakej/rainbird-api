package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ControllerInfo struct {
	StationNames      map[string]interface{}
	AvailableStations []interface{}
	Name              string
	ProgramNames      map[string]interface{}
	SerialNumber      string
	Model             string
	Version           string
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

var controllerTime string

var serialNumber SerialNumber

var modelVersion ModelVersion

var firmwareVersion FirmwareVersion

// ESP_RZXe("ESP-RZXe", "0003", a.b.esp_rzxe, false, 0, 6),
// ESP_ME("ESP-Me", "0007", a.b.esp_me, true, 4, 6),
// ST8X_WF("ST8x-WiFi", "0006", a.b.st8wifi, false, 0, 6),
// ESP_TM2("ESP-TM2", "0005", a.b.esp_tm2, true, 3, 4),
// ST8X_WF2("ST8x-WiFi2", "0008", a.b.st8wifi, false, 8, 6),
// ESP_ME3("ESP-ME3", "0009", a.b.esp_me2, true, 4, 6),
// MOCK_ESP_ME2("ESP=Me2", "0010", a.b.esp_me, true, 4, 6),
// ESP_TM2v2("ESP-TM2", "000A", a.b.esp_tm2v2, true, 3, 4),
// ESP_TM2v3("ESP-TM2", "010A", a.b.esp_tm2v3, true, 3, 4),
// TBOS_BT("TBOS-BT", "0099", a.b.tbos_bt_photo, true, 3, 8),
// ESP_MEv2("ESP-Me", "0107", a.b.esp_me, true, 4, 6),
// ESP_RZXev2("ESP-RZXe", "0103", a.b.esp_rzxe, false, 0, 6);

var controllerModels map[int]string = map[int]string{
	-1:  "Unknown",
	3:   "ESP-RZXE",
	5:   "ESP_TM2",
	6:   "ST8x-WiFi",
	7:   "ESP-Me",
	8:   "ST8x-WiFi2",
	9:   "ESP-ME3",
	11:  "ESP-TM2v2",
	16:  "ESP-MEv2 (Mock)",
	153: "TBOS-BT",
	259: "ESP-RZXEv2",
	263: "ESP-MEv2",
	266: "ESP-TM2v3",
}

func ControllerInfoHandler(w http.ResponseWriter, r *http.Request) {
	if controllerInfo.Name == "" {
		log.Info("Retrieving Controller Info from RainBird Cloud")

		if zipCode == 0 {
			zipCode = getZipCode()
		}

		response := cloudRPCCommand("requestWeatherAndStatus", map[string]interface{}{
			"StickId": viper.GetString("controller.mac"),
			"ZipCode": zipCode,
			"Country": "US",
		})

		controllerData := response.Result["Controller"].(map[string]interface{})

		_, serialData := sipCommand("SerialNumberRequest")

		_, modelVersionData := sipCommand("ModelAndVersionRequest")

		model, err := strconv.Atoi(modelVersionData["modelID"])

		if err != nil {
			model = -1
		}

		protocolRevisionMajor, err := strconv.Atoi(modelVersionData["protocolRevisionMajor"])

		if err != nil {
			protocolRevisionMajor = 0
		}

		protocolRevisionMinor, err := strconv.Atoi(modelVersionData["protocolRevisionMinor"])

		if err != nil {
			protocolRevisionMinor = 0
		}

		version := fmt.Sprintf("%d.%d", protocolRevisionMajor, protocolRevisionMinor)

		controllerInfo = ControllerInfo{
			StationNames:      controllerData["customStationNames"].(map[string]interface{}),
			AvailableStations: controllerData["availableStations"].([]interface{}),
			Name:              controllerData["customName"].(string),
			ProgramNames:      controllerData["customProgramNames"].(map[string]interface{}),
			SerialNumber:      serialData["serialNumber"],
			Model:             controllerModels[model],
			Version:           version,
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
			Model:                 model,
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

func ControllerTimeHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Retrieving Controller Time from Local API")

	code, responseData := sipCommand("CurrentDateRequest")

	log.Debug(code)
	log.Debug(responseData)

	// hour := "12"
	// minute := "00"
	// second := "00"
	year, _ := strconv.ParseInt(responseData["year"], 16, 64)
	month, _ := strconv.ParseInt(responseData["month"], 16, 64)
	day, _ := strconv.ParseInt(responseData["day"], 16, 64)

	code, responseData = sipCommand("CurrentTimeRequest")

	log.Debug(code)
	log.Debug(responseData)

	hour, _ := strconv.ParseInt(responseData["hour"], 16, 64)
	minute, _ := strconv.ParseInt(responseData["minute"], 16, 64)
	second, _ := strconv.ParseInt(responseData["second"], 16, 64)
	localLoc, _ := time.LoadLocation("Local")

	datetime := time.Date(int(year), time.Month(int(month)), int(day), int(hour), int(minute), int(second), 0, localLoc)

	json.NewEncoder(w).Encode(datetime.Format(time.UnixDate))
}
