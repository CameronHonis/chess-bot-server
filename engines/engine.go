package engines

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-bot-server/engines/random"
)

type Engine interface {
	Initialize(match *models.Match)
	GenerateMove(match *models.Match) (*chess.Move, error)
	Terminate()
}

func GetLocalEngine(engineName string) (Engine, error) {
	switch engineName {
	case "random":
		return &random.Engine{}, nil
	default:
		return nil, fmt.Errorf("unimplemented engine %s", engineName)
	}
}
