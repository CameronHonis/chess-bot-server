package main

import (
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	GetArbitratorClient()
	wg.Wait()
}
