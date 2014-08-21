package syslogService

//Syslog service that receives syslog messages and forwards messages to perceptor

import (
	"bufio"
	"errors"
	. "github.com/jeromer/syslogparser"
	"github.com/jeromer/syslogparser/rfc3164"
	// 	"github.com/jeromer/syslogparser/rfc5424"
	"code.google.com/p/go-uuid/uuid"
	"github.com/CapillarySoftware/goforward/messaging"
	log "github.com/cihub/seelog"
	"net"
	"sync"
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
	done      chan bool
	wg        *sync.WaitGroup
}

func NewSyslogService(cType ConnectionType, format Format, port string) (serv SyslogService) {
	done := make(chan bool, 1)
	wg := sync.WaitGroup{}
	serv = SyslogService{ConType: cType, RFCFormat: format, Port: port, done: done, wg: &wg}
	return
}

//Bind to syslog socket
func (this *SyslogService) Bind() (err error) {
	switch this.ConType {
	case TCP:
		{
			this.ln, err = net.Listen(string(this.ConType), ":"+this.Port)
		}
	case UDP:
		{
			var (
				udpAddr *net.UDPAddr
			)
			udpAddr, err = net.ResolveUDPAddr("udp", ":"+this.Port)
			if err != nil {
				return err
			}
			this.udpConn, err = net.ListenUDP(string(this.ConType), udpAddr)
		}
	default:
		{
			log.Warn("Failed to provide valid connection type : ", this.ConType)
		}

	}

	if err != nil {
		return
	}
	return
}

func (this *SyslogService) Close() {
	close(this.done)
	this.wg.Wait()
}

//Get message from syslog socket
func (this *SyslogService) SendMessages(msgsChan chan messaging.Food) (err error) {
	switch this.ConType {
	case TCP:
		{
			for {
				var conn net.Conn
				conn, err = this.ln.Accept()
				log.Trace("Accepted connection")
				if err != nil {
					return
				}
				this.wg.Add(1)
				go SendMessagesFromSocket(conn, msgsChan, this.RFCFormat, 240, this.done, this.wg)
			}
		}

	case UDP:
		{
			this.wg.Add(1)
			go SendMessagesFromSocket(this.udpConn, msgsChan, this.RFCFormat, 0, this.done, this.wg)
		}
	}
	return
}

//Scan and parse messages
func SendMessagesFromSocket(conn net.Conn, msgsChan chan messaging.Food, format Format, timeout int, done <-chan bool, wg *sync.WaitGroup) {
	if timeout > 0 {
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	}
	defer wg.Done()
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	var (
		proto *messaging.Food
		err   error
	)
main:
	for {
		select {
		default:
			if scanner.Scan() {
				switch format {
				case RFC3164:
					{
						proto, err = ProcessRfc3164(scanner)
					}
				case RFC5423:
					{
						errors.New("RFC5423 not implemented yet...")

					}
				}
				if nil != err {
					log.Error(err)
				} else {
					msgsChan <- *proto
				}

				if timeout > 0 {
					conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
				}
			} else {
				log.Info("Closing connection")
				break main
			}
		case _, ok := <-done:
			if !ok {
				log.Debug("Closing connection because of shutdown")
				break main
			} else {
				log.Trace("Unknown message")
			}

		}
	}

	return
}

//simple interface to allow use to easily test
type ScannerText interface {
	Text() string
}

//Process rfc3164 message from bufio scanner and return a proto message
func ProcessRfc3164(scanner ScannerText) (food *messaging.Food, err error) {
	msg := rfc3164.NewParser([]byte(scanner.Text()))
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
	id := uuid.NewRandom().String()
	proto.Id = &id
	food = new(messaging.Food)
	food.Type = &pType
	ts := time.Now().UTC().UnixNano()
	food.TimeNano = &ts
	food.Rfc3164 = append(food.Rfc3164, proto)
	// food.Rfc3164 = proto
	return
}
