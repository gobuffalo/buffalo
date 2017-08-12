package snaker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Snaker Suite")
}
