package client

import (
	"game-server/internal/player"
	"game-server/internal/server"
	"time"
)

type Client interface {
}

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
