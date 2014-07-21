package msgService

// Syslog Server with defined settings for RFC format and connection type.
import (
	"fmt"
	"time"
)

type ForwardMessage interface {
}

//SyslogServer interface
type Service interface {
	Bind() error
	SendMessages(chan *[]ForwardMessage) error
}

//Main server thread.
func Run(server Service, msgsChan chan *[]ForwardMessage) {
	for {
		err := server.Bind()
		if nil != err {
			fmt.Println("error Binding to service: ", err)
			time.Sleep(1000 * time.Millisecond)
		} else {
			break
		}
	}
	err := server.SendMessages(msgsChan)
	if nil != err {
		fmt.Println(err)
	}

}
