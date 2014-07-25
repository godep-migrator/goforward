package syslogService_test

import (
	. "github.com/CapillarySoftware/goforward/syslogService"
	. "github.com/jeromer/syslogparser"

	// "fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("SyslogService", func() {

	var l LogParts
	BeforeEach(func() {
		l = LogParts{
			"timestamp": time.Now(),
			"hostname":  "hostname",
			"tag":       "tag",
			"content":   "content",
			"priority":  1,
			"facility":  7,
			"severity":  2,
		}
	})
	Describe("Invalid setup tests", func() {

		It("Invalid port test", func() {
			serv := SyslogService{ConType: TCP,
				RFCFormat: RFC3164,
				Port:      "99999999999"}
			err := serv.Bind()
			Expect(err).ShouldNot(Equal(BeNil()))
		})

		It("RFC3164 extra field", func() {
			l["something"] = "something"
			msg, err := RFC3164ToProto(l)
			Expect(err).ShouldNot(Equal(BeNil()))
			Expect(msg).ShouldNot(Equal(BeNil()))

		})
		It("RFC3164 invalid timestamp", func() {
			l["timestamp"] = "not time.time"
			msg, err := RFC3164ToProto(l)
			Expect(err).ShouldNot(Equal(BeNil()))
			Expect(msg).ShouldNot(Equal(BeNil()))

		})
		It("RFC3164 invalid hostname", func() {
			l["hostname"] = 5
			msg, err := RFC3164ToProto(l)
			Expect(err).ShouldNot(Equal(BeNil()))
			Expect(msg).ShouldNot(Equal(BeNil()))

		})
		It("RFC3164 invalid tag", func() {
			l["tag"] = 5
			msg, err := RFC3164ToProto(l)
			Expect(err).ShouldNot(Equal(BeNil()))
			Expect(msg).ShouldNot(Equal(BeNil()))

		})
		It("RFC3164 invalid content", func() {
			l["content"] = 5
			msg, err := RFC3164ToProto(l)
			Expect(err).ShouldNot(Equal(BeNil()))
			Expect(msg).ShouldNot(Equal(BeNil()))

		})
		It("RFC3164 invalid priority", func() {
			l["priority"] = "not int"
			msg, err := RFC3164ToProto(l)
			Expect(err).ShouldNot(Equal(BeNil()))
			Expect(msg).ShouldNot(Equal(BeNil()))

		})
		It("RFC3164 invalid facility", func() {
			l["facility"] = "not int"
			msg, err := RFC3164ToProto(l)
			Expect(err).ShouldNot(Equal(BeNil()))
			Expect(msg).ShouldNot(Equal(BeNil()))

		})
		It("RFC3164 invalid severity", func() {
			l["severity"] = "not int"
			msg, err := RFC3164ToProto(l)
			Expect(err).ShouldNot(Equal(BeNil()))
			Expect(msg).ShouldNot(Equal(BeNil()))

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

		It("RFC3164 to proto message", func() {
			msg, err := RFC3164ToProto(l)
			Expect(err).Should(BeNil())
			Expect(msg).ShouldNot(Equal(BeNil()))

		})
	})

})
