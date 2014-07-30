package forward

import (
	"fmt"
	"github.com/CapillarySoftware/goforward/msgService"
)

func Run(channel <-chan msgService.ForwardMessage) {
	for msgs := range channel {
		fmt.Println("msg: ", msgs)
	}
}
