package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type envConfig struct {
	Port          int    `envconfig:"PORT" default:"8080"`
	Downstream    string `envconfig:"DISCOVERY_DOWNSTREAM"` // comma separated list of urls.
	Services      string `envconfig:"DISCOVERY_SERVICES_FILE"`
	Subscriptions string `envconfig:"SUBSCRIPTIONS_FILE"`
}

func main() {

	var configPath string
	switch len(os.Args) {
	case 1:
		configPath = "vanity.yaml"
	case 2:
		configPath = os.Args[1]
	default:
		log.Fatal("usage: govanityurls [CONFIG]")
	}
	vanity, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	h, err := newHandler(vanity)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", h)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
