package start

import (
	"fmt"
	"github.com/CapillarySoftware/goforward/forward"
	"github.com/CapillarySoftware/goforward/msgService"
	sys "github.com/CapillarySoftware/goforward/syslogService"
	"time"
)

//Main run loop for our package.
func Run() {
	fmt.Println("Starting goforward")

	msgForwardChan := make(chan *msgService.ForwardMessage, 1000)

	serv := sys.SyslogService{ConType: sys.TCP,
		RFCFormat: sys.RFC3164,
		Port:      "2024"}

	go msgService.Run(&serv, msgForwardChan)
	go forward.Run(msgForwardChan)
	for {
		time.Sleep(1000 * time.Millisecond)
	}
}
