package main

import (
	"main/internal/client"
	"main/internal/server"
)

func main() {
	rgs := server.New()

	rgs.Loop()
	client.Loop(rgs)
}
