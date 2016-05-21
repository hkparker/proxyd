package main

import (
	"encoding/json"
	"flag"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
)

func main() {
	var config_filename string
	flag.StringVar(&config_filename, "config", "proxyd_config.json", "name of configuration file")
	flag.Parse()

	config_data, err := ioutil.ReadFile(config_filename)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err.Error(),
			"config": config_filename,
		}).Error("unable to read configuration data")
		return
	}

	var service_pack ServicePack
	json.Unmarshal(config_data, &service_pack)

	service_pack.run()
}
