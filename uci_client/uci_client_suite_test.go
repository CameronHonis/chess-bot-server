package uci_client_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUciClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Suite")
}
