package main

import (
	"crypto/tls"
	"github.com/stretchr/testify/assert"
	"net"
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
