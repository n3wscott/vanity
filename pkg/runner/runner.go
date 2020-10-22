package runner

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/n3wscott/vanity/pkg/handler"
	"github.com/n3wscott/vanity/pkg/vanity"
	"log"
	"net/http"
)

func Run() {
	var env struct {
		Port         int    `envconfig:"PORT" default:"8080"`
		VanityConfig string `envconfig:"VANITY_CONFIG" default:"kodata/example.yaml" required:"true"`
	}

	if err := envconfig.Process("", &env); err != nil {
		log.Fatal("Failed to process env var: ", err)
	}

	cfg, err := vanity.Parse(env.VanityConfig)
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", env.Port), handler.New(cfg)); err != nil {
		log.Fatal(err)
	}
}
