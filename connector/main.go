package main

import (
	"coding.net/miraclew/pusher/xlog"
	logp "log"
	"os"
)

var (
	log *xlog.Logger
)

func main() {
	f, err := os.OpenFile("connector.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logp.Fatalf("error opening file: %v", err)
		return
	}
	defer f.Close()

	log = xlog.NewNamed(f, "aa")
	log.Debug = true

	log.Debugf("debugf %s", "ddd")
	log.Warnf("warn %s", "ddd")
	log.Infof("info %s", "ddd")
	log.Errorf("error %s", "ddd")
}
