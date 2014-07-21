package syslogService

import (
	// "errors"
	// 	"github.com/jeromer/syslogparser"
	"bufio"
	"github.com/jeromer/syslogparser/rfc3164"
	// 	"github.com/jeromer/syslogparser/rfc5424"
	"fmt"
	. "github.com/CapillarySoftware/goforward/msgService"
	"net"
)

//Define RFC syslog formats supported
type Format int

const (
	RFC3164 Format = 1
	RFC5423 Format = 2
)

//Define connection types supported.
type ConnectionType string

const (
	TCP ConnectionType = "tcp"
	UDP ConnectionType = "udp"
)

//Basic service struct.
type SyslogService struct {
	ConType   ConnectionType
	RFCFormat Format
	Port      string
	ln        net.Listener
	conn      net.Conn
}

//Bind to syslog socket
func (s *SyslogService) Bind() (err error) {
	s.ln, err = net.Listen(string(s.ConType), ":"+s.Port)
	if err != nil {
		return
	}
	return
}

//Get message from syslog socket
func (s *SyslogService) SendMessages(msgsChan chan *[]ForwardMessage) (err error) {
	for {
		s.conn, err = s.ln.Accept()
		if err != nil {
			return
		}
		go s.ScanForMsgs(s.conn)
	}
	return
}

//Scan and parse messages
func (s *SyslogService) ScanForMsgs(conn net.Conn) (msgs *[]ForwardMessage, err error) {

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		parser := rfc3164.NewParser([]byte(scanner.Text()))
		fmt.Println("Parser: ", parser)
	}
	conn.Close()

	return
}
