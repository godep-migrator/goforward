package start

import (
	"fmt"
	"github.com/CapillarySoftware/goforward/syslogServer"
)

//Main run loop for our package.
func Run() {
	fmt.Println("Starting goforward")
	_ = syslogServer
}
