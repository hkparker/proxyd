package main

import (
	"log"
	"os"
	"encoding/json"
)

const TTPD_CONFIG = "TTPD_CONFIG"

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	AssignLoggingVariables()

	var service_pack TTPDServicePack
	err := json.Unmarshal([]byte(os.Getenv(TTPD_CONFIG)), &service_pack)
	if err != nil {
		LogServiceParseFailed()
		return
	}

	service_pack.RunServices()
}
