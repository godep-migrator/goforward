package msgService_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMsgService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MsgService Suite")
}
