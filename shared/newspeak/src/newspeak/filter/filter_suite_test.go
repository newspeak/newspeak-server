package filter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "launchpad.net/gocheck"

	"testing"
)

// hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func TestFilter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filter Suite")
}
