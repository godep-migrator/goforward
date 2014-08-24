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
	"strconv"
	"sync"
	"time"
)

//Define RFC syslog formats supported
type Format string

const (
	RFC3164 Format = "RFC3164"
	RFC5423 Format = "RFC5423"
)

//Define connection types supported.
type ConnectionType string

const (
	TCP ConnectionType = "tcp"
	UDP ConnectionType = "udp"
)

//Basic service struct.
type SyslogService struct {
	wg     *sync.WaitGroup
	format Format
	port   int
	cType  ConnectionType
	ln     *net.TCPListener
	done   chan bool
}

//Create a new syslogService
func NewSyslogService(cType ConnectionType, format Format, port int) (sys SyslogService, err error) {
	wg := sync.WaitGroup{}
	done := make(chan bool, 1)
	sys = SyslogService{wg: &wg, format: format, port: port, cType: cType, done: done}
	err = sys.Bind()
	return
}

//Start the tcp server
func (this *SyslogService) startTCP(msgsChan chan *messaging.Food) {
	this.wg.Add(1)
main:
	for {
		var (
			conn net.Conn
			err  error
		)
		select {
		case <-this.done:
			break main
		default:
			this.ln.SetDeadline(time.Now().Add(2 * time.Second))
			conn, err = this.ln.Accept()
			if err != nil {
				log.Debug(err)
				//usually read / write timeout
				continue
			}
			log.Trace("Accepted connection")
			this.wg.Add(1)
			go ProcessSyslog(conn, msgsChan, this.format, 5, this.done, this.wg)
		}
	}
	log.Debug("StartTCP exiting")
	this.wg.Done()
}

//Initialize syslog service
func (this *SyslogService) Start(msgsChan chan *messaging.Food) {
	switch this.cType {
	case TCP:
		{
			go this.startTCP(msgsChan)
		}
	case UDP:
		{

			log.Info("UDP Called")
		}
	}
	return
}

//Bind to syslog socket
func (this *SyslogService) Bind() (err error) {
	var ln *net.TCPListener
	var listen net.Listener
	switch this.cType {
	case TCP:
		{
			listen, err = net.Listen(string(this.cType), ":"+strconv.Itoa(this.port))
			if nil != err {
				return
			}
			ln = listen.(*net.TCPListener)
		}
	case UDP:
		{
			// var (
			// 	udpAddr *net.UDPAddr
			// )
			// udpAddr, err = net.ResolveUDPAddr("udp", ":"+this.port)
			// if err != nil {
			// 	return err
			// }
			// this.udpConn, err = net.ListenUDP(string(this.cType), udpAddr)
		}
	default:
		{
			log.Warn("Failed to provide valid connection type : ", this.cType)
		}

	}
	this.ln = ln
	return err
}

func (this *SyslogService) Close() {
	log.Info("Waiting for syslog connections to finish")
	close(this.done)
	this.wg.Wait()
	log.Info("SyslogService closed")
}

// //Get message from syslog socket
// func (this *SyslogService) SendMessages(msgsChan chan messaging.Food) (err error) {
// 	switch this.ConType {
// 	case TCP:
// 		{
// 			for {
// 				var conn net.Conn
// 				conn, err = this.ln.Accept()
// 				log.Trace("Accepted connection")
// 				if err != nil {
// 					return
// 				}
// 				this.wg.Add(1)
// 				go ProcessSyslog(conn, msgsChan, this.RFCFormat, 240, this.done, this.wg)
// 			}
// 		}

// 	case UDP:
// 		{
// 			this.wg.Add(1)
// 			go ProcessSyslog(this.udpConn, msgsChan, this.RFCFormat, 0, this.done, this.wg)
// 		}
// 	}
// 	return
// }

//Scan and parse messages
func ProcessSyslog(conn net.Conn, msgsChan chan *messaging.Food, format Format, timeout int, done <-chan bool, wg *sync.WaitGroup) {
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
					continue
				} else {
					msgsChan <- proto
				}

				if timeout > 0 {
					conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
				}
			} else {
				break main
			}
		case _, ok := <-done:
			if !ok {
				log.Debug("Closing Syslog connection because of shutdown")
				break main
			} else {
				log.Trace("Unknown message")
			}

		}
	}
	log.Info("Closing Syslog connection")
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
