package msgService

// Syslog Server with defined settings for RFC format and connection type.
import (
	"github.com/CapillarySoftware/goforward/messaging"
	log "github.com/cihub/seelog"
	"sync"
	"time"
)

//SyslogServer interface
type Service interface {
	Bind() error
	SendMessages(chan messaging.Food) error
	Close()
}

//Main server thread.
func Run(server Service, msgsChan chan messaging.Food, wg *sync.WaitGroup) {
	defer wg.Done()
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
		log.Error(err)
	}

}
