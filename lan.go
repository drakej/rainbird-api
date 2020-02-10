package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Command is a SIP command object
type Command struct {
	Length   int    `json:"length"`
	Command  string `json:"command"`
	Response string `json:"response"`
}

// Response from SIP Command
type Response struct {
	Length int              `json:"length"`
	Type   string           `json:"type"`
	Params map[string]Param `json:"params"`
}

// RPCResponse from JSON RPC command
type RPCResponse struct {
	ID      int                    `json:"id"`
	Result  map[string]interface{} `json:"result"`
	JSONRPC string                 `json:"jsonrpc"`
}

// SIPConfig represents a lookup of SIP commands
type SIPConfig struct {
	Loaded    bool
	Commands  map[string]Command  `json:"ControllerCommands"`
	Responses map[string]Response `json:"ControllerResponses"`
}

// Param for  SIP Commands
type Param struct {
	Position int `json:"position"`
	Length   int `json:"length"`
}

var sipIndex SIPConfig

var headers = map[string]string{
	"Accept-Language": "en",
	"Accept-Encoding": "gzip, deflate",
	"User-Agent":      "RainBird/2.0 CFNetwork/811.5.4 Darwin/16.7.0",
	"Accept":          "*/*",
	"Connection":      "keep-alive",
	"Content-Type":    "application/octet-stream",
}

var localIPv4Address string

var zipCode int
var country string

func rpcCommand(method string, params map[string]interface{}) (RPCResponse, error) {
	now := time.Now()

	payload := RPCRequest{
		ID:      int(now.Unix()),
		Method:  method,
		Params:  params,
		JSONRPC: `2.0`,
	}

	log.Debug(payload)

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		log.Error(err)
	}

	encryptedPayload := Encrypt(string(jsonPayload), viper.GetString("controller.key"))

	reader := bytes.NewReader([]byte(encryptedPayload))

	client := &http.Client{}

	if viper.GetString("controller.ip") == "" {
		getStickIPAddress()
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/stick", viper.GetString("controller.ip")), reader)

	for name, value := range headers {
		req.Header.Add(name, value)
	}

	response, err := client.Do(req)

	log.Debug(response)

	if err != nil {
		log.Error(err)
	}

	encryptedResponse, _ := ioutil.ReadAll(response.Body)

	var rpcResponse RPCResponse

	if response.StatusCode == 200 {
		rpcRawResponse := strings.TrimRight(Decrypt(string(encryptedResponse), viper.GetString("controller.key")), "\x00\x0A\x10")

		err = json.Unmarshal([]byte(rpcRawResponse), &rpcResponse)

		if err != nil {
			log.Error(err)
		}
	} else if response.StatusCode == 503 {
		return rpcResponse, fmt.Errorf("rpcResponse: Error code %d received, the controller is most likely locked/busy", response.StatusCode)
	}

	return rpcResponse, nil
}

func getZipCode() (int, string) {
	response, err := rpcCommand("getZipCode", map[string]interface{}{})

	if err != nil {
		return -1, ""
	}

	zipCode, _ := strconv.Atoi(response.Result["code"].(string))

	return zipCode, response.Result["country"].(string)
}

func sipCommand(command string, args ...string) (string, map[string]string) {
	loadSipIndex()

	commandData := sipIndex.Commands[command]
	cmdResponse := make(map[string]string)

	data := commandData.Command
	for _, arg := range args {
		data = data + arg
	}

	rpcResponse, err := rpcCommand(`tunnelSip`, map[string]interface{}{
		"data":   data,
		"length": commandData.Length,
	})

	if err != nil {
		return "0", cmdResponse
	}

	responseData := ""
	responseCode := "00"
	responseType := sipIndex.Responses["00"]

	if rpcResponse.Result["data"] != nil {
		responseData = rpcResponse.Result["data"].(string)

		responseCode = responseData[:2]

		responseType = sipIndex.Responses[responseCode]

		for name, param := range responseType.Params {
			cmdResponse[name] = responseData[param.Position : param.Position+param.Length]
		}
	}

	return responseCode, cmdResponse
}

func loadSipIndex() {
	if !sipIndex.Loaded {
		contents, err := ioutil.ReadFile("sipCommands.json")

		if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal(contents, &sipIndex)

		sipIndex.Loaded = true
	}
}
