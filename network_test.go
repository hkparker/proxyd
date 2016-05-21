package main_test

import (
	"net"
	"testing"
)

func CreateConnection() (client, server net.Conn) {
	//socket := make(chan net.Conn, 1)
	//listener, err := net.Listen("tcp", "127.0.0.1:0")
	//Expect(err).To(BeNil())
	//go func() {
	//	conn, _ := listener.Accept()
	//	socket <- conn
	//}()
	//client, err = net.Dial("tcp", listener.Addr().String())
	//Expect(err).To(BeNil())
	//server = <-socket
	return
}

func TestExchangeDataCopiesBidirectionally(t *testing.T) {
	//client, server := CreateConnection()
	//defer client.Close()
	//defer server.Close()
	//go ExchangeData(server, client)
	//down_msg := "Going down"
	//up_msg := "Going up"
	//server.Write([]byte(down_msg))
	//client.Write([]byte(up_msg))
	//server_msg := make([]byte, 10)
	//n, err := client.Read(server_msg)
	//Expect(err).To(BeNil())
	//Expect(n).To(Equal(10))
	//Expect(string(server_msg)).To(Equal(down_msg))
	//client_msg := make([]byte, 8)
	//n, err = server.Read(client_msg)
	//Expect(err).To(BeNil())
	//Expect(n).To(Equal(8))
	//Expect(string(client_msg)).To(Equal(up_msg))
}

func TestExchangeDataClosesBothOnErr(t *testing.T) {
	//client, server := CreateConnection()
	//defer client.Close()
	//defer server.Close()
	//returned := make(chan bool, 1)
	//go func() {
	//	ExchangeData(server, client)
	//	returned <- true
	//}()
	//client.Close()
	//server.Write([]byte("closed"))
	//Expect(<-returned).To(Equal(true))
	//err := server.Close()
	//Expect(err).ToNot(BeNil())
}

func TestProxyBackDialsCorrectBackend(t *testing.T) {
	//backend := "127.0.0.1:3456"
	//listener, err := net.Listen("tcp", backend)
	//defer listener.Close()
	//Expect(err).To(BeNil())
	//connected := make(chan bool)
	//go func() {
	//	_, err := listener.Accept()
	//	Expect(err).To(BeNil())
	//	connected <- true
	//}()
	//client, _ := CreateConnection()
	//go ProxyBack(client, "tcp://"+backend, TLSConfig{})
	//Expect(<-connected).To(BeTrue())
}

func TestProxyBackClosesFrontendIfBackendIsClosed() {
	//client, _ := CreateConnection()
	//ProxyBack(client, "foo://", TLSConfig{})
	//err := client.Close()
	//Expect(err).ToNot(BeNil())
}
