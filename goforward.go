package main

import (
	"github.com/CapillarySoftware/goconsularis"
	"github.com/CapillarySoftware/goforward/start"
	log "github.com/cihub/seelog"
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
