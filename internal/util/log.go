package util

import (
	"github.com/crazygit/hpv-notification/config"
	log "github.com/sirupsen/logrus"
	"os"
)

func ConfigLog() {
	if config.AppConfig.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.InfoLevel)
	}
	log.SetOutput(os.Stdout)
}
