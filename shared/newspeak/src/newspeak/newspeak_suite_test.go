package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "launchpad.net/gocheck"

	"testing"
)

// hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type NewspeakSuite struct{}
var _ = Suite(&NewspeakSuite{})

// unit tests
func TestNewspeak(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Newspeak Suite")
}

// benchmarks
func (s *NewspeakSuite) BenchmarkLogic(c *C) {
	for i := 0; i < c.N; i++ {
		// logic to benchmark
	}
}
