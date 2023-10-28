package main

import (
	"fmt"
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
	engine *engines.Engine
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
	(*lbc.engine).Initialize(match)
	return nil
}

func (lbc LocalBotClient) GenerateMove(match *server.Match) (*chess.Move, error) {
	return (*lbc.engine).GenerateMove(match), nil
}

func (lbc LocalBotClient) Terminate(match *server.Match) error {
	(*lbc.engine).Terminate(match)
	return nil
}

func (lbc LocalBotClient) PushMatchUpdate(match *server.Match) error {
	if match.Board.IsTerminal {
		removeErr := GetBotClientManager().RemoveBotClient(match)
		if removeErr != nil {
			return removeErr
		}
	}
	isBotWhite := match.WhiteClientId == ""
	isBotTurn := isBotWhite == match.Board.IsWhiteTurn
	if !isBotTurn {
		return nil
	}
	move, moveErr := lbc.GenerateMove(match)
	if moveErr != nil {
		return moveErr
	}
	msg := server.Message{
		Topic:       server.MessageTopic(fmt.Sprintf("move-%s", match.Uuid)),
		ContentType: server.CONTENT_TYPE_MOVE,
		Content: &server.MoveMessageContent{
			Move: move,
		},
	}
	return GetArbitratorClient().SendMessage(&msg)
}
