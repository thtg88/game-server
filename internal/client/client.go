package client

import (
	"game-server/internal/player"
	"game-server/internal/server"
	"log"
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

		// Consume channels
		go func(p *player.Player) {
			gameOver := false

			for !gameOver {
				select {
				case <-p.MessagesCh:
					// log.Printf("[client-loop] got message %s from server", msg)
				case <-p.GameOverCh:
					log.Printf("[%s] [client-loop] it's game over!", p.ID)
					gameOver = true
				}
			}

			close(p.GameOverCh)
			close(p.MessagesCh)
		}(c.Player)

		time.Sleep(50 * time.Millisecond)
	}
}
