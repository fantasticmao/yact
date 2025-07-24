package infra

import (
	"context"
	"github.com/zeebo/assert"
	"os"
	"os/signal"
	"testing"
)

func TestAutoReConnClient_Run(t *testing.T) {
	wsClient, err := DialWebsocket(context.Background(), "ws://127.0.0.1:9090/traffic")
	assert.Nil(t, err)

	queue := make(chan string, 100)
	wsClient.Run(context.Background(), queue)

	go func() {
		for msg := range queue {
			Info("Received message: %s", msg)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	<-stop
}
