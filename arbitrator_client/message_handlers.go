package arbitrator_client

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-bot-server/app"
	"os"
)

func (ac *ArbitratorClient) HandleMsg(msg *models.Message) error {
	switch msg.ContentType {
	case models.CONTENT_TYPE_AUTH:
		return ac.HandleAuthMessage(msg)
	case models.CONTENT_TYPE_MATCH_UPDATE:
		return ac.HandleMatchUpdateMessage(msg)
	case models.CONTENT_TYPE_UPGRADE_AUTH_DENIED:
		return ac.HandleUpgradeAuthDeniedMessage(msg)
	case models.CONTENT_TYPE_UPGRADE_AUTH_GRANTED:
		return ac.HandleAuthUpgradeGrantedMessage(msg)
	case models.CONTENT_TYPE_SUBSCRIBE_REQUEST_DENIED:
		return ac.HandleSubscribeDeniedMessage(msg)
	case models.CONTENT_TYPE_CHALLENGE_PLAYER:
		return ac.HandleChallengePlayerMessage(msg)
	case models.CONTENT_TYPE_SUBSCRIBE_REQUEST_GRANTED:
		return nil
	case models.CONTENT_TYPE_MOVE:
		return nil
	default:
		return fmt.Errorf("unhandled message with content type %s", msg.ContentType)
	}
}

func (ac *ArbitratorClient) HandleAuthMessage(msg *models.Message) error {
	content, ok := msg.Content.(*models.AuthMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to AuthMessageContent")
	}
	ac.SetPublicPrivateKey(content.PublicKey, content.PrivateKey)
	botSecret, ok := os.LookupEnv("BOT_CLIENT_SECRET")
	if !ok {
		panic("could not determine bot client secret")
	}
	authUpgradeErr := ac.RequestAuthUpgrade(botSecret)
	if authUpgradeErr != nil {
		ac.LogService.LogRed(app.ENV_ARBITRATOR_CLIENT, "could not send upgrade auth request: ",
			authUpgradeErr.Error())
	}
	return nil
}

func (ac *ArbitratorClient) HandleAuthUpgradeGrantedMessage(msg *models.Message) error {
	//return GetArbitratorClient().RequestSubscribe("findBotMatch")
	return nil
}

func (ac *ArbitratorClient) HandleUpgradeAuthDeniedMessage(msg *models.Message) error {
	panic("arbitrator denied bot client auth")
}

func (ac *ArbitratorClient) HandleSubscribeDeniedMessage(msg *models.Message) error {
	content, ok := msg.Content.(*models.SubscribeRequestDeniedMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to SubscribeRequestDeniedMessageContent")
	}

	panic(fmt.Sprintf("arbitrator denied bot client subscription%s", content.Topic))
}

func (ac *ArbitratorClient) HandleMatchUpdateMessage(msg *models.Message) error {
	content, ok := msg.Content.(*models.MatchUpdateMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to MatchUpdateMessageContent")
	}

	match := content.Match
	if match.Board.IsTerminal {
		removeErr := ac.BotMngr.RemoveBotClient(match)
		if removeErr != nil {
			ac.LogService.LogRed(app.ENV_ARBITRATOR_CLIENT, "could not remove bot client: ", removeErr.Error())
		}
		return nil
	}

	if match.Board.IsWhiteTurn && match.BlackClientKey == ac.PublicKey() {
		return nil
	} else if !match.Board.IsWhiteTurn && match.WhiteClientKey == ac.PublicKey() {
		return nil
	}

	botClient, botClientErr := ac.BotMngr.BotClientFromMatchId(match.Uuid)
	if botClientErr != nil {
		return botClientErr
	}
	move, moveErr := botClient.Engine().GenerateMove(match)
	if moveErr != nil {
		return moveErr
	}
	return ac.SendMove(match.Uuid, move)
}

func (ac *ArbitratorClient) HandleChallengePlayerMessage(msg *models.Message) error {
	content, ok := msg.Content.(*models.ChallengePlayerMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to ChallengerPlayerMessageContent")
	}

	challenge := content.Challenge
	botInitErr := ac.BotMngr.AddNewBotClient(challenge., challenge.ChallengedKey)
	if botInitErr != nil {
		_ = GetArbitratorClient().FailInitBotRequest(matchId, botName, botInitErr.Error())
		return botInitErr
	}
}
