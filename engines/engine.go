package engines

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/server"
)

type Engine interface {
	Initialize(match *server.Match)
	GenerateMove(match *server.Match) *chess.Move
	Terminate(match *server.Match)
}

func GetLocalEngine(engineName string) (*Engine, error) {
	switch engineName {
	default:
		return nil, fmt.Errorf("unimplemented engine %s", engineName)
	}
}
