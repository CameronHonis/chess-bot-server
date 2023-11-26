package main

import (
	"github.com/CameronHonis/chess-bot-server/bot_client"
	"sync"
)

func main() {
	bot_client.ConfigLogger()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	bot_client.GetArbitratorClient()
	wg.Wait()
}
