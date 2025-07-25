package clash

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/fantasticmao/yact/infra"
)

func StartupTracing(duckdb *sql.DB, urls map[EventType]string) []error {
	initDatabase(duckdb)

	errCh := make(chan error, len(urls))
	for typo, url := range urls {
		go func() {
			err := startTracing(typo, url)
			errCh <- err
		}()
	}

	errs := make([]error, len(urls))
	for i := 0; i < len(urls); i++ {
		errs[i] = <-errCh
	}
	return errs
}

func startTracing(typo EventType, url string) error {
	wsClient, err := infra.DialWebsocket(context.Background(), url)
	if err != nil {
		return err
	}

	msgQueue := make(chan []byte, 1024)
	wsClient.Run(context.Background(), msgQueue)

	go func() {
		for msg := range msgQueue {
			err := handleMessage(typo, msg)
			if err != nil {
				infra.Error("handleMessage error: %v\n", err)
			}
		}
	}()
	return nil
}

func handleMessage(typo EventType, msg []byte) error {
	switch typo {
	case EventTraffic:
		return handleTrafficMessage(msg)
	case EventTracing:
		return handleTracingMessage(msg)
	default:
		return nil
	}
}

func handleTrafficMessage(msg []byte) error {
	traffic := &Traffic{}
	if err := json.Unmarshal(msg, traffic); err != nil {
		return err
	}
	//infra.Info("type: traffic, Up: %d, Down: %d", traffic.Up, traffic.Down)
	return saveEventTraffic(traffic)
}

func handleTracingMessage(msg []byte) error {
	infra.Info("type: tracing, msg: %s", string(msg))
	return nil
}
