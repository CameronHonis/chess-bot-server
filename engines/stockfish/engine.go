package stockfish

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/models"
	"math/rand"
	"os/exec"
)

type Engine struct {
	cmd *exec.Cmd
}

func (re *Engine) Initialize(match *models.Match) {

}

func (re *Engine) GenerateMove(match *models.Match) (*chess.Move, error) {
	moves, movesErr := chess.GetLegalMoves(match.Board)
	if movesErr != nil {
		return nil, fmt.Errorf("cannot generate move: %s", movesErr)
	}

	if len(moves) == 0 {
		return nil, fmt.Errorf("no legal moves")
	}
	randMoveIdx := rand.Intn(len(moves))
	return moves[randMoveIdx], nil
}

func (re *Engine) Terminate() {

}
