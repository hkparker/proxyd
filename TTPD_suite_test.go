package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTTPD(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TTPD Suite")
}
