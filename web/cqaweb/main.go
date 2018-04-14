package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/donyori/cqa/web"
)

func main() {
	exitC := make(chan os.Signal, 1)
	signal.Notify(exitC, os.Interrupt, os.Kill)

	defer func() {
		log.Println("cqaserver exits.")
	}()
	defer func() {
		web.CleanUp()
		log.Println("Web cleanup finish.")
	}()

	web.Init()
	log.Println("Web init finish.")
	serverC := web.LaunchBackground()
	log.Println("Web launched.")
	select {
	case err := <-serverC:
		log.Println("Server error:", err)
	case sig := <-exitC:
		log.Println("System signal", sig, "trapped.")
		err := web.Shutdown()
		log.Println("Shutdown error:", err)
		err = <-serverC
		if err != nil && err != http.ErrServerClosed {
			log.Println("Server error:", err)
		} else {
			log.Println("Server closed.")
		}
	}
}
