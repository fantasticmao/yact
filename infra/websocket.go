package infra

import (
	"context"
	"github.com/coder/websocket"
)

type WebsocketConn struct {
	url string
}

func NewWebsocketConn(url string) (*WebsocketConn, error) {
	return &WebsocketConn{
		url: url,
	}, nil
}

func (wsConn *WebsocketConn) Run(ctx context.Context, callback func([]byte)) error {
	conn, _, err := websocket.Dial(ctx, wsConn.url, nil)
	if err != nil {
		return err
	}
	defer conn.CloseNow()

	for {
		msgType, data, err := conn.Read(ctx)
		if err != nil {
			return err
		}

		if msgType == websocket.MessageBinary {
			continue
		}

		callback(data)
	}
}
