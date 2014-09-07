package main

import (
	"github.com/CapillarySoftware/goforward/start"
	"github.com/CapillarySofware/goconsularis"
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
	goconsularis.RegisterService("goforward", 2025, 15)
	start.Run()
}
