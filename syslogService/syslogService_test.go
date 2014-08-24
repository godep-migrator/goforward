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
	var scanText MockScannerText
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
		scanText = *new(MockScannerText)
		scanText.TValue = "Bad input"

	})
	Describe("Invalid setup tests", func() {

		It("Invalid port test", func() {
			// cType ConnectionType, format Format, port int
			_, err := NewSyslogService(TCP, RFC3164, 99999999999)
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
			_, err := NewSyslogService(TCP, RFC3164, 9019)
			Expect(err).Should(BeNil())

		})

		It("RFC3164 to proto message", func() {
			msg, err := RFC3164ToProto(l)
			Expect(err).Should(BeNil())
			Expect(msg).ShouldNot(Equal(BeNil()))

		})

		It("Process RFC3164 test", func() {
			scanText.TValue = "<31>Jul 27 21:01:57 Bens-MacBook-Pro.local com.apple.metadata.mdflagwriter[287]: Done with /Users/vrecan/Library/Application Support/Google/Chrome/Local State"
			msg, err := ProcessRfc3164(&scanText)
			Expect(err).ShouldNot(Equal(BeNil()))
			Expect(msg).ShouldNot(Equal(BeNil()))
			rfc := msg.Rfc3164[0]
			Expect(rfc.GetTimestamp()).Should(Equal(int64(1406494917)))
			Expect(rfc.GetHostname()).Should(Equal("Bens-MacBook-Pro.local"))
			Expect(rfc.GetTag()).Should(Equal("com.apple.metadata.mdflagwriter"))
			Expect(rfc.GetContent()).Should(Equal("Done with /Users/vrecan/Library/Application Support/Google/Chrome/Local State"))
			Expect(rfc.GetPriority()).Should(Equal(int32(31)))
			Expect(rfc.GetFacility()).Should(Equal(int32(3)))
			Expect(rfc.GetSeverity()).Should(Equal(int32(7)))

		})
	})

})
