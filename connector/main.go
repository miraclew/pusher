package main

import (
	log2 "coding.net/miraclew/pusher/log"
	logp "log"
	"os"
)

var (
	log *log2.Logger
)

func main() {
	f, err := os.OpenFile("connector.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logp.Fatalf("error opening file: %v", err)
		return
	}
	defer f.Close()

	log = log2.NewNamed(f, "aa")
	log.Debug = true

	log.Debugf("debugf %s", "ddd")
	log.Warnf("warn %s", "ddd")
	log.Infof("info %s", "ddd")
	log.Errorf("error %s", "ddd")
}
