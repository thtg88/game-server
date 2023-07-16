package player

import (
	"math/rand"

	"github.com/google/uuid"
)

type Player struct {
	ID    string
	Level uint64
}

func Random() Player {
	return Player{
		ID:    uuid.New().String(),
		Level: uint64(rand.Intn(1000)) + 1,
	}
}
