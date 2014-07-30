package forward

import (
	"github.com/CapillarySoftware/goforward/msgService"
	log "github.com/cihub/seelog"
)

func Run(channel <-chan msgService.ForwardMessage) {
	for msg := range channel {
		log.Trace(msg)
	}
}
