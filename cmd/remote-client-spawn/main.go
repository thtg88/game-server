package main

import (
	"game-server/internal/remoteclient"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for j := 0; j < 10000; j++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			rc := remoteclient.NewTcpSocketRandomClient()

			if err := rc.Join(); err != nil {
				log.Printf("%v", err)
			}
		}()
	}

	wg.Wait()
}
