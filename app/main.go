package main

import (
	"log"

	"github.com/estella-studio/leon-backend/internal/bootstrap"
)

func main() {
	log.Println("starting app")
	err := bootstrap.Start()
	if err != nil {
		log.Panic(err)
	}
}
