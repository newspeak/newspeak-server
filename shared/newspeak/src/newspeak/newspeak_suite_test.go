package newspeak_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNewspeak(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Newspeak Suite")
}
