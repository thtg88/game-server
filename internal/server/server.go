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
	Shutdown()
}

type RandomGameServer struct {
	canKillRandomWaitingPlayers bool
	canPrintStats               bool
	isAcceptingNewPlayers       bool

	gameOverCh chan string

	Games       cmap.ConcurrentMap[string, *game.RandomGame]
	gamesMutex  sync.RWMutex
	WaitingRoom *waitingroom.WaitingRoom
}

func New() *RandomGameServer {
	return &RandomGameServer{
		canKillRandomWaitingPlayers: true,
		canPrintStats:               true,
		isAcceptingNewPlayers:       true,
		gameOverCh:                  make(chan string),
		Games:                       cmap.New[*game.RandomGame](),
		WaitingRoom:                 waitingroom.New(),
	}
}

func (rgs *RandomGameServer) Shutdown() {
	rgs.isAcceptingNewPlayers = false
	rgs.canKillRandomWaitingPlayers = false

	for rgs.Games.Count() > 0 {
		time.Sleep(200 * time.Millisecond)
	}

	rgs.WaitingRoom.KillAll()

	rgs.canPrintStats = false
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

	go rgs.startNewGames()
	go rgs.endGamesOver()
	go rgs.killRandomWaitingPlayers()
	go rgs.cleanDanglingGamesOver()
	go rgs.printStats()
}

func (rgs *RandomGameServer) startNewGames() {
	for rgs.isAcceptingNewPlayers {
		for rgs.WaitingRoom.PlayersWaiting() >= 2 {
			// log.Default().Println("[game-starter] players waiting")

			rgs.gamesMutex.Lock()

			pair := rgs.WaitingRoom.Pair()
			g := game.New(pair)
			rgs.Games.Set(g.ID, g)
			go g.Start(rgs.gameOverCh)

			rgs.gamesMutex.Unlock()

			msg := fmt.Sprintf("[%s] [game-starter] new game started with players: %s and %s, it will end at %s", g.ID, pair[0].ID, pair[1].ID, g.EndDate.String())
			g.SendMsgs(msg)
			log.Default().Println(msg)
		}

		time.Sleep(4 * time.Millisecond)
	}
}

func (rgs *RandomGameServer) endGamesOver() {
	for {
		gameID := <-rgs.gameOverCh

		msg1 := fmt.Sprintf("[%s] [game-ender] received game over message", gameID)
		log.Default().Println(msg1)

		rgs.Games.Pop(gameID)

		msg2 := fmt.Sprintf("[%s] [game-ender] game removed, %d games left", gameID, rgs.Games.Count())
		log.Default().Printf(msg2)
	}
}

func (rgs *RandomGameServer) killRandomWaitingPlayers() {
	for rgs.canKillRandomWaitingPlayers {
		rgs.WaitingRoom.KillRandom()

		time.Sleep(10 * time.Second)
	}
}

func (rgs *RandomGameServer) cleanDanglingGamesOver() {
	for {
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

		time.Sleep(8 * time.Second)
	}
}

func (rgs *RandomGameServer) printStats() {
	for rgs.canPrintStats {
		log.Default().Printf("[stats-printer] %d games active, %d players waiting", rgs.Games.Count(), rgs.WaitingRoom.PlayersWaiting())

		time.Sleep(1 * time.Second)
	}
}
