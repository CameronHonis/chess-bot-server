package main

import (
	"github.com/CameronHonis/chess-bot-server/app"
	"github.com/CameronHonis/chess-bot-server/arbitrator_client"
	"sync"
)

func main() {
	app.LoggerConfig()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	arbitrator_client.GetArbitratorClient()
	wg.Wait()
}
