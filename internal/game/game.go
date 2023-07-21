package game

import (
	"fmt"
	"game-server/internal/player"
	"log"
	"time"

	"github.com/google/uuid"
)

type Game interface {
	Round()
	Start()
	IsOver() bool
}

type RandomGame struct {
	EndDate *time.Time
	ID      string
	Player1 *player.Player
	Player2 *player.Player
}

func New(players []*player.Player) *RandomGame {
	if len(players) != 2 {
		panic("only 2 players supported by this game")
	}

	// offset := rand.Intn(3) + 1
	offset := 1

	now := time.Now()
	endDate := now.Add(time.Second * time.Duration(offset))

	return &RandomGame{
		ID:      uuid.NewString(),
		Player1: players[0],
		Player2: players[1],
		EndDate: &endDate,
	}
}

func (rg *RandomGame) Start(gameOverCh chan<- string) {
	for i := 0; !rg.IsOver(); i++ {
		rg.SendMsgs(fmt.Sprintf("round %d starting...", i))
		rg.Round()
		rg.SendMsgs(fmt.Sprintf("round %d over!", i))
	}

	msg := fmt.Sprintf("[%s] [game] game over", rg.ID)
	log.Default().Println(msg)

	// TODO: Increment player levels

	// First let the clients disconnect
	rg.Player1.GameOverCh <- true
	rg.Player2.GameOverCh <- true

	// Then let the server clean up the game
	gameOverCh <- rg.ID
}

func (rg *RandomGame) Round() {
	time.Sleep(500 * time.Millisecond)
}

func (rg *RandomGame) IsOver() bool {
	// end date is in the past
	return rg.EndDate.Before(time.Now())
}

func (rg *RandomGame) SendMsgs(msgs ...string) {
	rg.Player1.SendMsgs(msgs...)
	rg.Player2.SendMsgs(msgs...)
}
