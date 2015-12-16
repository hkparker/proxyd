package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"bytes"
	"log"
	"net"
	"os"
)

const FRONT_SERVICE = "TTPD_FRONT_SERVICE"
const BACK_SERVICE = "TTPD_BACK_SERVICE"
const CERT_PEM = "TTPD_CERT"
const CERT_KEY = "TTPD_KEY"
const ENVIRONMENT = "RACK_ENV"
const SLACK_ENDPOINT = "TTPD_SLACK_ENDPOINT"

var slack_rate_limit = make(map[string]bool)

func Snag(internal bool, service string, err error) {
	hostname, host_err := os.Hostname()
	if host_err != nil {
		hostname = ""
	}
	environment := os.Getenv(ENVIRONMENT)
	fallback_error := ""
	title := ""
	description := ""
	err_string := fmt.Sprintf("%v", err)

	if internal {
		fallback_error = fmt.Sprintf(
			"Service down on %s!  Could not contact internal service %s: %v",
			hostname,
			service,
			err,
		)
		title = "Internal service down!"
		description = fmt.Sprintf(
			"TLS connections are being accepted on %s, but the backend service %s is down.  This is likely a docker container on the host that has stopped responding.  More information about this host on the <https://aws.cbhq.net/|AWS Metadata Search>.",
			os.Getenv(FRONT_SERVICE),
			service,
		)
	} else {
		fallback_error = fmt.Sprintf(
			"Service down on %s!  External service %s has gone down: %v",
			hostname,
			service,
			err,
		)
		title = "External service down!"
		description = fmt.Sprintf(
			"The TLS listener on %s has died is no longer accepting connections for the backend service %s.  It is likely the TLS Terminator ran into an operating system limitation.  More information about this host on the <https://aws.cbhq.net/|AWS Metadata Search>.",
			service,
			os.Getenv(BACK_SERVICE),
		)
	}
	log.Println(fallback_error)

	client := http.Client{}
	flat := []byte(fmt.Sprintf(
		`{
			"username": "TLS Terminator",
			"icon_url": "http://rocketdock.com/images/screenshots/thumbnails/Terminator-Head.png",
			"attachments": [
				{
					"fallback":	"%s",
					"pretext":	"Service down on %s",
					"title":	"%s",
					"text":		"%s",
					"color":	"#FF0000",
					"fields": [
						{
							"title":	"Host",
							"value":	"%s",
							"short":	true
						},
						{
							"title":	"Service",
							"value":	"%s",
							"short":	true
						},
						{
							"title":	"Error",
							"value":	"%s",
							"short":	true
						},
						{
							"title":	"Environment",
							"value":	"%s",
							"short":	true
						}
					]
				}
			]
		}`,
		fallback_error,
		hostname,
		title,
		description,
		hostname,
		service,
		err_string,
		environment,
	))
	req, _ := http.NewRequest(
		"POST",
		os.Getenv(SLACK_ENDPOINT),
		bytes.NewBuffer(flat),
	)
	client.Do(req)
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
		Snag(
			true,
			os.Getenv(BACK_SERVICE),
			err,
		)
		external.Close()
		return
	}
	ExchangeData(external, internal)
}

func main() {
	log.SetOutput(os.Stdout)

	certificate, err := tls.X509KeyPair(
		[]byte(os.Getenv(CERT_PEM)),
		[]byte(os.Getenv(CERT_KEY)),
	)
	if err != nil {
		log.Println(fmt.Sprintf(
			"error creating TLS configuration for %s: %v",
			os.Getenv(FRONT_SERVICE),
			err,
		))
		return
	}
	config := tls.Config{
		Certificates: []tls.Certificate{certificate},
		// cipher list and curve preferences
	}

	listener, err := tls.Listen("tcp", os.Getenv(FRONT_SERVICE), &config)
	if err != nil {
		log.Println(fmt.Sprintf(
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
			Snag(
				false,
				os.Getenv(FRONT_SERVICE),
				err,
			)
			return
		}
		go ProxyBack(external, os.Getenv(BACK_SERVICE))
	}
}
