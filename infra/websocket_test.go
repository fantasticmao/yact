package infra

import (
	"context"
	"github.com/zeebo/assert"
	"os"
	"os/signal"
	"testing"
)

func TestWebsocketClient_Run(t *testing.T) {
	wsClient, err := DialWebsocket(context.Background(), "ws://127.0.0.1:9090/traffic")
	assert.Nil(t, err)

	queue := make(chan []byte, 100)
	wsClient.Run(context.Background(), queue)

	go func() {
		for msg := range queue {
			Info("Received message: %s", string(msg))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	<-stop
}
