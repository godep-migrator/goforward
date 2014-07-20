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
	GetMsg() (ForwardMessage, error)
}

//Main server thread.
func Run(server Service) {
	for {
		err := server.Bind()
		if nil != err {
			fmt.Println("error Binding to service: ", err)
			time.Sleep(1000 * time.Millisecond)
		} else {
			break
		}
	}
	for {
		msg, err := server.GetMsg()
		if nil != err {
			fmt.Println(err)
		}
		fmt.Println("Msg: ", msg)

	}

}
