package arbitrator_client

import (
	"fmt"
	mainMods "github.com/CameronHonis/chess-arbitrator/models"
	"os"
)

func (ac *ArbitratorClient) HandleMsg(msg *mainMods.Message) error {
	switch msg.ContentType {
	case mainMods.CONTENT_TYPE_AUTH:
		return ac.HandleAuthMessage(msg)
	case mainMods.CONTENT_TYPE_MATCH_UPDATED:
		return ac.HandleMatchUpdateMessage(msg)
	case mainMods.CONTENT_TYPE_UPGRADE_AUTH_DENIED:
		return ac.HandleUpgradeAuthDeniedMessage(msg)
	case mainMods.CONTENT_TYPE_UPGRADE_AUTH_GRANTED:
		return ac.HandleUpgradeAuthGrantedMessage(msg)
	case mainMods.CONTENT_TYPE_SUBSCRIBE_REQUEST_DENIED:
		return ac.HandleSubscribeDeniedMessage(msg)
	case mainMods.CONTENT_TYPE_CHALLENGE_UPDATED:
		return HandleChallengeUpdatedMessage(ac, msg)
	case mainMods.CONTENT_TYPE_SUBSCRIBE_REQUEST_GRANTED:
		return nil
	case mainMods.CONTENT_TYPE_MOVE:
		return nil
	case mainMods.CONTENT_TYPE_ACCEPT_CHALLENGE:
		return nil
	case mainMods.CONTENT_TYPE_DECLINE_CHALLENGE:
		return nil
	case mainMods.CONTENT_TYPE_REVOKE_CHALLENGE:
		return nil
	default:
		return fmt.Errorf("unhandled message with content type %s", msg.ContentType)
	}
}

func (ac *ArbitratorClient) HandleAuthMessage(msg *mainMods.Message) error {
	content, ok := msg.Content.(*mainMods.AuthMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to AuthMessageContent")
	}
	ac.SetPublicPrivateKey(content.PublicKey, content.PrivateKey)
	botSecret, ok := os.LookupEnv("BOT_CLIENT_SECRET")
	if !ok {
		panic("could not determine bot client secret")
	}
	authUpgradeErr := RequestAuthUpgrade(ac.SendMessage, mainMods.BOT, botSecret)
	if authUpgradeErr != nil {
		ac.LogService.LogRed(ENV_ARBITRATOR_CLIENT, "could not send upgrade auth request: ",
			authUpgradeErr.Error())
	}
	return nil
}

func (ac *ArbitratorClient) HandleUpgradeAuthGrantedMessage(msg *mainMods.Message) error {
	//return GetArbitratorClient().RequestSubscribe("findBotMatch")
	return nil
}

func (ac *ArbitratorClient) HandleUpgradeAuthDeniedMessage(msg *mainMods.Message) error {
	panic("arbitrator denied bot client auth")
}

func (ac *ArbitratorClient) HandleSubscribeDeniedMessage(msg *mainMods.Message) error {
	content, ok := msg.Content.(*mainMods.SubscribeRequestDeniedMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to SubscribeRequestDeniedMessageContent")
	}

	panic(fmt.Sprintf("arbitrator denied bot client subscription%s", content.Topic))
}

func (ac *ArbitratorClient) HandleMatchUpdateMessage(msg *mainMods.Message) error {
	content, ok := msg.Content.(*mainMods.MatchUpdateMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to MatchUpdateMessageContent")
	}

	match := content.Match
	botClient, botClientErr := ac.BotMngr.ClientByMatch(match)
	if botClientErr != nil {
		return botClientErr
	}

	if match.Result != mainMods.MATCH_RESULT_IN_PROGRESS {
		removeErr := ac.BotMngr.RemoveBot(botClient.Key())
		if removeErr != nil {
			return removeErr
		}
	}

	if match.Board.IsWhiteTurn && match.BlackClientKey == ac.PublicKey() {
		return nil
	} else if !match.Board.IsWhiteTurn && match.WhiteClientKey == ac.PublicKey() {
		return nil
	}

	move, moveErr := botClient.Engine().GenerateMove(match)
	if moveErr != nil {
		return moveErr
	}
	return SendMove(ac.SendMessage, match.Uuid, move)
}

var HandleChallengeUpdatedMessage = func(ac *ArbitratorClient, msg *mainMods.Message) error {
	content, ok := msg.Content.(*mainMods.ChallengeUpdatedMessageContent)
	if !ok {
		return fmt.Errorf("could not cast message to ChallengerPlayerMessageContent")
	}

	if content.Challenge == nil {
		return nil
	}

	// NOTE: assumes that only update to challenge is when challenge is created
	challenge := content.Challenge
	_, botInitErr := ac.BotMngr.InitBot(challenge)
	if botInitErr != nil {
		return DeclineChallengeRequest(ac.SendMessage, msg.Topic, challenge.ChallengerKey)
	}

	return AcceptChallengeRequest(ac.SendMessage, msg.Topic, challenge.ChallengerKey)
}
