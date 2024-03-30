package bot_manager

import (
	"github.com/CameronHonis/chess-arbitrator/auth"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-bot-server/engines"
)

const ENV_BOT_CLIENT = "BOT_CLIENT"

type BotClient struct {
	key    models.Key
	engine engines.Engine
}

func NewLocalBotClient(engineName string) (*BotClient, error) {
	engine, err := engines.EngineFromName(engineName)
	if err != nil {
		return nil, err
	}
	pubKey, _ := auth.GenerateKeyset()
	botClient := &BotClient{
		key:    pubKey,
		engine: engine,
	}

	return botClient, nil
}

func (c *BotClient) Key() models.Key {
	return c.key
}

func (c *BotClient) Engine() engines.Engine {
	return c.engine
}
