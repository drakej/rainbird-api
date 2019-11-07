package main

import "encoding/json"

type CloudRPCRequest struct {
	Id      int                    `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	JsonRPC string                 `json:"jsonrpc"`
}

type CloudRPCResponse struct {
	Id      int         `json:"id"`
	Result  CloudResult `json:"result"`
	JsonRPC string      `json:"jsonrpc"`
}

type CloudResult struct {
	Weather         Weather         `json:"Weather"`
	ForecastedRain  ForecastedRain  `json:"ForecastedRain"`
	Controller      CloudController `json:"Controller"`
	ConnectedStatus []CloudService  `json:"ConnectedStatus"`
}

type Weather struct {
	City              string      `json:"city"`
	TimeZoneId        string      `json:"timeZoneId"`
	Forecasts         []Forecast  `json:"forecast"`
	Location          string      `json:"location"`
	TimeZoneRawOffset json.Number `json:"timeZoneRawOffset"`
}

type Forecast struct {
	DateTime     json.Number `json:"dateTime"`
	High         json.Number `json:"high"`
	ChanceOfRain json.Number `json:"chaneofrain"`
	Precip       json.Number `json:"precip"`
	Low          json.Number `json:"low"`
	Icon         string      `json:"icon"`
	Description  string      `json:"description"`
}

type ForecastedRain struct {
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
