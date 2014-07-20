package syslogServer

// Syslog Server with defined settings for RFC format and connection type.
import (
	"bufio"
	"errors"
	"github.com/jeromer/syslogparser"
	"github.com/jeromer/syslogparser/rfc3164"
	"github.com/jeromer/syslogparser/rfc5424"
	"net"
)

//Define RFC syslog formats supported
type Format int

const (
	RFC3164 Format = 1
	RFC5423 Format = 2
)

//Define connection types supported.
type ConnectionType int

const (
	TCP ConnectionType = 1
	UDP ConnectionType = 2
)

//SyslogServer interface
type SyslogService interface {
	getMsg() (msg SyslogMessage, err error)
}

//Basic service struct.
type Service struct {
	ConType   ConnectionType
	RFCFormat Format

	// scanners    []*bufio.Scanner
	// listeners   []*net.TCPListener
	// connections []net.Conn
	// format      Format
	// handler     Handler
	// lastError   error
}

type SyslogMessage interface {
}

//Main server thread.
func Run(server *SyslogService) {
	for {
		msg, err := server.getMsg()
		if nil != err {
			fmt.Prinln(err)
		}
		fmt.Println("Got Mesg: ", msg)
	}

}
