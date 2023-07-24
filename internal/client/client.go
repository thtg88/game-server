package client

import (
	"game-server/internal/player"
	"game-server/internal/server"
	"log"
	"sync"
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
	var wg sync.WaitGroup

	defer wg.Wait()

	for i := 0; i < 10000; i++ {
		c := New()

		err := gs.Join(c.Player)
		if err != nil {
			log.Printf("[client-loop] error joining the server: %v", err)
			break
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			c.gameLoop()
		}()

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
}
