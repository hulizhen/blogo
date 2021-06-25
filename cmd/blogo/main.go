package main

import (
	"blogo/config"
	"blogo/router"
	"blogo/store"
	"log"
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
