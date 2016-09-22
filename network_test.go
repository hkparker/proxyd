package main

import (
	"crypto/rand"
	"crypto/tls"
	b58 "github.com/jbenet/go-base58"
	"github.com/stretchr/testify/assert"
	"net"
	"os"
	"testing"
)

func CreateConnection(t *testing.T) (client, server net.Conn) {
	assert := assert.New(t)

	socket := make(chan net.Conn, 1)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	assert.Nil(err)
	go func() {
		conn, _ := listener.Accept()
		socket <- conn
	}()
	client, err = net.Dial("tcp", listener.Addr().String())
	assert.Nil(err)
	server = <-socket
	return
}

func TestListenAnyReportsURITooShort(t *testing.T) {
	assert := assert.New(t)

	_, err := listenAny("asdf", tls.Config{})
	if assert.NotNil(err) {
		assert.Equal("uri too short", err.Error())
	}
}

func TestListenAnyReportsUnrecognizesProtocols(t *testing.T) {
	assert := assert.New(t)

	_, err := listenAny("foobar://bizbaz", tls.Config{})
	if assert.NotNil(err) {
		assert.Equal("unrecognized protocol", err.Error())
	}
}

func TestListenAnyDialAnyTLS(t *testing.T) {

}

func TestListenAnyDialAnyTCP(t *testing.T) {
	assert := assert.New(t)

	results := make(chan error, 1)
	listener, err := listenAny("tcp://127.0.0.1:0", tls.Config{})
	assert.Nil(err)
	go func() {
		_, lerr := listener.Accept()
		results <- lerr
	}()
	_, err = dialAny("tcp://"+listener.Addr().String(), tls.Config{})
	assert.Nil(err)
	assert.Nil(<-results)
}

func TestListenAnyDialAnyUnix(t *testing.T) {
	assert := assert.New(t)

	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	assert.Nil(err)

	ipc_file := "/tmp/" + b58.Encode(bytes)
	ipc := "unix://" + ipc_file
	results := make(chan error, 1)
	listener, err := listenAny(ipc, tls.Config{})
	assert.Nil(err)
	go func() {
		_, lerr := listener.Accept()
		results <- lerr
	}()
	_, err = dialAny(ipc, tls.Config{})

	assert.Nil(<-results)
	err = os.Remove(ipc_file)
	assert.Nil(err)
}

func TestDialAnyReportsURITooShort(t *testing.T) {
	assert := assert.New(t)

	_, err := dialAny("asdf", tls.Config{})
	if assert.NotNil(err) {
		assert.Equal("uri too short", err.Error())
	}
}

func TestDialAnyReportsUnrecognizesProtocols(t *testing.T) {
	assert := assert.New(t)

	_, err := dialAny("foobar://bizbaz", tls.Config{})
	if assert.NotNil(err) {
		assert.Equal("unrecognized protocol", err.Error())
	}
}

func TestExchangeDataCopiesBidirectionally(t *testing.T) {
	assert := assert.New(t)

	client, server := CreateConnection(t)
	defer client.Close()
	defer server.Close()
	go exchangeData(server, client)
	down_msg := "Going down"
	up_msg := "Going up"
	server.Write([]byte(down_msg))
	client.Write([]byte(up_msg))
	server_msg := make([]byte, 10)
	n, err := client.Read(server_msg)
	assert.Nil(err)
	assert.Equal(10, n)
	assert.Equal(down_msg, string(server_msg))
	client_msg := make([]byte, 8)
	n, err = server.Read(client_msg)
	assert.Nil(err)
	assert.Equal(8, n)
	assert.Equal(up_msg, string(client_msg))
}

func TestExchangeDataClosesBothOnErr(t *testing.T) {
	assert := assert.New(t)

	client, server := CreateConnection(t)
	defer client.Close()
	defer server.Close()
	returned := make(chan bool, 1)
	go func() {
		exchangeData(server, client)
		returned <- true
	}()
	client.Close()
	server.Write([]byte("closed"))
	assert.Equal(true, <-returned)
	err := server.Close()
	assert.NotNil(err)
}

func TestProxyBackDialsCorrectBackend(t *testing.T) {
	assert := assert.New(t)

	backend := "127.0.0.1:3456"
	listener, err := net.Listen("tcp", backend)
	defer listener.Close()
	assert.Nil(err)
	connected := make(chan bool)
	go func() {
		_, err := listener.Accept()
		assert.Nil(err)
		connected <- true
	}()
	client, _ := CreateConnection(t)
	go proxyBack(client, "tcp://"+backend, tls.Config{})
	assert.Equal(true, <-connected)
}

func TestProxyBackClosesFrontendIfBackendIsClosed(t *testing.T) {
	assert := assert.New(t)

	client, _ := CreateConnection(t)
	proxyBack(client, "foo://foobar", tls.Config{})
	err := client.Close()
	assert.NotNil(err)
}
