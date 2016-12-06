package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGoReactor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoReactor Suite")
}
