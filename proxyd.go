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
		}).Error("unable to read configuration file")
		return
	}

	var service_pack ServicePack
	err = json.Unmarshal(config_data, &service_pack)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err.Error(),
			"config": config_filename,
		}).Error("unable to parse configuration data")
		return
	}

	service_pack.run()
}
