package main

import (
	"github.com/lukas-w/consul-checks-logger/daemon"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := daemon.DefaultConfig()
	d, err := daemon.NewDaemon(config)

	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		d.Stop()
	}()

	d.Run()
}
