package main

import (
	"game-server/internal/grpcclient"
	"log"
)

func main() {
	rc := grpcclient.New()

	if err := rc.Join(); err != nil {
		log.Fatal(err)
	}
}
