package syslogService

//Syslog service that receives syslog messages and forwards messages to perceptor

import (
	"bufio"
	"errors"
	. "github.com/jeromer/syslogparser"
	"github.com/jeromer/syslogparser/rfc3164"
	// 	"github.com/jeromer/syslogparser/rfc5424"
	"github.com/CapillarySoftware/goforward/messaging"
	log "github.com/cihub/seelog"
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
			s.ln, err = net.Listen(string(s.ConType), ":"+s.Port)
		}
	case UDP:
		{
			var (
				udpAddr *net.UDPAddr
			)
			udpAddr, err = net.ResolveUDPAddr("udp", ":"+s.Port)
			if err != nil {
				return err
			}
			s.udpConn, err = net.ListenUDP(string(s.ConType), udpAddr)
		}
	default:
		{
			log.Warn("Failed to provide valid connection type : ", s.ConType)
		}

	}

	if err != nil {
		return
	}
	return
}

//Get message from syslog socket
func (s *SyslogService) SendMessages(msgsChan chan messaging.Food) (err error) {
	switch s.ConType {
	case TCP:
		{
			for {
				var conn net.Conn
				conn, err = s.ln.Accept()
				log.Trace("Accepted connection")
				if err != nil {
					return
				}
				go SendMessagesFromSocket(conn, msgsChan, s.RFCFormat, 240)
			}
		}

	case UDP:
		{
			go SendMessagesFromSocket(s.udpConn, msgsChan, s.RFCFormat, 0)
		}
	}
	return
}

//Scan and parse messages
func SendMessagesFromSocket(conn net.Conn, msgsChan chan messaging.Food, format Format, timeout int) {
	if timeout > 0 {
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	}

	scanner := bufio.NewScanner(conn)

	var (
		proto *messaging.Food
		err   error
	)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		switch format {
		case RFC3164:
			{

				proto, err = ProcessRfc3164(&bytes)
			}
		case RFC5423:
			{
				errors.New("RFC5423 not implemented yet...")

			}
		}
		if nil != err {
			log.Error(err)
			continue
		} else {
			msgsChan <- *proto
		}

		if timeout > 0 {
			conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		}
	}
	err = scanner.Err()
	if nil != err {
		log.Error(err)
	}
	log.Debug("Closing connection")
	conn.Close()

	return
}

//simple interface to allow use to easily test
type ScannerText interface {
	Text() string
}

//Process rfc3164 message from bufio scanner and return a proto message
func ProcessRfc3164(bytes *[]byte) (food *messaging.Food, err error) {
	msg := rfc3164.NewParser(*bytes)
	msg.Parse()
	food, err = RFC3164ToProto(msg.Dump())
	return
}

//RFC3164 conversion to protobuffers
func RFC3164ToProto(lParts LogParts) (food *messaging.Food, err error) {
	proto := new(messaging.Rfc3164)
	pType := messaging.RFC3164
	for k, v := range lParts {
		switch k {
		case "timestamp":
			{

				time, ok := v.(time.Time)
				if !ok {
					errors.New("Invalid timestamp, not of rfc3164 type")
					return
				}
				unix := time.Unix()
				proto.Timestamp = &unix
			}
		case "hostname":
			{
				hostname, ok := v.(string)
				if !ok {
					errors.New("Invalid hostname, not of rfc3164 type")
					return
				}
				proto.Hostname = &hostname

			}
		case "tag":
			{
				tag, ok := v.(string)
				if !ok {
					errors.New("Invalid tag, not of rfc3164 type")
					return
				}
				proto.Tag = &tag
			}
		case "content":
			{
				content, ok := v.(string)
				if !ok {
					errors.New("Invalid content, not of rfc3164 type")
					return
				}
				proto.Content = &content

			}
		case "priority":
			{
				priority, ok := v.(int)
				if !ok {
					errors.New("Invalid priority, not of rfc3164 type")
					return
				}
				temp := int32(priority)
				proto.Priority = &temp

			}
		case "facility":
			{
				facility, ok := v.(int)
				if !ok {
					errors.New("Invalid facility, not of rfc3164 type")
					return
				}
				temp := int32(facility)
				proto.Facility = &temp

			}
		case "severity":
			{
				severity, ok := v.(int)
				if !ok {
					errors.New("Invalid severity, not of rfc3164 type")
					return
				}
				temp := int32(severity)
				proto.Severity = &temp

			}
		default:
			{
				errors.New("Invalid message, not of rfc3164 type")
			}
		}
	}
	food = new(messaging.Food)
	food.Type = &pType
	food.Rfc3164 = append(food.Rfc3164, proto)
	// food.Rfc3164 = proto
	return
}
