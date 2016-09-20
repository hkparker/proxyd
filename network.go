package main

import (
	"crypto/tls"
	"errors"
	log "github.com/Sirupsen/logrus"
	"io"
	"net"
)

func listenAny(uri string, tls_config tls.Config) (net.Listener, error) {
	if len(uri) < 7 {
		err_str := "uri too short"
		log.WithFields(log.Fields{
			"error": err_str,
			"uri":   uri,
		}).Error("cannot dial")
		return nil, errors.New(err_str)
	}

	if uri[:6] == "tls://" {
		address := uri[6:]
		return tls.Listen("tcp", address, &tls_config)
	} else if uri[:6] == "tcp://" {
		address := uri[6:]
		return net.Listen("tcp", address)
	} else if uri[:7] == "unix://" {
		address := uri[7:]
		return net.Listen("unix", address)
	}

	err_str := "unrecognized protocol"
	log.WithFields(log.Fields{
		"error": err_str,
		"uri":   uri,
	}).Error("cannot listen")
	return nil, errors.New(err_str)
}

func dialAny(uri string, tls_config tls.Config) (net.Conn, error) {
	if len(uri) < 7 {
		err_str := "uri too short"
		log.WithFields(log.Fields{
			"error": err_str,
		}).Error("cannot dial")
		return nil, errors.New(err_str)
	}

	if uri[:6] == "tls://" {
		address := uri[6:]
		return tls.Dial("tcp", address, &tls_config)
	} else if uri[:6] == "tcp://" {
		address := uri[6:]
		return net.Dial("tcp", address)
	} else if uri[:7] == "unix://" {
		address := uri[7:]
		return net.Dial("unix", address)
	}

	err_str := "unrecognized protocol"
	log.WithFields(log.Fields{
		"error": err_str,
		"uri":   uri,
	}).Error("cannot dial")
	return nil, errors.New(err_str)
}

func exchangeData(external, internal net.Conn) {
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
	<-errs
}

func proxyBack(external net.Conn, addr string, tls_config tls.Config) {
	internal, err := dialAny(addr, tls_config)
	if err != nil {
		external.Close()
		return
	}
	exchangeData(external, internal)
}
