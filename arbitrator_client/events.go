package arbitrator_client

import "github.com/CameronHonis/service"

const (
	CONN_SUCCESS service.EventVariant = "CONN_SUCCESS"
	CONN_FAILED                       = "CONN_FAILED"
)

type ConnSuccessPayload struct {
}

type ConnSuccessEvent struct{ service.Event }

func NewConnSuccessEvent() *ConnSuccessEvent {
	return &ConnSuccessEvent{
		Event: *service.NewEvent(CONN_SUCCESS, &ConnSuccessPayload{}),
	}
}
