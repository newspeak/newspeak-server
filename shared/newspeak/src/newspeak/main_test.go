package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {
        Context("example context", func() {
                It("should pass", func() {
                        example := true
                        Expect(example).To(Equal(true))
                })
        })
})
