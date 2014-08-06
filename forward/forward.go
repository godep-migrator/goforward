package forward

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/CapillarySoftware/goforward/messaging"
	log "github.com/cihub/seelog"
	nano "github.com/op/go-nanomsg"
)

const ()

func Run(channel <-chan messaging.Food) {
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
	for msg := range channel {
		index := "document"
		indexType := "all"
		id := string(uuid.NewRandom())
		msg.Index = &index
		msg.IndexType = &indexType
		msg.Id = &id
		log.Trace(msg)
		bytes, err := msg.Marshal()
		if nil != err {
			log.Error(err)
			continue
		}
		_, err = socket.Send(bytes, nano.DontWait) //blocking
		if nil != err {
			log.Error(err)
		}
	}
}
