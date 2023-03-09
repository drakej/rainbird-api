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

// RPCRequest for JSONRPC
type RPCRequest struct {
	ID      int                    `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	JSONRPC string                 `json:"jsonrpc"`
}

// cloudRPCCommand sends a JSON RPC command request to the RB cloud
func cloudRPCCommand(method string, params map[string]interface{}) RPCResponse {
	requestData := RPCRequest{
		ID:      int(time.Now().Unix()),
		Method:  method,
		Params:  params,
		JSONRPC: "2.0",
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

// getStickIPAddress returns the local IP of the WiFi module by requesting it from the RB cloud
func getStickIPAddress() string {
	response := cloudRPCCommand("getStickIpAddress", map[string]interface{}{
		"StickId": viper.GetString("controller.mac"),
	})

	stickIP := ""

	if response.Result["IpAddress"] != nil {
		stickIP = response.Result["IpAddress"].(string)
	}

	viper.Set("controller.ip", stickIP)

	return stickIP
}
