package server

import (
	"errors"
	"fmt"
	"game-server/internal/game"
	"game-server/internal/player"
	"game-server/internal/waitingroom"
	"log"
	"sync"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type GameServer interface {
	Join(*player.Player) error
	Loop()
}

type RandomGameServer struct {
	isAcceptingNewPlayers       bool
	Games       cmap.ConcurrentMap[string, *game.RandomGame]
	gamesMutex  sync.RWMutex
	WaitingRoom *waitingroom.WaitingRoom
}

func New() *RandomGameServer {
	return &RandomGameServer{
		Games:       cmap.New[*game.RandomGame](),
		WaitingRoom: waitingroom.New(),
		isAcceptingNewPlayers:       true,
	}
}

func (rgs *RandomGameServer) Join(p *player.Player) error {
	if !rgs.isAcceptingNewPlayers {
		return errors.New("not accepting new players")
	}

	rgs.WaitingRoom.Sit([]*player.Player{p})

	return nil
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

				msg := fmt.Sprintf("[%s] [game-starter] new game started with players: %s and %s, it will end at %s", g.ID, pair[0].ID, pair[1].ID, g.EndDate.String())
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

			msg1 := fmt.Sprintf("[%s] [game-ender] received game over message", gameID)
			log.Default().Println(msg1)

			rgs.Games.Pop(gameID)

			msg2 := fmt.Sprintf("[%s] [game-ender] game removed, %d games left", gameID, rgs.Games.Count())
			log.Default().Printf(msg2)
		}
	}()

	// Kill random player
	go func() {
		for {
			time.Sleep(10 * time.Second)

			rgs.WaitingRoom.KillRandom()
		}
	}()

	// Games over cleaner
	go func() {
		for {
			time.Sleep(8 * time.Second)

			var ids []string

			if rgs.Games.Count() == 0 {
				log.Printf("[game-over-cleaner] no games dangling")
				continue
			}

			rgs.gamesMutex.Lock()

			for game := range rgs.Games.IterBuffered() {
				if game.Val.IsOver() {
					ids = append(ids, game.Val.ID)
				}
			}

			log.Printf("[game-over-cleaner] removing %d games dangling...", len(ids))

			for _, v := range ids {
				rgs.Games.Pop(v)
			}

			rgs.gamesMutex.Unlock()
		}
	}()

	// Print Stats
	go func() {
		for {
			log.Default().Printf("[stats-printer] %d games active, %d players waiting", rgs.Games.Count(), rgs.WaitingRoom.PlayersWaiting())

			time.Sleep(1 * time.Second)
		}
	}()
}
