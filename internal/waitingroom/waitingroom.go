package waitingroom

import (
	"log"
	"main/internal/player"
	"math/rand"
	"sync"
)

type WaitingRoom struct {
	Players      []*player.Player
	PlayersMutex sync.RWMutex
}

func (wg *WaitingRoom) Sit(players []*player.Player) {
	wg.PlayersMutex.Lock()
	wg.Players = append(wg.Players, players...)
	// for _, p := range players {
	// 	log.Default().Printf("player %s (level %d) sat at the waiting room", p.ID, p.Level)
	// }
	wg.PlayersMutex.Unlock()

}

func (wg *WaitingRoom) PlayersWaiting() int {
	wg.PlayersMutex.RLock()
	defer wg.PlayersMutex.RUnlock()

	// the poor player left alone cant be paired with anybody else anyway
	return len(wg.Players)
}

func (wg *WaitingRoom) KillRandom() {
	wg.PlayersMutex.Lock()
	defer wg.PlayersMutex.Unlock()

	playersWaiting := len(wg.Players)

	if playersWaiting == 0 {
		return
	}

	killedIdx := rand.Intn(playersWaiting)

	log.Default().Printf("killing %s", wg.Players[killedIdx].ID)

	wg.Players[killedIdx] = wg.Players[playersWaiting-1]

	wg.Players = wg.Players[:playersWaiting-1]
}
