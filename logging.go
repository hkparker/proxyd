package main

import (
	"log"
	"net/http"
	"encoding/json"
	"os"
	"bytes"
	"fmt"
	"time"
)

const ENVIRONMENT = "ENVIRONMENT"
const SLACK_ENDPOINT = "TTPD_SLACK_ENDPOINT"

var hostname = "unknown"
var environment = "unknown"
var slack_endpoint = ""
var service_up map[string]bool
var listener_recreated map[string]string

func AssignLoggingVariables() {
	service_up = make(map[string]bool)
	listener_recreated = make(map[string]string)
	hostname, _ = os.Hostname()
	if env := os.Getenv(ENVIRONMENT); env != "" {
		environment = env
	} else if env = os.Getenv("RACK_ENV"); env != "" {
		environment = env
	} else if env = os.Getenv("NODE_ENV"); env != "" {
		environment = env
	}
	slack_endpoint = os.Getenv(SLACK_ENDPOINT)
}

func SlackPost(chat_message []byte) {
	if slack_endpoint != "" {
		client := http.Client{}
		req, _ := http.NewRequest(
			"POST",
			slack_endpoint,
			bytes.NewBuffer(chat_message),
		)
		client.Do(req)
	}
}

type ServicesStartedLog struct {
	Now		string
	Hostname	string
	Event		string
	Environment	string
}

func LogServicesStarted() {
	msg, _ := json.Marshal(ServicesStartedLog{
		Now:		time.Now().String(),
		Hostname:	hostname,
		Event:		"services_started",
		Environment:	environment,
	})
	log.Println(string(msg))
}

type ServiceDownLog struct {
	Now		string
	Service		string
	Hostname	string
	Environment	string
	Event		string
	Error		error
}

func LogServiceDown(service string, err error) {
	msg, _ := json.Marshal(ServiceDownLog{
		Now:		time.Now().String(),
		Service:	service,
		Hostname:	hostname,
		Environment:	environment,
		Event:		"service_down",
		Error:		err,
	})
	log.Println(string(msg))

	err_string := fmt.Sprintf("%v", err)
	fallback_error := fmt.Sprintf(
		"Service down on %s!  Could not contact internal service %s: %v",
		hostname,
		service,
		err,
	)
	description := fmt.Sprintf(
		"Internal service %s is down.  This is likely a docker container on the host that has stopped responding.  More information about this host on the <https://aws.cbhq.net/#%s|AWS Metadata Search>.",
		service,
		hostname,
	)
	SlackPost([]byte(fmt.Sprintf(
		`{
			"username": "TLS Terminator",
			"icon_url": "http://i.imgur.com/64l3NXn.png",
			"attachments": [
				{
					"fallback":	"%s",
					"pretext":	"Service down on %s",
					"title":	"Internal service down!",
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
		description,
		hostname,
		service,
		err_string,
		environment,
	)))
}

type ServiceRecoveredLog struct {
	Now		string
	Service		string
	Hostname	string
	Environment	string
	Event		string
}

func LogServiceRecovered(service string) {
	msg, _ := json.Marshal(ServiceRecoveredLog{
		Now:		time.Now().String(),
		Service:	service,
		Hostname:	hostname,
		Environment:	environment,
		Event:		"service_recovered",
	})
	log.Println(string(msg))

	fallback_message := fmt.Sprintf(
		"Service recovered on %s: %s",
		hostname,
		service,
	)
	description := fmt.Sprintf(
		"Internal service %s has recovered.  More information about this host on the <https://aws.cbhq.net/#%s|AWS Metadata Search>.",
		service,
		hostname,
	)
	SlackPost([]byte(fmt.Sprintf(
		`{
			"username": "TLS Terminator",
			"icon_url": "http://i.imgur.com/64l3NXn.png",
			"attachments": [
				{
					"fallback":	"%s",
					"pretext":	"Service recovered on %s",
					"title":	"Internal service recovered!",
					"text":		"%s",
					"color":	"#00BB00",
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
							"title":	"Environment",
							"value":	"%s",
							"short":	true
						}
					]
				}
			]
		}`,
		fallback_message,
		hostname,
		description,
		hostname,
		service,
		environment,
	)))
}

type ServiceParseFailedLog struct {
	Now		string
	Hostname	string
	Event		string
	Environment	string
}

func LogServiceParseFailed() {
	msg, _ := json.Marshal(ServiceParseFailedLog{
		Now:		time.Now().String(),
		Hostname:	hostname,
		Event:		"service_parse_failed",
		Environment:	environment,
	})
	log.Println(string(msg))

	fallback_error := fmt.Sprintf(
		"Could not parse services for %s",
		hostname,
	)
	description := fmt.Sprintf(
		"Unable to parse service configuration for %s.  No services have been started.  More information about this host on the <https://aws.cbhq.net/#%s|AWS Metadata Search>.",
		hostname,
		hostname,
	)
	SlackPost([]byte(fmt.Sprintf(
		`{
			"username": "TLS Terminator",
			"icon_url": "http://i.imgur.com/64l3NXn.png",
			"attachments": [
				{
					"fallback":	"%s",
					"pretext":	"Service parse failed on %s",
					"title":	"Services not parsed!",
					"text":		"%s",
					"color":	"#FF0000",
					"fields": [
						{
							"title":	"Host",
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
		description,
		hostname,
		environment,
	)))
}
