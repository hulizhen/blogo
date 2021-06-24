package main

import (
	"blogo/config"
	"blogo/router"
	"blogo/store"
	"log"
)

func main() {
	c, err := config.New(config.ConfigFilePath)
	if err != nil {
		log.Panicf("Failed to parse the configurations with error: %v.\n", err)
	}

	s, err := store.New(c)
	if err != nil {
		log.Panicf("Failed to create store with error: %v.", err)
	}

	r := router.New(c, s)
	r.Run()
}
