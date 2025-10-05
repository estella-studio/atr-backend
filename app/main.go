package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/estella-studio/atr-backend/internal/bootstrap"
)

func main() {
	log.Println("starting app")
	app, port, err := bootstrap.Start()
	if err != nil {
		log.Panic(err)
	}

	go func() {
		err := app.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Println("gracefully shutting down")

	err = app.Shutdown()
	if err != nil {
		log.Println("graceful shutdown failed")
	}

	log.Println("shutdown complete")
}
