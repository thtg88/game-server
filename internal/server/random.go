package server

import (
	"log"
	"main/internal/game"
	"main/internal/player"
	"main/internal/waitingroom"
	"sync"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type RandomGameServer struct {
	Games       cmap.ConcurrentMap[string, *game.RandomGame]
	gamesMutex  sync.RWMutex
	WaitingRoom *waitingroom.WaitingRoom
}

func New() *RandomGameServer {
	return &RandomGameServer{
		Games:       cmap.New[*game.RandomGame](),
		WaitingRoom: waitingroom.New(),
	}
}

func (rgs *RandomGameServer) Join(p *player.Player) {
	rgs.WaitingRoom.Sit([]*player.Player{p})
}

func (rgs *RandomGameServer) Loop() {
	log.Default().Println("server started")

	gameOverCh := make(chan string)

	// Start new games
	go func() {
		for {
			for rgs.WaitingRoom.PlayersWaiting() >= 2 {
				// log.Default().Println("[game-starter] players waiting")

				rgs.gamesMutex.Lock()

				pair := rgs.WaitingRoom.Pair()
				g := game.New(pair)
				rgs.Games.Set(g.ID, g)
				go g.Start(gameOverCh)

				rgs.gamesMutex.Unlock()

				// log.Default().Printf("[game-starter] new game started with players: %s and %s, it will end at %s", pair[0].ID, pair[1].ID, g.EndDate.String())
			}

			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Game over channel
	go func() {
		for {
			gameID := <-gameOverCh
			// log.Default().Printf("[game-ender] received game over message for ID: %s", gameID)
			game, ok := rgs.Games.Pop(gameID)
			// log.Default().Printf("[game-ender] game %s removed", gameID)
			if ok {
				// log.Default().Printf("[game-ender] sitting players %s and %s...", game.Player1.ID, game.Player2.ID)
				rgs.WaitingRoom.Sit([]*player.Player{game.Player1, game.Player2})
			}
		}
	}()

	// Kill random player
	// go func() {
	// 	for {
	// 		rgs.WaitingRoom.KillRandom()

	// 		time.Sleep(1500 * time.Millisecond)
	// 	}
	// }()

	// Print Stats
	go func() {
		for {
			log.Default().Printf("[stats-printer] %d games active, %d players waiting", rgs.Games.Count(), rgs.WaitingRoom.PlayersWaiting())

			time.Sleep(1 * time.Second)
		}
	}()
}
