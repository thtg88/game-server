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

func Spawn(gs server.GameServer) {
	for {
		c := New()

		// log.Println("new client")

		err := gs.Join(c.Player)
		if err != nil {
			log.Printf("error joining the server: %v", err)
			break
		}

		go c.gameLoop()

		time.Sleep(50 * time.Millisecond)
	}
}

func (c *RandomClient) gameLoop() {
	gameOver := false

	for !gameOver {
		select {
		case <-c.Player.MessagesCh:
			// log.Printf("[client-loop] got message %s from server", msg)
		case <-c.Player.GameOverCh:
			log.Printf("[%s] [client-loop] it's game over!", c.Player.ID)
			gameOver = true
		}
	}

	close(c.Player.GameOverCh)
	close(c.Player.MessagesCh)
}
