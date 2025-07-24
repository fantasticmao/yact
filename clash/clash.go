package clash

import (
	"context"
	"database/sql"
	"github.com/fantasticmao/yact/infra"
)

func StartupTracing(db *sql.DB, urls map[EventType]string) []error {
	errCh := make(chan error, len(urls))
	for typo, url := range urls {
		go func() {
			err := startTracing(db, typo, url)
			errCh <- err
		}()
	}

	errs := make([]error, len(urls))
	for i := 0; i < len(urls); i++ {
		errs[i] = <-errCh
	}
	return errs
}

func startTracing(db *sql.DB, typo EventType, url string) error {
	wsClient, err := infra.DialWebsocket(context.Background(), url)
	if err != nil {
		return err
	}

	msgQueue := make(chan string, 1024)
	wsClient.Run(context.Background(), msgQueue)

	go func() {
		for msg := range msgQueue {
			err := handleMessage(db, typo, msg)
			if err != nil {
				infra.Error("handleMessage error: %v\n", err)
			}
		}
	}()
	return nil
}

func handleMessage(db *sql.DB, typo EventType, msg string) error {
	infra.Info("type: %v, msg: %s", typo, msg)
	return nil
}
