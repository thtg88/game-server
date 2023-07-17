package waitingroom

import (
	"main/internal/player"
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
		// log.Default().Printf("player %s (level %d) sat at the waiting room", p.ID, p.Level)
	}

	wr.Players.MSet(pMap)
}

func (wr *WaitingRoom) Pair() []*player.Player {
	wr.playersMutex.Lock()
	defer wr.playersMutex.Unlock()

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
	return wr.Players.Count()
}

func (wr *WaitingRoom) KillRandom() {
	wr.playersMutex.Lock()
	defer wr.playersMutex.Unlock()

	condemned := wr.RandomPlayerWaiting()

	if condemned == nil {
		return
	}

	// log.Default().Printf("killing %s", condemned.ID)

	wr.Players.Remove(condemned.ID)
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
	killedIdx := rand.Intn(len(ids))

	return ids[killedIdx]
}
