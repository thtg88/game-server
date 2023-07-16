package client

import (
	"main/internal/player"
	"main/internal/server"
	"time"
)

type RandomClient struct {
	Player *player.Player
}

func New() *RandomClient {
	p := player.Random()

	return &RandomClient{Player: &p}
}

func Loop(rgs *server.RandomGameServer) {
	for {
		c := New()

		// log.Default().Println("new client")

		rgs.WaitingRoom.Sit([]*player.Player{c.Player})

		time.Sleep(1 * time.Second)
	}
}
