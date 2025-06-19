package main

import (
	"os"
	"os/signal"
	"syscall"
	"test/cron"
	"test/handler"
	"test/httpserver"
)

func main() {
	httpServer := httpserver.New()
	go httpServer.Run(handler.GinHandler())

	cronServer := cron.New()
	go cronServer.Run()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sc)

	<-sc
}
