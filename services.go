package main

import (
	"crypto/tls"
	log "github.com/Sirupsen/logrus"
)

type TLSConfig map[string]string

type ServiceConfig struct {
	Front       string
	Back        string
	FrontConfig TLSConfig
	BackConfig  TLSConfig
}

type ServicePack []ServiceConfig

func (service_pack ServicePack) run() {
	log.Info("proxyd starting")

	listener_failed := make(chan error)
	for _, service_config := range service_pack {
		go listenAndProxy(service_config, listener_failed)
	}
	for _ = range service_pack {
		service_error := <-listener_failed
		log.WithFields(log.Fields{
			"error": service_error,
		}).Warn("service died")
	}

	log.Warn("all services died")
}

func listenAndProxy(config ServiceConfig, failed chan error) {

	tls_config, err := populateTLSConfig(config.FrontConfig)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("error populating TLS configuration")
	}

	listener, err := listenAny(config.Front, tls_config)
	if err != nil {
		failed <- err
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Info("unable to accept connection")
			failed <- err
			return
		} else {

			tls_config, err := populateTLSConfig(config.BackConfig)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Error("error populating TLS configuration")
			}

			go proxyBack(conn, config.Back, tls_config)
		}
	}
}

func populateTLSConfig(tls_config TLSConfig) (tls.Config, error) {
	config := tls.Config{}
	cert_data := ""
	key_data := ""
	for tls_config_key, _ := range tls_config {
		if tls_config_key == "CERT" {
		} else if tls_config_key == "KEY" {
		}
	}
	certificate, err := tls.X509KeyPair(
		[]byte(cert_data),
		[]byte(key_data),
	)
	if err != nil {
		return config, err
	}
	config.Certificates = []tls.Certificate{certificate}
	return config, nil
}
