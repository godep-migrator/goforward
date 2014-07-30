package msgService

// Syslog Server with defined settings for RFC format and connection type.
import (
	log "github.com/cihub/seelog"
	"time"
)

type ForwardMessage interface {
	String() string
}

//SyslogServer interface
type Service interface {
	Bind() error
	SendMessages(chan ForwardMessage) error
}

//Main server thread.
func Run(server Service, msgsChan chan ForwardMessage) {
	for {
		err := server.Bind()
		if nil != err {
			log.Error("error Binding to service: ", err)
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	err := server.SendMessages(msgsChan)
	if nil != err {
		fmt.Println(err)
	}

}
