package main

import (
	"game-server/internal/remoteclient"
	"log"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	for {
		for i := 0; i < 10000; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				rc := remoteclient.New()

				if err := rc.Join(); err != nil {
					log.Printf("%v", err)
				}
			}()
		}

		wg.Wait()

		time.Sleep(10 * time.Second)
	}
}
