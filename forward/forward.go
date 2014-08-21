package forward

import (
	"github.com/CapillarySoftware/goforward/messaging"
	log "github.com/cihub/seelog"
	nano "github.com/op/go-nanomsg"
	"sync"
	"time"
)

const ()

func Run(channel <-chan messaging.Food, wg *sync.WaitGroup) {
	defer wg.Done()
	socket, err := nano.NewPushSocket()
	if nil != err {
		log.Error(err)
	}
	defer socket.Close()
	_, err = socket.Connect("tcp://localhost:2025")
	if nil != err {
		log.Error(err)
		return
	}
	socket.SetSendTimeout(1 * time.Minute)
	for msg := range channel {
		//add time here
		log.Trace(msg)
		bytes, err := msg.Marshal()
		if nil != err {
			log.Error(err)
			continue
		}
		_, err = socket.Send(bytes, 0) //blocking until timeout hit
		if nil != err {
			log.Error(err)
		}
	}
}
