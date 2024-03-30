package bot_manager

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/builders"
	arb_mods "github.com/CameronHonis/chess-arbitrator/models"
	mods "github.com/CameronHonis/chess-bot-server/models"
	"github.com/CameronHonis/log"
	. "github.com/CameronHonis/marker"
	"github.com/CameronHonis/service"
	"sync"
)

const ENV_BOT_MANAGER = "BOT_MANAGER"

// BotManager A key is generated and assigned to a new bot. This key is mapped to the player
// key that is opposing the bot in the match. This allows a unique bot to be assigned
// for the duration of a match.
type BotManager struct {
	service.Service
	__dependencies__ Marker
	LogService       log.LoggerServiceI

	__state__         Marker
	clientKeyByOppKey map[mods.PlrClientKey]mods.BotClientKey
	oppKeyByClientKey map[mods.BotClientKey]mods.PlrClientKey
	clientByKey       map[mods.BotClientKey]*BotClient
	mu                sync.Mutex
}

func NewBotManager(config *BotManagerConfig) *BotManager {
	m := &BotManager{
		clientKeyByOppKey: make(map[mods.PlrClientKey]mods.BotClientKey),
		oppKeyByClientKey: make(map[mods.BotClientKey]mods.PlrClientKey),
		clientByKey:       make(map[mods.BotClientKey]*BotClient),
		mu:                sync.Mutex{},
	}
	m.Service = *service.NewService(m, config)
	return m
}

func (bm *BotManager) Client(key mods.BotClientKey) (*BotClient, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	if botClient, ok := bm.clientByKey[key]; ok {
		return botClient, nil
	}
	return nil, fmt.Errorf("no client found with key %s", key)
}

func (bm *BotManager) ClientByOppKey(oppKey mods.PlrClientKey) (*BotClient, error) {
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

func (bm *BotManager) ClientByMatch(match *arb_mods.Match) (*BotClient, error) {
	var client *BotClient
	if client, _ = bm.ClientByOppKey(match.WhiteClientKey); client != nil {
		return client, nil
	}
	if client, _ = bm.ClientByOppKey(match.BlackClientKey); client != nil {
		return client, nil
	}
	return nil, fmt.Errorf("no client could be resolved from match")
}

func (bm *BotManager) InitBot(challenge *arb_mods.Challenge) (*BotClient, error) {
	botClient, botClientErr := NewLocalBotClient(challenge.BotName)
	if botClientErr != nil {
		return nil, fmt.Errorf("could not create bot: %s", botClientErr)
	}

	match := builders.NewMatchBuilder().FromChallenge(challenge).Build()
	initErr := botClient.Engine().Initialize(match)
	if initErr != nil {
		return nil, fmt.Errorf("could not init bot: %s", initErr)
	}

	botKey := botClient.Key()
	bm.addBot(challenge.ChallengerKey, botKey, botClient)
	return botClient, nil
}

func (bm *BotManager) RemoveBot(key mods.BotClientKey) error {
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

func (bm *BotManager) addBot(playerKey mods.PlrClientKey, botKey mods.BotClientKey, botClient *BotClient) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.clientByKey[botKey] = botClient
	bm.clientKeyByOppKey[playerKey] = botKey
	bm.oppKeyByClientKey[botKey] = playerKey
}
