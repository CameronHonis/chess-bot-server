package engines

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-bot-server/engines/mila"
	"github.com/CameronHonis/chess-bot-server/engines/random"
	"github.com/CameronHonis/chess-bot-server/engines/stockfish"
	"os"
	"os/exec"
)

type Engine interface {
	Initialize(match *models.Match) error
	GenerateMove(match *models.Match) (*chess.Move, error)
	Terminate()
}

func EngineFromName(engineName string) (Engine, error) {
	switch engineName {
	case "random":
		return &random.Engine{}, nil
	case "stockfish":
		cmd, stockfishErr := StockfishCmd()
		if stockfishErr != nil {
			return nil, stockfishErr
		}

		engine, engineErr := stockfish.NewEngine(cmd)
		if engineErr != nil {
			return nil, fmt.Errorf("could not make stockfish engine: %s", engineErr)
		}
		return engine, nil
	case "mila":
		cmd, milaErr := MilaCmd()
		if milaErr != nil {
			return nil, milaErr
		}

		engine, engineErr := mila.NewEngine(cmd)
		if engineErr != nil {
			return nil, fmt.Errorf("could not make mila engine: %s", engineErr)
		}
		return engine, nil
	default:
		return nil, fmt.Errorf("unimplemented engine %s", engineName)
	}
}

func StockfishCmd() (*exec.Cmd, error) {
	path, pathExists := os.LookupEnv("STOCKFISH_PATH")
	if !pathExists {
		return nil, fmt.Errorf("stockfish path not found")
	}

	return exec.Command(path), nil
}

func MilaCmd() (*exec.Cmd, error) {
	path, pathExists := os.LookupEnv("MILA_PATH")
	if !pathExists {
		return nil, fmt.Errorf("mila path not found")
	}

	return exec.Command(path), nil
}
