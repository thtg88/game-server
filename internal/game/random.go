package game

import (
	"main/internal/player"
	"time"

	"github.com/google/uuid"
)

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
	for !rg.IsOver() {
		rg.Round()
	}

	// log.Default().Printf("[game] game %s over", rg.ID)

	// TODO: Increment player levels

	gameOverCh <- rg.ID
}

func (rg *RandomGame) Round() {
	time.Sleep(5 * time.Second)
}

func (rg *RandomGame) IsOver() bool {
	// end date is in the past
	return rg.EndDate.Before(time.Now())
}
