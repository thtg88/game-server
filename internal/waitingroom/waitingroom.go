package waitingroom

import (
	"fmt"
	"game-server/internal/player"
	"log"
	"math/rand"
	"sync"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type WaitingRoom struct {
	Players      cmap.ConcurrentMap[string, *player.Player]
	playersMutex sync.RWMutex
}

func New() *WaitingRoom {
	return &WaitingRoom{
		Players: cmap.New[*player.Player](),
	}
}

func (wr *WaitingRoom) Sit(players []*player.Player) {
	pMap := make(map[string]*player.Player)
	for _, p := range players {
		pMap[p.ID] = p
		msg := fmt.Sprintf("player %s (level %d) sat at the waiting room", p.ID, p.Level)
		p.SendMsgs(msg)
		log.Default().Printf(msg)
	}

	wr.Players.MSet(pMap)
}

func (wr *WaitingRoom) Pair() []*player.Player {
	pair := []*player.Player{wr.RandomPlayerWaiting()}
	wr.Players.Remove(pair[0].ID)

	for len(pair) < 2 {
		p := wr.RandomPlayerWaiting()

		if p != nil && pair[0].ID != p.ID {
			pair = append(pair, p)
			wr.Players.Remove(pair[1].ID)
		}
	}

	return pair
}

func (wr *WaitingRoom) PlayersWaiting() int {
	return wr.Players.Count()
}

func (wr *WaitingRoom) KillRandom() {
	wr.Players.Remove(wr.RandomPlayerKey())
}

func (wr *WaitingRoom) RandomPlayerWaiting() *player.Player {
	playersWaiting := wr.PlayersWaiting()

	if playersWaiting == 0 {
		return nil
	}

	playerID := wr.RandomPlayerKey()

	player, _ := wr.Players.Get(playerID)

	return player
}

func (wr *WaitingRoom) RandomPlayerKey() string {
	ids := wr.Players.Keys()

	if len(ids) == 0 {
		return ""
	}

	killedIdx := rand.Intn(len(ids))

	return ids[killedIdx]
}
