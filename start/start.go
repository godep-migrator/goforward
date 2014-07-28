package start

import (
	"flag"
	"fmt"
	"github.com/CapillarySoftware/goforward/forward"
	"github.com/CapillarySoftware/goforward/msgService"
	sys "github.com/CapillarySoftware/goforward/syslogService"
	"strconv"
	"strings"
	"time"
)

//Main run loop for our package.

var port = flag.Int("port", 514, "Syslog port you are going to listen on.")
var protocol = flag.String("protocol", "udp", "Syslog protocol options (udp,tcp)")

func processProtocol(proto string) (protocol sys.ConnectionType, err error) {
	protocol = sys.ConnectionType(strings.ToLower(proto))
	return
}

func Run() {
	flag.Parse()
	fmt.Println("Starting goforward")

	msgForwardChan := make(chan msgService.ForwardMessage, 1000)

	serv := sys.SyslogService{ConType: sys.UDP,
		RFCFormat: sys.RFC3164,
		Port:      strconv.Itoa(*port)}

	go msgService.Run(&serv, msgForwardChan)
	go forward.Run(msgForwardChan)
	for {
		time.Sleep(1000 * time.Millisecond)
	}
}
