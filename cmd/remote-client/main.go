package main

import (
	"game-server/internal/remoteclient"
	"log"
)

func main() {
	rc := remoteclient.New()

	if err := rc.Join(); err != nil {
		log.Fatal(err)
	}
}
