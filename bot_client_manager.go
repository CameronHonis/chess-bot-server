package main

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/server"
)

var botClientManager *BotClientManager

type BotClientManager struct {
	botClientByMatchId map[string]*BotClient
}

func GetBotClientManager() *BotClientManager {
	if botClientManager == nil {
		botClientManager = &BotClientManager{
			botClientByMatchId: make(map[string]*BotClient),
		}
	}
	return botClientManager
}

func (bm *BotClientManager) AddBotClient(match *server.Match, botName string) error {
	// TODO: add lookups for remote bot clients
	botClient, botClientErr := NewLocalBotClient(botName)
	if botClientErr != nil {
		return botClientErr
	}
	initErr := botClient.Initialize(match)
	if initErr != nil {
		return initErr
	}
	bm.botClientByMatchId[match.Uuid] = &botClient
	return nil
}

func (bm *BotClientManager) GetBotClient(matchId string) (*BotClient, error) {
	botClient, ok := bm.botClientByMatchId[matchId]
	if !ok {
		return nil, fmt.Errorf("no bot client found for match %s", matchId)
	}
	return botClient, nil
}

func (bm *BotClientManager) RemoveBotClient(match *server.Match) error {
	botClient, fetchBotClientErr := bm.GetBotClient(match.Uuid)
	if fetchBotClientErr != nil {
		return fetchBotClientErr
	}
	delete(bm.botClientByMatchId, match.Uuid)
	return (*botClient).Terminate(match)
}
