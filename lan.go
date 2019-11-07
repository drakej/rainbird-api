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

type Command struct {
	RequestCode  string `json:"command"`
	ResponseCode string `json:"response"`
	Parameter    int    `json:"parameter"`
	ParameterOne int    `json:"parameterOne"`
	ParameterTwo int    `json:"parameterTwo"`
	Length       int    `json:"length"`
}

type Response struct {
}

var sipIndex map[string]map[string]Command

var sipCommandIndex map[string]Command

var headers map[string]string

func sipCommand(command string, args ...string) {
	loadSipIndex()

	now := time.Now()

	commandData := sipCommandIndex[command]

	data := commandData.RequestCode
	for _, arg := range args {
		data = data + arg
	}

	payload := CloudRPCRequest{
		Id:     int(now.Unix()),
		Method: `tunnelSip`,
		Params: map[string]interface{}{
			"data":   data,
			"length": commandData.Length,
		},
		JsonRPC: `2.0`,
	}

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		log.Error(err)
	}

	encryptedPayload := Encrypt(string(jsonPayload), viper.GetString("controller.key"))

	reader := bytes.NewReader([]byte(encryptedPayload))

	client := &http.Client{}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/stick", viper.GetString("controller.ip")), reader)

	for name, value := range headers {
		req.Header.Add(name, value)
	}

	response, err := client.Do(req)

	if err != nil {
		log.Error(err)
	}

	responseData, _ := ioutil.ReadAll(response.Body)

	log.Info(string(Decrypt(string(responseData), viper.GetString("controller.key"))))
}

func loadSipIndex() {
	if len(sipIndex) == 0 {
		headers = map[string]string{
			"Accept-Language": "en",
			"Accept-Encoding": "gzip, deflate",
			"User-Agent":      "RainBird/2.0 CFNetwork/811.5.4 Darwin/16.7.0",
			"Accept":          "*/*",
			"Connection":      "keep-alive",
			"Content-Type":    "application/octet-stream",
		}

		contents, err := ioutil.ReadFile("sipCommands.json")

		if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal(contents, &sipIndex)

		sipCommandIndex = sipIndex["ControllerCommands"]
	}
}
