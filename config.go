package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Config holds configuration for this application
type Config struct {
	TWHOST string `required:"true"`
}

var conf = &Config{}

func init() {
	err := envconfig.Process("", conf)
	if err != nil {
		log.Fatalf("Couldn't set configuration: %s", err)
	}
}
