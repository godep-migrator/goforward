package start

//Start manages the main run loop of the application
import (
	"flag"
	"github.com/CapillarySoftware/goforward/forward"
	"github.com/CapillarySoftware/goforward/messaging"
	"github.com/CapillarySoftware/goforward/msgService"
	sys "github.com/CapillarySoftware/goforward/syslogService"
	log "github.com/cihub/seelog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
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
	log.Info("Starting goforward")
	flag.Parse()
	var wg sync.WaitGroup

	proto := ProcessProtocol(*protocol)

	msgForwardChan := make(chan messaging.Food, 1000)
	serv := sys.NewSyslogService(proto, sys.RFC3164, strconv.Itoa(*port))

	wg.Add(2)
	go msgService.Run(&serv, msgForwardChan, &wg)
	go forward.Run(msgForwardChan, &wg)

	c := make(chan os.Signal, 1)
	s := make(chan int, 1)
	signal.Notify(c)
	go Death(c, s)
	death := <-s //time for shutdown
	log.Info("Closing syslog server")
	serv.Close()
	//close only after all senders are done
	close(msgForwardChan)
	log.Info("Waiting for everything to come down gracefully...")
	wg.Wait()
	log.Debug("Death return code: ", death)
	log.Info("Exiting")
}
