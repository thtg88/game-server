package server

import (
	"fmt"
	"log"
	"main/internal/game"
	"main/internal/player"
	"main/internal/waitingroom"
	"sync"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type GameServer interface {
	Join(*player.Player)
	Loop()
}

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

				msg := fmt.Sprintf("[game-starter] new game started with players: %s and %s, it will end at %s", pair[0].ID, pair[1].ID, g.EndDate.String())
				g.SendMsgs(msg)
				log.Default().Println(msg)
			}

			time.Sleep(4 * time.Millisecond)
		}
	}()

	// Game over channel
	go func() {
		for {
			gameID := <-gameOverCh

			msg1 := fmt.Sprintf("[game-ender] received game over message for ID: %s", gameID)
			log.Default().Println(msg1)

			game, ok := rgs.Games.Pop(gameID)

			msg2 := fmt.Sprintf("[game-ender] game %s removed", gameID)
			log.Default().Printf(msg2)

			if ok {
				msg3 := fmt.Sprintf("[game-ender] sitting players %s and %s...", game.Player1.ID, game.Player2.ID)
				log.Default().Printf(msg3)
				game.SendMsgs(msg1, msg2, msg3)

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
