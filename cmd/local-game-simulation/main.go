package main

import (
	"game-server/internal/client"
	"game-server/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	rgs := server.New()

	rgs.Loop()
	go client.Spawn(rgs)

	shutdownCh := handleShutdownSignal(func() {
		rgs.Shutdown()
	})

	<-shutdownCh

	log.Println("[local-game-simulation] shutting down gracefully")
}

func handleShutdownSignal(onShutdownReceived func()) chan struct{} {
	shutdownCh := make(chan struct{})

	go func() {
		sigNotifier := make(chan os.Signal, 1)

		signal.Notify(sigNotifier, os.Interrupt, syscall.SIGTERM)

		// Park here until a signal is received
		<-sigNotifier

		onShutdownReceived()
		close(shutdownCh)
	}()

	return shutdownCh
}
