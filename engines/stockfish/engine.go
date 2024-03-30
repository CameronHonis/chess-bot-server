package stockfish

import (
	"context"
	"fmt"
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-bot-server/uci_client"
	"github.com/CameronHonis/chess-bot-server/uci_client/cmd_client"
	"os"
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

	startErr := cmd.Start()
	if startErr != nil {
		return nil, fmt.Errorf("could not start stockfish: %s", startErr)
	}

	return &Engine{
		cmd:    cmd,
		client: uci_client.NewUciClient(cmdClient),
	}, nil
}

func (e *Engine) Initialize(match *models.Match) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), time.Second)
	defer cancelCtx()

	_, initErr := e.client.Init(ctx)
	if initErr != nil {
		return initErr
	}

	optErr := e.SetOption(ctx, "Threads", "32")
	if optErr != nil {
		return fmt.Errorf("error setting option 'Threads' to 32: %s", optErr)
	}

	return nil
}

func (e *Engine) GenerateMove(match *models.Match) (*chess.Move, error) {
	time.Sleep(time.Second)
	searchOpts := uci_client.NewSearchOptionsBuilder().
		WithWhiteMs(uint(match.WhiteTimeRemainingSec * 1000.)).
		WithBlackMs(uint(match.BlackTimeRemainingSec * 1000.)).
		Build()
	setPosErr := e.client.SetPosition(match.Board.ToFEN())
	if setPosErr != nil {
		return nil, fmt.Errorf("could not set position: %s", setPosErr)
	}
	ctx, cancelCtx := context.WithTimeout(context.Background(), time.Second)
	defer cancelCtx()
	for {
		isReady, isReadyErr := e.client.IsReady(ctx)
		if isReadyErr != nil {
			return nil, fmt.Errorf("could not read ready state of engine: %s", isReadyErr)
		}
		if isReady {
			break
		}
	}
	bestMoveLAlg, searchErr := e.client.Go(ctx, searchOpts)
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
	interruptErr := e.cmd.Process.Signal(os.Interrupt)
	if interruptErr != nil {
		_ = e.cmd.Process.Kill()
	}
}

func (e *Engine) SetOption(ctx context.Context, optName, optValue string) error {
	return e.client.SetOption(ctx, optName, optValue)
}
