package main

import (
	"log"

	"github.com/hulizhen/blogo/config"
	"github.com/hulizhen/blogo/router"
	"github.com/hulizhen/blogo/store"
)

func main() {
	c, err := config.New()
	if err != nil {
		log.Panicf("Failed to parse the configurations with error: %v.\n", err)
	}

	s, err := store.New(c)
	if err != nil {
		log.Panicf("Failed to create store with error: %v.", err)
	}

	r := router.New(c, s)
	err = r.Run()
	if err != nil {
		log.Panicf("Failed to run the router with error: %v", err)
	}
}
