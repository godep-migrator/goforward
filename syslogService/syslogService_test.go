package syslogService_test

import (
	. "github.com/CapillarySoftware/goforward/syslogService"

	// "fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SyslogService", func() {

	Describe("Invalid setup tests", func() {

		It("Invalid port test", func() {
			serv := SyslogService{ConType: TCP,
				RFCFormat: RFC3164,
				Port:      "99999999999"}
			err := serv.Bind()
			Expect(err).ShouldNot(Equal(BeNil()))
		})
	})

	Describe("Valid Tests", func() {
		It("Bind to valid port", func() {
			serv := SyslogService{ConType: TCP,
				RFCFormat: RFC3164,
				Port:      "9099"}
			err := serv.Bind()
			Expect(err).Should(BeNil())

		})
	})

})
