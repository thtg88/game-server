package main

import (
	"game-server/internal/remoteserver"
	"log"
)

func main() {
	rrgs := remoteserver.NewGrpcRandomGameServer()

	if err := rrgs.Serve(); err != nil {
		log.Fatal(err)
	}
}
