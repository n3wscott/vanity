package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/n3wscott/vanity/pkg/handler"
	"github.com/n3wscott/vanity/pkg/vanity"
)

type envConfig struct {
	Port         int    `envconfig:"PORT" default:"8080"`
	VanityConfig string `envconfig:"VANITY_CONFIG" required:"true"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}

	cfg, err := vanity.Parse(env.VanityConfig)
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", env.Port), handler.New(cfg)); err != nil {
		log.Fatal(err)
	}
}
