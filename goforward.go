package main

import (
	"github.com/CapillarySoftware/goforward/start"
	consul "github.com/armon/consul-api"
	log "github.com/cihub/seelog"
)

func main() {
	defer log.Flush()
	logger, err := log.LoggerFromConfigAsFile("seelog.xml")

	if err != nil {
		log.Warn("Failed to load config", err)
	}

	client, err := consul.NewClient(consul.DefaultConfig())
	if nil != err {
		log.Error("Failed to get consul client")
	} else {

		agent := client.Agent()

		reg := &consul.AgentServiceRegistration{
			Name: "goforward",
			Check: &consul.AgentServiceCheck{
				TTL: "10s",
			},
		}
		if err := agent.ServiceRegister(reg); err != nil {
			log.Error("err: ", err)
		}
	}
	log.ReplaceLogger(logger)
	start.Run()
}
