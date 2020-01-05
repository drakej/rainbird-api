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
)

type RPCRequest struct {
	Id      int                    `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	JsonRPC string                 `json:"jsonrpc"`
}

type CloudController struct {
	StationNames      map[int]string `json:"customStationNames"`
	AvailableStations []int          `json:"availableStations"`
	Name              string         `json:"customName"`
	ProgramNames      map[int]string `json:"customProgramNames"`
}

type ConnectedStatus struct {
}

type CloudService struct {
	Status    int               `json:"Status"`
	CompanyId int               `json:"CompanyId"`
	Enabled   bool              `json:"Enabled"`
	Params    map[string]string `json:"Params"`
	Name      string            `json:"Name"`
}

func cloudRPCCommand(method string, params map[string]interface{}) RPCResponse {
	requestData := RPCRequest{
		Id:      int(time.Now().Unix()),
		Method:  method,
		Params:  params,
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

	var rpcResponse RPCResponse

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Error(err)
	}

	log.Debug(string(respBody))

	err = json.Unmarshal(respBody, &rpcResponse)

	return rpcResponse
}
