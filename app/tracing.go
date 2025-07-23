package app

import (
	"context"
	"database/sql"
	"github.com/fantasticmao/yact/infra"
)

type Tracing struct {
	db    *sql.DB
	queue chan *ClashEvent
}

func NewTracing(db *sql.DB) *Tracing {
	return &Tracing{
		db:    db,
		queue: make(chan *ClashEvent, 1000),
	}
}

func (t *Tracing) Listen() error {
	// 1. receive a ClashEvent from the channel

	// 2. save the events to the database

	// Create a new WebSocket client to connect to the traffic endpoint
	wsConn, err := infra.NewWebsocketConn("ws://127.0.0.1:9090/traffic")
	if err != nil {
		return err
	}
	err = wsConn.Run(context.Background(), receive)
	return err
}

func (t *Tracing) Save() error {
	return nil
}

func receive(data []byte) {

}
