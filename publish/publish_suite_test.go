package publish_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPublish(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Publish Suite")
}
