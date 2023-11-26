package bot_client

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/server"
	. "github.com/CameronHonis/log"
)

var botClientManager *BotClientManager

type BotClientManager struct {
	botClientByMatchId map[string]BotClient
}

func GetBotClientManager() *BotClientManager {
	if botClientManager == nil {
		botClientManager = &BotClientManager{
			botClientByMatchId: make(map[string]BotClient),
		}
	}
	return botClientManager
}

func (bm *BotClientManager) AddNewBotClient(matchId string, botName string) error {
	GetLogManager().Log(ENV_BOT_CLIENT_MANAGER, fmt.Sprintf("adding bot client for match %s", matchId))
	// TODO: add lookups for remote bot clients
	botClient, botClientErr := NewLocalBotClient(botName)
	if botClientErr != nil {
		return botClientErr
	}
	//initErr := botClient.Initialize(match)
	//if initErr != nil {
	//	return initErr
	//}
	bm.botClientByMatchId[matchId] = botClient
	return nil
}

func (bm *BotClientManager) GetBotClient(matchId string) (BotClient, error) {
	botClient, ok := bm.botClientByMatchId[matchId]
	if !ok {
		return nil, fmt.Errorf("no bot client found for match %s", matchId)
	}
	return botClient, nil
}

func (bm *BotClientManager) RemoveBotClient(match *server.Match) error {
	GetLogManager().Log(ENV_BOT_CLIENT_MANAGER, fmt.Sprintf("removing bot client for match %s", match.Uuid))
	botClient, fetchBotClientErr := bm.GetBotClient(match.Uuid)
	if fetchBotClientErr != nil {
		return fetchBotClientErr
	}
	delete(bm.botClientByMatchId, match.Uuid)
	return botClient.Terminate(match)
}
