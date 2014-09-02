package main

import (
	"github.com/CapillarySoftware/goforward/start"
	consul "github.com/armon/consul-api"
	log "github.com/cihub/seelog"
	"strconv"
	"time"
)

func main() {

	defer log.Flush()
	logger, err := log.LoggerFromConfigAsFile("seelog.xml")

	if err != nil {
		log.Warn("Failed to load config", err)
	}

	log.ReplaceLogger(logger)
	go RegisterService("goforward", 2025, 15)
	start.Run()
}

func RegisterService(name string, port int, ttl int) {

	reportInterval := make(chan bool, 1)
	go func() {
		for {
			time.Sleep(time.Duration(ttl) / 2 * time.Second)
			reportInterval <- true
		}
	}()

	client, err := consul.NewClient(consul.DefaultConfig())
	if nil != err {
		log.Error("Failed to get consul client")
	}

	for {
		select {
		case <-reportInterval: //report registration
			{

				agent := client.Agent()

				reg := &consul.AgentServiceRegistration{
					Name: "goforward",
					Port: port,
					Check: &consul.AgentServiceCheck{
						TTL: strconv.Itoa(ttl) + "s",
					},
				}
				if err := agent.ServiceRegister(reg); err != nil {
					log.Error("err: ", err)
				}
				checks, err := agent.Checks()
				if err != nil {
					log.Error("err: ", err)
				}
				if _, ok := checks["goforward"]; !ok {
					log.Error("Checks failed:, ", checks)
				}
			}
		}
	}

}
