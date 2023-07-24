package main

import (
	"game-server/internal/client"
	"game-server/internal/server"
	"log"
)

func main() {
	rgs := server.New()

	shutdownCh := server.HandleShutdownSignal(func() {
		rgs.Shutdown()
	})

	rgs.Loop()
	go client.Spawn(rgs)

	<-shutdownCh

	log.Println("[local-game-simulation] shutting down gracefully")
}
