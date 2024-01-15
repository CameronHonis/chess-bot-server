package arbitrator_client

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/models"
	mods "github.com/CameronHonis/chess-bot-server/models"
)

type Sender func(msg *models.Message) error

func RequestAuthUpgrade(send Sender, role models.RoleName, secret string) error {
	msg := &models.Message{
		Topic:       "",
		ContentType: models.CONTENT_TYPE_UPGRADE_AUTH_REQUEST,
		Content: &models.UpgradeAuthRequestMessageContent{
			Role:   role,
			Secret: secret,
		},
	}
	return send(msg)
}

func RequestSubscribe(send Sender, topic models.MessageTopic) error {
	msg := &models.Message{
		Topic:       "subscribe",
		ContentType: models.CONTENT_TYPE_SUBSCRIBE_REQUEST,
		Content: &models.SubscribeRequestMessageContent{
			Topic: topic,
		},
	}
	return send(msg)
}

func DeclineChallengeRequest(send Sender, topic models.MessageTopic, challengerKey mods.PlrClientKey) error {
	msg := &models.Message{
		Topic:       topic,
		ContentType: models.CONTENT_TYPE_DECLINE_CHALLENGE,
		Content: &models.DeclineChallengeMessageContent{
			ChallengerClientKey: challengerKey,
		},
	}
	return send(msg)
}

func AcceptChallengeRequest(send Sender, topic models.MessageTopic, challengerKey mods.PlrClientKey) error {
	msg := &models.Message{
		Topic:       topic,
		ContentType: models.CONTENT_TYPE_ACCEPT_CHALLENGE,
		Content: &models.AcceptChallengeMessageContent{
			ChallengerClientKey: challengerKey,
		},
	}
	return send(msg)
}

func SendMove(send Sender, matchId string, move *chess.Move) error {
	msg := &models.Message{
		Topic:       models.MessageTopic(fmt.Sprintf("match-%s", matchId)),
		ContentType: models.CONTENT_TYPE_MOVE,
		Content: &models.MoveMessageContent{
			MatchId: matchId,
			Move:    move,
		},
	}
	return send(msg)
}
