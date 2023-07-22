package main

import (
	"game-server/internal/grpcserver"
	"log"
)

func main() {
	if err := grpcserver.Serve(); err != nil {
		log.Fatal(err)
	}
}
