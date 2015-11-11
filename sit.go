package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/bugsnag/bugsnag-go"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

const FRONT_SERVICE = "FRONT_SERVICE"
const BACK_SERVICE = "BACK_SERVICE"
const CERT_PEM = "CERT_PEM"
const CERT_KEY = "CERT_KEY"
const BUGSNAG_ENDPOINT = "BUGSNAG_ENDPOINT"
const BUGSNAG_API_KEY = "BUGSNAG_API_KEY"
const ENVIRONMENT = "ENVIRONMENT"

func Snag(msg string) {
	log.Println(msg)
	bugsnag.Notify(
		errors.New(msg),
	)
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

func ProxyBack(external net.Conn, addr string) {
	internal, err := net.Dial("tcp", addr)
	if err != nil {
		Snag(fmt.Sprintf(
			"could not contact internal service %s: %v",
			os.Getenv(BACK_SERVICE),
			err,
		))
		external.Close()
		return
	}
	ExchangeData(external, internal)
}

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}
	bugsnag.Configure(bugsnag.Configuration{
		Endpoint:     os.Getenv(BUGSNAG_ENDPOINT),
		APIKey:       os.Getenv(BUGSNAG_API_KEY),
		ReleaseStage: os.Getenv(ENVIRONMENT),
		Hostname:     hostname,
		Logger:	      log.New(ioutil.Discard, log.Prefix(), log.Flags()),
	})
	log.SetOutput(os.Stdout)

	certificate, err := tls.X509KeyPair(
		[]byte(os.Getenv(CERT_PEM)),
		[]byte(os.Getenv(CERT_KEY)),
	)
	if err != nil {
		Snag(fmt.Sprintf(
			"error creating TLS configuration for %s: %v",
			os.Getenv(FRONT_SERVICE),
			err,
		))
		return
	}
	config := tls.Config{
		Certificates: []tls.Certificate{certificate},
	}

	listener, err := tls.Listen("tcp", os.Getenv(FRONT_SERVICE), &config)
	if err != nil {
		Snag(fmt.Sprintf(
			"unable to start service %s: %v",
			os.Getenv(FRONT_SERVICE),
			err,
		))
		return
	}
	log.Println(
		"started listening for TLS connections on",
		os.Getenv(FRONT_SERVICE),
	)

	for {
		external, err := listener.Accept()
		if err != nil {
			Snag(fmt.Sprintf(
				"external facing service %s has gone down: %v",
				os.Getenv(FRONT_SERVICE),
				err,
			))
			return
		}
		go ProxyBack(external, os.Getenv(BACK_SERVICE))
	}
}
