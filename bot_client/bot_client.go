package bot_client

import (
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/server"
	"github.com/CameronHonis/chess-bot-server/engines"
)

type BotClient interface {
	Initialize(match *server.Match) error
	GenerateMove(match *server.Match) (*chess.Move, error)
	Terminate(match *server.Match) error
}

type LocalBotClient struct {
	engine engines.Engine
}

func NewLocalBotClient(engineName string) (BotClient, error) {
	engine, err := engines.GetLocalEngine(engineName)
	if err != nil {
		return nil, err
	}
	botClient := LocalBotClient{
		engine: engine,
	}
	return botClient, nil
}

func (lbc LocalBotClient) Initialize(match *server.Match) error {
	lbc.engine.Initialize(match)
	return nil
}

func (lbc LocalBotClient) GenerateMove(match *server.Match) (*chess.Move, error) {
	return lbc.engine.GenerateMove(match)
}

func (lbc LocalBotClient) Terminate(match *server.Match) error {
	lbc.engine.Terminate(match)
	return nil
}
