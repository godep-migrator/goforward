package syslogService

import (
	// "errors"
	"bufio"
	// "github.com/jeromer/syslogparser"
	"github.com/jeromer/syslogparser/rfc3164"
	// 	"github.com/jeromer/syslogparser/rfc5424"
	"fmt"
	. "github.com/CapillarySoftware/goforward/msgService"
	"net"
	"time"
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
	udpConn   *net.UDPConn
}

//Bind to syslog socket
func (s *SyslogService) Bind() (err error) {
	switch s.ConType {
	case TCP:
		{
			s.ln, err = net.Listen(string(s.ConType), "localhost:"+s.Port)
		}
	case UDP:
		{
			var (
				udpAddr *net.UDPAddr
			)
			udpAddr, err = net.ResolveUDPAddr("udp", "127.0.0.1:"+s.Port)
			if err != nil {
				return err
			}
			s.udpConn, err = net.ListenUDP(string(s.ConType), udpAddr)
		}
	default:
		{
			fmt.Println("Failed to provide valid connection type : ", s.ConType)
		}

	}

	if err != nil {
		return
	}
	return
}

//Get message from syslog socket
func (s *SyslogService) SendMessages(msgsChan chan *ForwardMessage) (err error) {
	switch s.ConType {
	case TCP:
		{

			for {
				var conn net.Conn
				conn, err = s.ln.Accept()
				if err != nil {
					return
				}
				go ScanForMsgs(conn, msgsChan, s.RFCFormat, 240)
			}
		}

	case UDP:
		{
			go ScanForMsgs(s.udpConn, msgsChan, s.RFCFormat, 0)
		}
	}
	return
}

//Scan and parse messages
func ScanForMsgs(conn net.Conn, msgsChan chan *ForwardMessage, format Format, timeout int) {
	if timeout > 0 {
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	}

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		txt := scanner.Text()
		fmt.Println(txt)
		msg := ForwardMessage(rfc3164.NewParser([]byte(txt))) //TODO: Create interface for parsers and pass it to func
		msgsChan <- &msg

		if timeout > 0 {
			conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		}
	}
	conn.Close()

	return
}
