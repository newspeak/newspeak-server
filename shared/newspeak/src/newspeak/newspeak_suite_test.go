package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "launchpad.net/gocheck"

	"testing"
)

// hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func TestNewspeak(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Newspeak Suite")
}
