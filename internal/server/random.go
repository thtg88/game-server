package server

import (
	"log"
	"main/internal/game"
	"main/internal/player"
	"main/internal/waitingroom"
	"sync"
	"time"
)

type RandomGameServer struct {
	Games            []*game.RandomGame
	gamesMutex       sync.RWMutex
	WaitingRoom      *waitingroom.WaitingRoom
	waitingRoomMutex sync.RWMutex
}

func New() *RandomGameServer {
	return &RandomGameServer{
		WaitingRoom: waitingroom.New(),
	}
}

func (rgs *RandomGameServer) Join(p *player.Player) {
	rgs.WaitingRoom.Sit([]*player.Player{p})
}

func (rgs *RandomGameServer) Loop() {
	log.Default().Println("server started")

	// Start new games
	go func() {
		for {
			for rgs.WaitingRoom.PlayersWaiting() >= 2 {
				log.Default().Println("players waiting")

				pair := rgs.WaitingRoom.Pair()

				rgs.gamesMutex.Lock()
				g := game.New(pair)
				rgs.Games = append(rgs.Games, g)
				log.Default().Printf("new game started with players: %s and %s, it will end at %s", pair[0].ID, pair[1].ID, g.EndDate.String())
				rgs.gamesMutex.Unlock()
			}

			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Game rounds
	go func() {
		for {
			var gameOverIdxs = make(map[int]bool)
			rgs.gamesMutex.RLock()
			for idx, game := range rgs.Games {
				// TODO: wrap iteration in a goroutine?

				// go func() { game.Round() }()

				if game.IsOver() {
					log.Default().Println("game over")
					// TODO: Increment player levels

					rgs.WaitingRoom.Sit([]*player.Player{game.Player1, game.Player2})
					gameOverIdxs[idx] = true
				}
			}
			rgs.gamesMutex.RUnlock()

			if len(gameOverIdxs) > 0 {
				log.Default().Printf("clearing up %d game(s) over", len(gameOverIdxs))
				rgs.gamesMutex.Lock()
				newGames := []*game.RandomGame{}
				for idx, game := range rgs.Games {
					if _, ok := gameOverIdxs[idx]; !ok {
						newGames = append(newGames, game)
					}
				}
				rgs.Games = newGames
				rgs.gamesMutex.Unlock()
			}

			time.Sleep(200 * time.Millisecond)
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
			rgs.gamesMutex.RLock()
			log.Default().Printf("%d games active, %d players waiting", len(rgs.Games), rgs.WaitingRoom.PlayersWaiting())
			rgs.gamesMutex.RUnlock()

			time.Sleep(1 * time.Second)
		}
	}()
}
