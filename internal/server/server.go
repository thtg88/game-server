package server

import (
	"main/internal/player"
)

type GameServer interface {
	Join(*player.Player)
	Loop()
}
