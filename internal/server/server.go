package server

import (
	"main/internal/player"
	"main/internal/waitingroom"
)

type GameServer interface {
	Join()
	Loop()
	Pair() (*waitingroom.WaitingRoom, []*player.Player)
}
