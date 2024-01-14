package main

import (
	"github.com/CameronHonis/chess-bot-server/app"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	appService := app.Setup()
	appService.ArbitratorClient.
		wg.Wait()
}
