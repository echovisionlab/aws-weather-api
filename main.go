package main

import (
	"github.com/echovisionlab/aws-weather-api/pkg/app"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	service, err := app.New()

	if err != nil {
		log.Println(err.Error())
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	cancel := service.Run()

	<-exit
	cancel()
	log.Println("bye")
}
