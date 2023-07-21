package main

import (
	"game-server/internal/client"
	"game-server/internal/server"
)

func main() {
	rgs := server.New()

	rgs.Loop()
	client.Loop(rgs)
}
