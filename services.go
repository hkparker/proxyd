package main

import (
	"crypto/tls"
	log "github.com/Sirupsen/logrus"
	"strings"
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
	log.WithFields(log.Fields{
		"service_count": len(service_pack),
	}).Info("proxyd starting")

	listener_failed := make(chan error)
	for _, service_config := range service_pack {
		go listenAndProxy(service_config, listener_failed)
	}
	for i := 0; i < len(service_pack); i++ {
		service_error := <-listener_failed
		log.WithFields(log.Fields{
			"error": service_error,
		}).Warn("service died")
	}

	log.Fatal("all services died")
}

func listenAndProxy(config ServiceConfig, failed chan error) {
	front_tls_config, err := populateTLSConfig(config.FrontConfig)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("error populating front TLS configuration")
	}
	back_tls_config, err := populateTLSConfig(config.BackConfig)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("error populating back TLS configuration")
	}

	listener, err := listenAny(config.Front, front_tls_config)
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
			go proxyBack(conn, config.Back, back_tls_config)
		}
	}
}

func populateTLSConfig(tls_config TLSConfig) (tls.Config, error) {
	config := tls.Config{}
	if len(tls_config) == 0 {
		return config, nil
	}

	cert_data := ""
	key_data := ""
	insecure := false
	for tls_config_key, value := range tls_config {
		if tls_config_key == "CERT" {
			cert_data = value
		} else if tls_config_key == "KEY" {
			key_data = value
		} else if tls_config_key == "InsecureSkipVerify" {
			insecure = value == "true"
		} else if tls_config_key == "NextProtos" {
			config.NextProtos = strings.Split(value, ",")
		}
	}

	if cert_data != "" && key_data != "" {
		certificate, err := tls.X509KeyPair(
			[]byte(cert_data),
			[]byte(key_data),
		)
		if err != nil {
			return config, err
		}
		config.Certificates = []tls.Certificate{certificate}
	}
	if insecure {
		config.InsecureSkipVerify = true
	}
	// root CAs
	// cypher suites
	// PreferServerCipherSuites
	// SessionTicketsDisabled
	// SessionTicketKey
	// MinVersion
	// MaxVersion
	// CurvePreferences
	// DynamicRecordSizingDisabled
	return config, nil
}
