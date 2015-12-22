package main

import (
	"net"
	"io"
	"errors"
	"crypto/tls"
	"os"
)

func PopulateTLSConfig(config TLSConfig) tls.Config {
	config := tls.Config{}
	cert_data := ""
	key_data := ""
	for tls_config_key, envar_name := range(config) {
		if tls_config_key == "CERT" {
			cert_data = os.Getenv(envar_name)
		} else if tls_config_key == "KEY" {
			key_data := os.Getenv(envar_name)
		}
	}
	certificate, err := tls.X509KeyPair(
		[]byte(cert_data),
		[]byte(key_data),
	)
	config.Certificates = []tls.Certificate{certificate}
	return config
}

func ListenEither(addr string, config TLSConfig) (net.Listener, error) {
	proto := addr[:6]
	address := addr[6:]
	if proto == "tls://" {
		tls_config := PopulateTLSConfig(config)
		return tls.Listen("tcp", address, tls_config)
	} else if proto == "tcp://" {
		return net.Listen("tcp", address)
	}
	return nil, errors.New("unrecognized protocol")
}

func DialEither(addr string, config TLSConfig) (net.Conn, error) {
	proto := addr[:6]
	address := addr[6:]
	if proto == "tls://" {
		tls_config := PopulateTLSConfig(config)
		return tls.Dial("tcp", address, config)
	} else if proto == "tcp://" {
		return net.Dial("tcp", address)
	}
	return nil, errors.New("unrecognized protocol")
}

func ExchangeData(external, internal net.Conn) {
	errs := make(chan error, 2)

	go func() {
		_, err := io.Copy(external, internal)
		errs <- err
	}()
	go func() {
		_, err := io.Copy(internal, external)
		errs <- err
	}()

	<-errs
	external.Close()
	internal.Close()
}

func ProxyBack(external net.Conn, addr string, config TLSConfig) {
	internal, err := DialEither(addr, config)
	if err != nil {
		if service_up[addr] {
			service_up[addr] = false
			LogServiceDown(addr, err)
		}
		external.Close()
		return
	} else if !service_up[addr] {
		service_up[addr] = true
		LogServiceRecovered(addr)
	}
	ExchangeData(external, internal)
}
