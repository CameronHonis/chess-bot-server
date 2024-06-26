package arbitrator_client

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-bot-server/bot_manager"
	"github.com/CameronHonis/log"
	. "github.com/CameronHonis/marker"
	"github.com/CameronHonis/service"
	"github.com/gorilla/websocket"
	"time"
)

const ENV_ARBITRATOR_CLIENT = "ARBITRATOR_CLIENT"

type ArbitratorClientI interface {
	service.ServiceI
	PublicKey() models.Key
	SendMessage(msg *models.Message) error
}

type ArbitratorClient struct {
	service.Service

	__dependencies__ Marker
	LogService       log.LoggerServiceI
	BotMngr          *bot_manager.BotManager

	__state__ Marker
	conn      *websocket.Conn
	pubKey    models.Key
	priKey    models.Key
}

func NewArbitratorClient(config *ArbitratorClientConfig) *ArbitratorClient {
	s := &ArbitratorClient{}
	s.Service = *service.NewService(s, config)
	return s
}

func (ac *ArbitratorClient) OnBuild() {
	ac.AddEventListener(CONN_SUCCESS, OnConnSuccess)
}

func (ac *ArbitratorClient) OnStart() {
	for {
		ac.Connect()
		ac.ListenOnWebsocket()
	}
}

func (ac *ArbitratorClient) PublicKey() models.Key {
	return ac.pubKey
}

func (ac *ArbitratorClient) SendMessage(msg *models.Message) error {
	if ac.conn == nil {
		return fmt.Errorf("cannot send message, connection is nil")
	}
	ac.SignMessage(msg)
	msgBytes, marshalErr := msg.Marshal()
	if marshalErr != nil {
		return marshalErr
	}
	ac.LogService.Log(ENV_ARBITRATOR_CLIENT, "<< ", string(msgBytes))
	return ac.conn.WriteMessage(websocket.TextMessage, msgBytes)
}

func (ac *ArbitratorClient) Connect() {
	config := ac.Config().(*ArbitratorClientConfig)
	wsUrl := fmt.Sprintf("ws://%s/ws", config.Url())
	for ac.conn == nil {
		var err error
		ac.conn, _, err = websocket.DefaultDialer.Dial(wsUrl, nil)
		if err != nil {
			ac.LogService.Log(ENV_ARBITRATOR_CLIENT, fmt.Sprintf("could not Connect to arbitrator, retrying in 1 second: %s", err))
			ac.conn = nil
			time.Sleep(time.Second)
		}
	}
	ac.LogService.Log(ENV_ARBITRATOR_CLIENT, "successfully connected to arbitrator")
	go ac.Dispatch(NewConnSuccessEvent())
}

func (ac *ArbitratorClient) ListenOnWebsocket() {
	ac.LogService.Log(ENV_ARBITRATOR_CLIENT, "listening on arbitrator websocket connection...")
	for {
		_, rawMsg, readErr := ac.conn.ReadMessage()
		if readErr != nil {
			ac.LogService.LogRed(ENV_ARBITRATOR_CLIENT, fmt.Sprintf("error reading message from websocket: %s", readErr))
			// assume all readErrs are disconnects
			ac.conn = nil
			break
		}
		ac.LogService.Log(ENV_ARBITRATOR_CLIENT, ">> ", string(rawMsg))

		msg, unmarshalErr := models.UnmarshalToMessage(rawMsg)
		if unmarshalErr != nil {
			ac.LogService.LogRed(ENV_ARBITRATOR_CLIENT, fmt.Sprintf("could not unmarshal message: %s", unmarshalErr))
			continue
		}
		handleMsgErr := ac.HandleMsg(msg)
		if handleMsgErr != nil {
			ac.LogService.LogRed(ENV_ARBITRATOR_CLIENT, fmt.Sprintf("could not handle message: %s", handleMsgErr))
			continue
		}
	}
}

func (ac *ArbitratorClient) SetPublicPrivateKey(publicKey models.Key, privateKey models.Key) {
	ac.pubKey = publicKey
	ac.priKey = privateKey
}

func (ac *ArbitratorClient) SignMessage(msg *models.Message) {
	msg.SenderKey = ac.pubKey
	msg.PrivateKey = ac.priKey
}
