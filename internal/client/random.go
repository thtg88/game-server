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

func Loop(gs server.GameServer) {
	for {
		c := New()

		// log.Default().Println("new client")

		gs.Join(c.Player)

		time.Sleep(50 * time.Millisecond)
	}
}
