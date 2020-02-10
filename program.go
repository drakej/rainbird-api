package main

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var soilTypes map[int]string = map[int]string{
	0: "None",
	1: "Clay",
	2: "Sand",
	3: "Other",
}

// ProgramInfoHandler returns information about the request program
func ProgramInfoHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Request Program Info")

	responseData, code := rpcCommand("getProgramInfo", map[string]interface{}{})

	log.Debug(code)
	log.Debug(responseData)
	programSoilTypes := responseData.Result["SoilTypes"].([]interface{})

	programs := map[int]map[string]interface{}{
		0: {
			"FlowRate": 0,
			"FlowUnit": 0,
			"SoilType": "",
		},
		1: {
			"FlowRate": 0,
			"FlowUnit": 0,
			"SoilType": "",
		},
		2: {
			"FlowRate": 0,
			"FlowUnit": 0,
			"SoilType": "",
		},
		3: {
			"FlowRate": 0,
			"FlowUnit": 0,
			"SoilType": "",
		},
	}

	//soilTypeList := make([]string, 0)

	log.Debug(programSoilTypes)

	for index, soilType := range programSoilTypes {
		log.Debug(soilType)

		soilTypeInt := int(soilType.(float64))

		programs[index]["SoilType"] = soilTypes[soilTypeInt]

		//soilTypeList = append(soilTypeList, soilTypes[index])
	}

	json.NewEncoder(w).Encode(programs)
}
