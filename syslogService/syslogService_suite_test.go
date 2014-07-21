package syslogService_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSyslogService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SyslogService Suite")
}
