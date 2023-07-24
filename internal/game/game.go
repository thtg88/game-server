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
	OverCh  chan struct{}
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
		OverCh:  make(chan struct{}),
		Player1: players[0],
		Player2: players[1],
		EndDate: &endDate,
	}
}

func (rg *RandomGame) Start() {
	for i := 0; !rg.IsOver(); i++ {
		rg.Round(i)
	}

	msg := fmt.Sprintf("[%s] [game] game over", rg.ID)
	log.Println(msg)

	// TODO: Increment player levels

	// First let the clients disconnect
	close(rg.Player1.GameOverCh)
	close(rg.Player1.MessagesCh)
	close(rg.Player2.GameOverCh)
	close(rg.Player2.MessagesCh)

	// Then let the server clean up the game
	close(rg.OverCh)
}

func (rg *RandomGame) Round(round int) {
	rg.SendMsgs(fmt.Sprintf("[%s] [game] round %d starting...", rg.ID, round))
	time.Sleep(500 * time.Millisecond)
	rg.SendMsgs(fmt.Sprintf("[%s] [game] round %d over!", rg.ID, round))
}

func (rg *RandomGame) IsOver() bool {
	// end date is in the past
	return rg.EndDate.Before(time.Now())
}

func (rg *RandomGame) SendMsgs(msgs ...string) {
	rg.Player1.SendMsgs(msgs...)
	rg.Player2.SendMsgs(msgs...)
}
