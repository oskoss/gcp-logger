package consume_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConsume(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Consume Suite")
}
