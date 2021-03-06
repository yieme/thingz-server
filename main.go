package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/mchmarny/thingz-server/config"
	"github.com/mchmarny/thingz-server/load"
	"github.com/mchmarny/thingz-server/server"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
}

func main() {

	// make sure we can shutdown gracefully
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	doneDone := make(chan bool)
	errCh := make(chan error, 1)

	// REST
	go func() {
		router := server.NewRouter()
		address := fmt.Sprintf(":%d", config.Config.APIPort)
		log.Printf("Starting API Server - http://localhost%s", address)
		errCh <- http.ListenAndServe(address, router)
	}()

	// UI
	go func() {
		address := fmt.Sprintf(":%d", config.Config.UIPort)
		log.Printf("Starting UI Server - http://localhost%s", address)
		errCh <- http.ListenAndServe(address, http.FileServer(http.Dir("./ui")))
	}()

	// If configured, load Messages from Kafka
	if config.Config.Load {
		go load.LoadFromKafka()
	}

	go func() {

		for {
			select {
			case err := <-errCh:
				log.Printf("Error: %v", err)
			case ex := <-sigCh:
				log.Println("Shutting down due to: ", ex)
				doneDone <- true
			default:
				// nothing to do
			}
		}

	}()
	<-doneDone

}
