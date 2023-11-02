package random

import (
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/server"
	"math/rand"
	"time"
)

type Engine struct {
}

func (re *Engine) Initialize(match *server.Match) {

}

func (re *Engine) GenerateMove(match *server.Match) *chess.Move {
	movesBySquare, _ := chess.GetLegalMoves(match.Board, false)
	moves := make([]*chess.Move, 0)
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			for _, move := range movesBySquare[r][c] {
				moves = append(moves, move)
			}
		}
	}
	rand.Seed(time.Now().UnixNano())
	randMoveIdx := rand.Intn(len(moves))
	return moves[randMoveIdx]
}

func (re *Engine) Terminate(match *server.Match) {

}