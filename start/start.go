package start

//Start manages the main run loop of the application
import (
	"flag"
	"fmt"
	"github.com/CapillarySoftware/goforward/forward"
	"github.com/CapillarySoftware/goforward/msgService"
	sys "github.com/CapillarySoftware/goforward/syslogService"
	"os"
	"os/signal"
	// "reflect"
	"strconv"
	"strings"
	// "time "
)

var port = flag.Int("port", 514, "Syslog port you are going to listen on.")
var protocol = flag.String("protocol", "udp", "Syslog protocol options (udp,tcp)")

//Process protocol from input flags
func ProcessProtocol(proto string) (protocol sys.ConnectionType) {
	protocol = sys.ConnectionType(strings.ToLower(proto))
	return
}

//Manage death of application by signal
func Death(c <-chan os.Signal, death chan int) {
	for sig := range c {
		switch sig.String() {
		case "terminated":
			{
				death <- 1
			}
		case "interrupt":
			{
				death <- 2
			}
		default:
			{
				death <- 3
			}
		}

	}
}

//Run the app.
func Run() {
	flag.Parse()
	fmt.Println("Starting goforward")
	proto := ProcessProtocol(*protocol)

	msgForwardChan := make(chan msgService.ForwardMessage, 1000)

	serv := sys.SyslogService{ConType: proto,
		RFCFormat: sys.RFC3164,
		Port:      strconv.Itoa(*port)}

	go msgService.Run(&serv, msgForwardChan)
	go forward.Run(msgForwardChan)
	c := make(chan os.Signal, 1)
	s := make(chan int, 1)
	signal.Notify(c)
	go Death(c, s)
	death := <-s //time for shutdown
	close(msgForwardChan)
	fmt.Println(death)
	fmt.Println("Exiting")
}
