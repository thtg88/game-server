package player

import (
	"math/rand"

	"github.com/google/uuid"
)

type Player struct {
	GameOverCh chan struct{}
	ID         string
	Level      uint64
	MessagesCh chan string
}

func Random() Player {
	return Player{
		GameOverCh: make(chan struct{}),
		ID:         uuid.NewString(),
		Level:      uint64(rand.Intn(1000)) + 1,
		MessagesCh: make(chan string),
	}
}

func (p *Player) SendMsgs(msgs ...string) {
	for _, msg := range msgs {
		p.MessagesCh <- msg
	}
}
