package start

import (
	"fmt"
	"github.com/CapillarySoftware/goforward/msgService"
	. "github.com/CapillarySoftware/goforward/syslogService"
	"time"
)

//Main run loop for our package.
func Run() {
	fmt.Println("Starting goforward")
	serv := SyslogService{ConType: TCP,
		RFCFormat: RFC3164,
		Port:      "514"}

	go msgService.Run(&serv)

	for {
		time.Sleep(1000 * time.Millisecond)
	}
}
