package main_test

import (
	"net"
	"log"
	"io/ioutil"
	. "github.com/hkparker/TTPD"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	//"github.com/onsi/gomega/gbytes"
	//"os"
	//"os/exec"
	//"fmt"
	//"crypto/tls"
)

func CreateConnection(service string) (client, server net.Conn) {
	socket := make(chan net.Conn, 1)
	listener, err := net.Listen("tcp", service)
	Expect(err).To(BeNil())
	go func() {
		conn, _ := listener.Accept()
		socket <- conn
	}()
	client, err = net.Dial("tcp", service)
	Expect(err).To(BeNil())
	server = <-socket
	return client, server
}

var _ = Describe("TTPD", func() {
	log.SetOutput(ioutil.Discard)

	var binary_path string
	BeforeSuite(func() {
		var err error
		binary_path, err = gexec.Build("github.com/hkparker/TTPD")
		Expect(err).To(BeNil())
	})

	AfterSuite(func() {
	    gexec.CleanupBuildArtifacts()
	})

	Describe("ExchangeData", func() {

		It("copies data in both directions", func() {
			client, server := CreateConnection("127.0.0.1:1234")
			defer client.Close()
			defer server.Close()
			go ExchangeData(server, client)
			down_msg := "Going down"
			up_msg := "Going up"
			server.Write([]byte(down_msg))
			client.Write([]byte(up_msg))
			server_msg := make([]byte, 10)
			n, err := client.Read(server_msg)
			Expect(err).To(BeNil())
			Expect(n).To(Equal(10))
			Expect(string(server_msg)).To(Equal(down_msg))
			client_msg := make([]byte, 8)
			n, err = server.Read(client_msg)
			Expect(err).To(BeNil())
			Expect(n).To(Equal(8))
			Expect(string(client_msg)).To(Equal(up_msg))
		})

		It("closes both sockets and returns when one closes", func() {
			client, server := CreateConnection("127.0.0.1:2345")
			defer client.Close()
			defer server.Close()
			returned := make(chan bool, 1)
			go func() {
				ExchangeData(server, client)
				returned <- true
			}()
			client.Close()
			server.Write([]byte("closed"))
			Expect(<-returned).To(Equal(true))
			err := server.Close()
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("ProxyBack", func() {
		It("dials the backend service when frontend connection received", func() {
			backend := "127.0.0.1:3456"
			listener, err := net.Listen("tcp", backend)
			defer listener.Close()
			Expect(err).To(BeNil())
			connected := make(chan bool)
			go func() {
				_, err := listener.Accept()
				Expect(err).To(BeNil())
				connected <- true
			}()
			client, _ := CreateConnection("127.0.0.1:3457")
			go ProxyBack(client, "tcp://"+backend, TLSConfig{})
			Expect(<-connected).To(BeTrue())
		})

		It("closes the frontend connection when it can't dial the backend service", func() {
			client, _ := CreateConnection("127.0.0.1:4567")
			ProxyBack(client, "tcp://127.0.0.1:4568", TLSConfig{})
			err := client.Close()
			Expect(err).ToNot(BeNil())
		})
	})


})
