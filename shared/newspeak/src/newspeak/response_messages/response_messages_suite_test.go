package response_messages_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestResponse_messages(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Response_messages Suite")
}
