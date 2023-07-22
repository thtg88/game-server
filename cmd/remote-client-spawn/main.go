package main

import (
	"game-server/internal/remoteclient"
	"log"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	for j := 0; j < 10000; j++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			rc := remoteclient.NewGrpcRandomClient()

			if err := rc.Join(); err != nil {
				log.Printf("%v", err)
			}
		}()
	}

	wg.Wait()

	time.Sleep(10 * time.Second)
}
