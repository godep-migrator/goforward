package forward

import (
	"fmt"
	"github.com/CapillarySoftware/goforward/msgService"
)

func Run(channel <-chan msgService.ForwardMessage) {
	for {
		msgs := <-channel
		fmt.Println("msg: ", msgs)
	}
}
