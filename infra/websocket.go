package infra

import (
	"context"
	"github.com/coder/websocket"
)

type WebsocketClient struct {
	conn *websocket.Conn
}

func DialWebsocket(ctx context.Context, url string) (*WebsocketClient, error) {
	// TODO support auto reconnect when connection is closed
	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return nil, err
	}

	return &WebsocketClient{
		conn: conn,
	}, nil
}

func (wsClient *WebsocketClient) Run(ctx context.Context, queue chan string) {
	go wsClient.receiveMessage(ctx, queue)
}

func (wsClient *WebsocketClient) receiveMessage(ctx context.Context, queue chan string) {
	for {
		msgType, msg, err := wsClient.conn.Read(ctx)
		if err != nil {
			Error("websocket read message error: %v\n", err)
		}

		if msgType == websocket.MessageText {
			queue <- string(msg)
		} else {
			Error("websocket received non-text message: %d\n", msgType)
		}
	}
}

func (wsClient *WebsocketClient) Close() error {
	return wsClient.conn.CloseNow()
}
