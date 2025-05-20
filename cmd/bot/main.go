package main

import (
	"github.com/dsrosen6/cw-ticket-bot/internal/api"
	"log"
)

func main() {
	s, err := api.NewServer()
	if err != nil {
		log.Fatal("error initializing server: ", err)
	}

	if err := s.Run(); err != nil {
		log.Fatal("error running server: ", err)
	}
}
