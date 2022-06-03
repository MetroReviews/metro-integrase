package lib

import (
	"metro-integrase/types"

	log "github.com/sirupsen/logrus"
)

/*
 * Package ``lib`` provides the main code for the integrase system.
 */

// Starts a web server handling all core integrase functions
func StartServer(adapter types.ListAdapter) {
	cfg := adapter.GetConfig()

	if cfg.StartupLogs {
		log.Info("Starting integrase server")
	}
}
