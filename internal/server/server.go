package server

import (
	"errors"
	"fmt"
	"game-server/internal/game"
	"game-server/internal/player"
	"game-server/internal/waitingroom"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type GameServer interface {
	Join(*player.Player) error
	Loop()
	Shutdown()
}

type RandomGameServer struct {
	config      *config
	Games       cmap.ConcurrentMap[string, *game.RandomGame]
	gamesMutex  sync.RWMutex
	WaitingRoom *waitingroom.WaitingRoom
}

type config struct {
	canCleanGamesDangling       bool
	canKillRandomWaitingPlayers bool
	canPrintStats               bool
	isAcceptingNewPlayers       bool
}

func New() *RandomGameServer {
	return &RandomGameServer{
		config: &config{
			canCleanGamesDangling:       true,
			canKillRandomWaitingPlayers: true,
			canPrintStats:               true,
			isAcceptingNewPlayers:       true,
		},
		Games:       cmap.New[*game.RandomGame](),
		WaitingRoom: waitingroom.New(),
	}
}

func HandleShutdownSignal(onShutdownReceived func()) chan struct{} {
	shutdownCh := make(chan struct{})

	go func() {
		sigNotifier := make(chan os.Signal, 1)

		signal.Notify(sigNotifier, os.Interrupt, syscall.SIGTERM)

		// Park here until a signal is received
		<-sigNotifier

		onShutdownReceived()
		close(shutdownCh)
	}()

	return shutdownCh
}

func (rgs *RandomGameServer) Shutdown() {
	log.Println("[random-game-server] shutting down...")

	rgs.config.isAcceptingNewPlayers = false
	rgs.config.canKillRandomWaitingPlayers = false

	// wait for all games to be over
	for rgs.Games.Count() > 0 {
		time.Sleep(100 * time.Millisecond)
	}

	rgs.config.canCleanGamesDangling = false

	rgs.WaitingRoom.KillAll()

	rgs.config.canPrintStats = false
}

func (rgs *RandomGameServer) Join(p *player.Player) error {
	if !rgs.config.isAcceptingNewPlayers {
		return errors.New("not accepting new players")
	}

	rgs.WaitingRoom.Sit([]*player.Player{p})

	return nil
}

func (rgs *RandomGameServer) Loop() {
	log.Println("server started")

	go rgs.startNewGames()
	go rgs.killRandomWaitingPlayers()
	go rgs.cleanDanglingGamesOver()
	go rgs.printStats()
}

func (rgs *RandomGameServer) startNewGames() {
	for rgs.config.isAcceptingNewPlayers {
		for rgs.WaitingRoom.PlayersWaiting() >= 2 {
			// log.Println("[game-starter] players waiting")

			rgs.gamesMutex.Lock()

			pair := rgs.WaitingRoom.Pair()
			g := game.New(pair)
			rgs.Games.Set(g.ID, g)
			go g.Start()
			go rgs.waitForGameOver(g)

			rgs.gamesMutex.Unlock()

			msg := fmt.Sprintf("[%s] [game-starter] new game started with players: %s and %s, it will end at %s", g.ID, pair[0].ID, pair[1].ID, g.EndDate.String())
			g.SendMsgs(msg)
			log.Println(msg)
		}

		time.Sleep(4 * time.Millisecond)
	}

	log.Println("[game-starter] stopped accepting new players")
}

func (rgs *RandomGameServer) waitForGameOver(g *game.RandomGame) {
	<-g.OverCh

	rgs.Games.Remove(g.ID)

	log.Printf("[%s] [game-ender] game removed", g.ID)
}

func (rgs *RandomGameServer) killRandomWaitingPlayers() {
	for rgs.config.canKillRandomWaitingPlayers {
		rgs.WaitingRoom.KillRandom()

		time.Sleep(10 * time.Second)
	}

	log.Println("[random-player-killer] stopped killing")
}

func (rgs *RandomGameServer) cleanDanglingGamesOver() {
	for rgs.config.canCleanGamesDangling {
		if rgs.Games.Count() == 0 {
			log.Println("[game-over-cleaner] no games dangling")
			time.Sleep(8 * time.Second)
			continue
		}

		var ids []string

		rgs.gamesMutex.Lock()

		for game := range rgs.Games.IterBuffered() {
			if game.Val.IsOver() {
				ids = append(ids, game.Val.ID)
			}
		}

		if len(ids) == 0 {
			log.Println("[game-over-cleaner] no games dangling")
			rgs.gamesMutex.Unlock()
			time.Sleep(8 * time.Second)
			continue
		}

		log.Printf("[game-over-cleaner] removing %d games dangling...", len(ids))

		for _, v := range ids {
			rgs.Games.Pop(v)
		}

		rgs.gamesMutex.Unlock()

		time.Sleep(8 * time.Second)
	}

	log.Println("[game-over-cleaner] stopped cleaning")
}

func (rgs *RandomGameServer) printStats() {
	for rgs.config.canPrintStats {
		log.Printf("[stats-printer] %d games active, %d players waiting", rgs.Games.Count(), rgs.WaitingRoom.PlayersWaiting())

		time.Sleep(1 * time.Second)
	}

	log.Println("[stats-printer] stopped printing")
}
