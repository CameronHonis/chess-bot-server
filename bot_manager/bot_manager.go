package bot_manager

import (
	"fmt"
	mainMods "github.com/CameronHonis/chess-arbitrator/models"
	mods "github.com/CameronHonis/chess-bot-server/models"
	"github.com/CameronHonis/log"
	. "github.com/CameronHonis/marker"
	"github.com/CameronHonis/service"
	"sync"
)

const ENV_BOT_MANAGER = "BOT_MANAGER"

type BotManagerI interface {
	service.ServiceI
	Client(key mods.BotClientKey) (BotClient, error)
	ClientByOppKey(oppKey mods.PlrClientKey) (BotClient, error)
	ClientByMatch(match *mainMods.Match) (BotClient, error)

	InitBotClient(botName string, oppKey mods.PlrClientKey) (BotClient, error)
	RemoveClient(key mainMods.Key) error
}

type BotManager struct {
	service.Service
	__dependencies__ Marker
	LogService       log.LoggerServiceI

	__state__         Marker
	clientKeyByOppKey map[mods.PlrClientKey]mods.BotClientKey
	oppKeyByClientKey map[mods.BotClientKey]mods.PlrClientKey
	clientByKey       map[mods.BotClientKey]BotClient
	mu                sync.Mutex
}

func NewBotManager(config *BotManagerConfig) *BotManager {
	m := &BotManager{
		clientKeyByOppKey: make(map[mods.PlrClientKey]mods.BotClientKey),
		oppKeyByClientKey: make(map[mods.BotClientKey]mods.PlrClientKey),
		clientByKey:       make(map[mods.BotClientKey]BotClient),
		mu:                sync.Mutex{},
	}
	m.Service = *service.NewService(m, config)
	return m
}

func (bm *BotManager) Client(key mods.BotClientKey) (BotClient, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	if botClient, ok := bm.clientByKey[key]; ok {
		return botClient, nil
	}
	return nil, fmt.Errorf("no client found with key %s", key)
}

func (bm *BotManager) ClientByOppKey(oppKey mods.PlrClientKey) (BotClient, error) {
	clientKey, clientKeyErr := bm.ClientKeyByOppKey(oppKey)
	if clientKeyErr != nil {
		return nil, clientKeyErr
	}

	bm.mu.Lock()
	defer bm.mu.Unlock()
	if botClient, ok := bm.clientByKey[clientKey]; ok {
		return botClient, nil
	}
	return nil, fmt.Errorf("no bot client found by oppKey %s", oppKey)
}

func (bm *BotManager) ClientByMatch(match *mainMods.Match) (BotClient, error) {
	var client BotClient
	if client, _ = bm.ClientByOppKey(match.WhiteClientKey); client != nil {
		return client, nil
	}
	if client, _ = bm.ClientByOppKey(match.BlackClientKey); client != nil {
		return client, nil
	}
	return nil, fmt.Errorf("no client could be resolved from match")
}

func (bm *BotManager) InitBotClient(botName string, oppKey mods.PlrClientKey) (BotClient, error) {
	// TODO: add lookups for remote bot clients
	botClient, botClientErr := NewLocalBotClient(botName)
	if botClientErr != nil {
		return nil, botClientErr
	}

	clientKey := botClient.Key()
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.clientByKey[clientKey] = botClient
	bm.clientKeyByOppKey[oppKey] = clientKey
	bm.oppKeyByClientKey[clientKey] = oppKey
	return botClient, nil
}

func (bm *BotManager) RemoveClient(key mods.BotClientKey) error {
	client, clientErr := bm.Client(key)
	if clientErr != nil {
		return clientErr
	}

	oppKey, oppKeyErr := bm.OppKeyByClientKey(key)
	if oppKeyErr != nil {
		return oppKeyErr
	}

	bm.mu.Lock()
	defer bm.mu.Unlock()

	delete(bm.clientByKey, key)
	delete(bm.oppKeyByClientKey, key)
	delete(bm.clientKeyByOppKey, oppKey)

	client.Engine().Terminate()
	return nil
}

func (bm *BotManager) ClientKeyByOppKey(oppKey mods.PlrClientKey) (mods.BotClientKey, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	if clientKey, ok := bm.clientKeyByOppKey[oppKey]; ok {
		return clientKey, nil
	}
	return "", fmt.Errorf("no client key exists for opp key %s", oppKey)
}

func (bm *BotManager) OppKeyByClientKey(key mods.BotClientKey) (mods.PlrClientKey, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	if oppKey, ok := bm.oppKeyByClientKey[key]; ok {
		return oppKey, nil
	}
	return "", fmt.Errorf("no opp key exists for client key %s", key)
}
