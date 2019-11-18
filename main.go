package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/gorilla/mux"
)

type Info struct {
	Name    string
	Version string
	Build   int
}

type ControllerInfo struct {
	StationNames      map[int]string
	AvailableStations []int
	Name              string
	ProgramNames      map[int]string
}

var controllerInfo ControllerInfo

func main() {
	LoadConfig()

	viper.SetDefault("controller.ip", "192.168.1.1")
	viper.SetDefault("rest.port", 8080)

	log.SetLevel(log.DebugLevel)

	log.Info(viper.GetString("controller.ip"))
	log.Info(viper.GetString("controller.key"))
	log.Info(viper.GetString("rest.port"))

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/controller", Controller)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("rest.port")), router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Info{Name: "rainbird-api", Version: "0.0.1", Build: 1})
}

func Controller(w http.ResponseWriter, r *http.Request) {
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

		code, responseData := sipCommand("SerialNumberRequest")

		log.Debug(code)
		log.Debug(responseData)

		err, otherRPCResponse := rpcCommand("getWifiParams", map[string]interface{}{})

		if err != nil {
			log.Error(err)
		}

		log.Debug(otherRPCResponse)

		controllerInfo = ControllerInfo{
			StationNames:      RPCResponse.Result.Controller.StationNames,
			AvailableStations: RPCResponse.Result.Controller.AvailableStations,
			Name:              RPCResponse.Result.Controller.Name,
			ProgramNames:      RPCResponse.Result.Controller.ProgramNames,
		}
	}

	json.NewEncoder(w).Encode(controllerInfo)
}
