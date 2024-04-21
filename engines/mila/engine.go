package mila

import (
	"context"
	"fmt"
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-bot-server/uci_client"
	"github.com/CameronHonis/chess-bot-server/uci_client/cmd_client"
	"os/exec"
	"time"
)

type Engine struct {
	cmd    *exec.Cmd
	client *uci_client.Client
}

func NewEngine(cmd *exec.Cmd) (*Engine, error) {
	cmdClient, cmdClientErr := cmd_client.ClientFromCmd(cmd)
	if cmdClientErr != nil {
		return nil, fmt.Errorf("could not construct CmdClient: %s", cmdClientErr)
	}

	return &Engine{
		cmd:    cmd,
		client: uci_client.NewUciClient(cmdClient),
	}, nil
}

func (e *Engine) Initialize(match *models.Match) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), time.Second)
	defer cancelCtx()

	startErr := e.cmd.Start()
	if startErr != nil {
		return fmt.Errorf("could not start stockfish: %s", startErr)
	}

	_, readErr := e.client.CmdClient.ReadLine(ctx)
	if readErr != nil {
		return fmt.Errorf("could not read startup msg: %s", readErr)
	}

	_, initErr := e.client.Init(ctx)
	if initErr != nil {
		return initErr
	}

	return nil
}

func (e *Engine) GenerateMove(match *models.Match) (*chess.Move, error) {
	searchOpts := uci_client.NewSearchOptionsBuilder().
		WithWhiteMs(uint(match.WhiteTimeRemainingSec * 1000.)).
		WithBlackMs(uint(match.BlackTimeRemainingSec * 1000.)).
		Build()
	setPosErr := e.client.SetPosition(match.Board.ToFEN())
	if setPosErr != nil {
		return nil, fmt.Errorf("could not set position: %s", setPosErr)
	}

	readyCtx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancelCtx()
	for {
		isReady, isReadyErr := e.client.IsReady(readyCtx)
		if isReadyErr != nil {
			return nil, fmt.Errorf("could not read ready state of engine: %s", isReadyErr)
		}
		if isReady {
			break
		}
	}

	var secsRemaining float64
	if match.Board.IsWhiteTurn {
		secsRemaining = match.WhiteTimeRemainingSec
	} else {
		secsRemaining = match.BlackTimeRemainingSec
	}
	genMoveCtx, cancelGenMoveCtx := context.WithTimeout(context.Background(), time.Duration(secsRemaining+1)*time.Second)
	defer cancelGenMoveCtx()

	bestMoveLAlg, searchErr := e.client.Go(genMoveCtx, searchOpts)
	if searchErr != nil {
		return nil, fmt.Errorf("error reading best move: %s", searchErr)
	}
	bestMove, moveConvertErr := chess.MoveFromAlgebraic(bestMoveLAlg, match.Board)
	if moveConvertErr != nil {
		return nil, fmt.Errorf("could not convert to move: %s", bestMove)
	}
	return bestMove, nil
}

func (e *Engine) Terminate() {
	if endErr := e.client.End(); endErr != nil {
		fmt.Println("WARN: could not end client: ", endErr)
	}
}

func (e *Engine) SetOption(ctx context.Context, optName, optValue string) error {
	return e.client.SetOption(ctx, optName, optValue)
}
