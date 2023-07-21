package main

import "game-server/internal/remoteclient"

func main() {
	rc := remoteclient.New()

	rc.Join()
}
