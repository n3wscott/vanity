package main

import (
	"github.com/n3wscott/vanity/pkg/runner"
	"log"
)

func main() {
	log.Print("here")
	// Blocking call.
	runner.Run()
}
