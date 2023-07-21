package main

import (
	"game-server/internal/remoteserver"
	"log"
)

func main() {
	if err := remoteserver.Serve(); err != nil {
		log.Fatal(err)
	}
}
