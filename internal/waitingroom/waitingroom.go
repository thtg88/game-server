package waitingroom

import (
	"log"
	"main/internal/player"
	"math/rand"
	"sync"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type WaitingRoom struct {
	Players      cmap.ConcurrentMap[string, *player.Player]
	PlayersMutex sync.RWMutex
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
		// log.Default().Printf("player %s (level %d) sat at the waiting room", p.ID, p.Level)
	}

	wr.PlayersMutex.Lock()
	defer wr.PlayersMutex.Unlock()

	wr.Players.MSet(pMap)
}

func (wr *WaitingRoom) Pair() []*player.Player {
	wr.PlayersMutex.Lock()
	defer wr.PlayersMutex.Unlock()

	pair := []*player.Player{wr.RandomPlayerWaiting()}
	for len(pair) < 2 {
		p := wr.RandomPlayerWaiting()

		if pair[0].ID != p.ID {
			pair = append(pair, p)
		}
	}

	wr.Players.Remove(pair[0].ID)
	wr.Players.Remove(pair[1].ID)

	return pair
}

func (wr *WaitingRoom) PlayersWaiting() int {
	wr.PlayersMutex.RLock()
	defer wr.PlayersMutex.RUnlock()

	return wr.Players.Count()
}

func (wr *WaitingRoom) KillRandom() {
	wr.PlayersMutex.Lock()
	defer wr.PlayersMutex.Unlock()

	condemned := wr.RandomPlayerWaiting()

	if condemned == nil {
		return
	}

	log.Default().Printf("killing %s", condemned.ID)

	wr.Players.Remove(condemned.ID)
}

func (wr *WaitingRoom) RandomPlayerWaiting() *player.Player {
	wr.PlayersMutex.RLock()
	defer wr.PlayersMutex.RUnlock()

	playersWaiting := wr.PlayersWaiting()

	if playersWaiting == 0 {
		return nil
	}

	ids := wr.Players.Keys()
	killedIdx := rand.Intn(len(ids))

	playerID := ids[killedIdx]

	player, _ := wr.Players.Get(playerID)

	return player
}
