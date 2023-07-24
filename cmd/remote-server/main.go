package main

import (
	"game-server/internal/remoteserver"
	"game-server/internal/server"
	"log"
)

func main() {
	rrgs := remoteserver.NewGrpcRandomGameServer()

	shutdownCh := server.HandleShutdownSignal(func() {
		rrgs.RandomGameServer.Shutdown()
	})

	go func() {
		if err := rrgs.Serve(); err != nil {
			log.Println(err)
		}
	}()

	<-shutdownCh

	log.Println("[remote-server] shutting down gracefully")
}
