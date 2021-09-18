package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
)

const (
	route = "/System/Info/Public"
)

// SystemInfo is response returned from /System/Info/Public endpoint
type SystemInfo struct {
	LocalAddress           string `json:"LocalAddress"`
	ServerName             string `json:"ServerName"`
	Version                string `json:"Version"`
	ProductName            string `json:"ProductName"`
	OperatingSystem        string `json:"OperatingSystem"`
	ID                     string `json:"Id"`
	StartupWizardCompleted bool   `json:"StartupWizardCompleted"`
}

// systemInfoHandler returns modified system information
func systemInfoHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
	info, err := getSystemInfo()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("original info: %+v", info)
	overrideSystemInfo(info)
	log.Printf("modified info: %+v", info)
	if err := json.NewEncoder(w).Encode(info); err != nil {
		log.Printf("failed to respond with system information")
	}
}

// getSystemInfo retrieves system information from internal jellyfin service
func getSystemInfo() (*SystemInfo, error) {
	info := &SystemInfo{}
	resp, err := http.Get(internalServer + route)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(info); err != nil {
		return nil, err
	}
	return info, nil
}

// overrideSystemInfo overrides LocalAddress of given system information struct
func overrideSystemInfo(info *SystemInfo) {
	info.LocalAddress = externalServer
}

// flags
var (
	internalServer string
	externalServer string
	port           string
)

func init() {
	flag.StringVar(&internalServer, "internal", lookupEnv("INTERNAL", "jellyfin"), "internal jellyfin service")
	flag.StringVar(&externalServer, "external", lookupEnv("EXTERNAL", ""), "external jellyfin service")
	flag.StringVar(&port, "port", lookupEnv("PORT", "8080"), "port")
}

func lookupEnv(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func main() {
	flag.Parse()

	http.HandleFunc(route, systemInfoHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("internal server: %s", internalServer)
	log.Printf("external server: %s", externalServer)
	info, err := getSystemInfo()
	if err != nil {
		log.Printf("ERROR: %+v", err)
		return
	}
	log.Printf("internal json: %+v", info)
	overrideSystemInfo(info)
	log.Printf("external json: %+v", info)

	http.ListenAndServe(":"+port, nil)
}
