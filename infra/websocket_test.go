package infra

import (
	"context"
	"fmt"
	"github.com/zeebo/assert"
	"testing"
)

func TestAutoReConnClient_Run(t *testing.T) {
	client, err := NewWebsocketConn("ws://127.0.0.1:9090/traffic")
	//client, err := NewWebsocketConn("ws://127.0.0.1:9090/profile/app")
	assert.Nil(t, err)

	err = client.Run(context.Background(), func(data []byte) {
		fmt.Println(string(data))
	})
	assert.Nil(t, err)
}
