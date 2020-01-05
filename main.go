package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/gorilla/mux"
)

type Info struct {
	Name    string
	Version string
	Build   int
}

func main() {
	LoadConfig()

	viper.SetDefault("controller.ip", "192.168.1.1")
	viper.SetDefault("rest.port", 8080)

	log.SetLevel(log.DebugLevel)

	log.Info(viper.GetString("controller.ip"))
	log.Info(viper.GetString("controller.key"))
	log.Info(viper.GetString("rest.port"))

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", Index).Methods("GET")
	r.HandleFunc("/controller", ControllerInfoHandler).Methods("GET")

	// Serial Number of the controller (or maybe the WiFi module? should confirm), comes from local API
	r.HandleFunc("/controller/serial", ControllerSerialNumberHandler).Methods("GET")

	// The model (ESP-Me, ESP-RZXe, ST8X-Wf, etc.) and the version (major and minor e.g. 2.1), comes from local API
	// Note: This is from manfacturing and I believe will never change. RB official apps seem to ignore the version data anyways
	r.HandleFunc("/controller/model", ControllerModelVersionHandler).Methods("GET")

	// This is the firmware of the controller itself (I'm guessing though since there's no reference to confirm)
	// This is different than the model itself with it's own version from manufacturing
	r.HandleFunc("/controller/firmware", ControllerFWVersionHandler).Methods("GET")

	// The currently configured time on the controller
	r.HandleFunc("/controller/time", ControllerTimeHandler).Methods("GET")

	// Whether irrigation is currently running
	r.HandleFunc("/irrigation/state", IrrigationStateHandler).Methods("GET")

	r.HandleFunc("/wifi/ip", WifiIPHandler).Methods("GET")

	// Wifi configuration setup in the RainBird app initially
	r.HandleFunc("/wifi/config", WifiConfigHandler).Methods("GET")

	// Wifi Network Status
	r.HandleFunc("/wifi/network", WifiNetworkHandler).Methods("GET")

	// The actual stick has it's own firmware which is reported via cloud for some reason and parsed out
	// Note: I believe there's only one version of the Wifi module, so right now the FW is at 1.41 as I write this
	//router.HandleFunc("/wifimodule", WiFiModuleFWVersion)

	r.Use(LoggingMiddleware)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("rest.port")), r))
}

// LoggingMiddleware is responsible for logging mux requests for the REST api
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Do stuff here
		log.Infof("%s - %s - %s %s %s %s", r.Host, r.RemoteAddr, r.Method, r.RequestURI, r.Proto, r.UserAgent())
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func Index(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Info{Name: "rainbird-api", Version: "0.0.1"})
}
