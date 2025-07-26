package clash

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/fantasticmao/yact/infra"
)

func Startup(duckdb *sql.DB, urls map[EventMode]string) []error {
	initDatabase(duckdb)

	errCh := make(chan error, len(urls))
	for mode, url := range urls {
		go func() {
			err := startTracing(mode, url)
			errCh <- err
		}()
	}

	errs := make([]error, len(urls))
	for i := 0; i < len(urls); i++ {
		errs[i] = <-errCh
	}
	return errs
}

func startTracing(mode EventMode, url string) error {
	wsClient, err := infra.DialWebsocket(context.Background(), url)
	if err != nil {
		return err
	}

	msgQueue := make(chan []byte, 1024)
	wsClient.Run(context.Background(), msgQueue)

	go func() {
		for msg := range msgQueue {
			err := handleMessage(mode, msg)
			if err != nil {
				infra.Error("handleMessage error: %v\n", err)
			}
		}
	}()
	return nil
}

func handleMessage(mode EventMode, msg []byte) error {
	switch mode {
	case EventModeTraffic:
		return handleTrafficMessage(msg)
	case EventModeTracing:
		return handleTracingMessage(msg)
	default:
		return nil
	}
}

func handleTrafficMessage(msg []byte) error {
	traffic := &EventTraffic{}
	if err := json.Unmarshal(msg, traffic); err != nil {
		return err
	}
	//infra.Info("type: traffic, Up: %d, Down: %d", traffic.Up, traffic.Down)
	return insertEventTraffic(traffic)
}

func handleTracingMessage(msg []byte) error {
	basic := &EventBasic{}
	if err := json.Unmarshal(msg, basic); err != nil {
		return err
	}

	switch basic.Type {
	case EventTypeRuleMatch:
		ruleMatch := &EventRuleMatch{}
		if err := json.Unmarshal(msg, ruleMatch); err != nil {
			return err
		}
		return insertEventRuleMatch(ruleMatch)
	case EventTypeProxyDial:
		proxyDial := &EventProxyDial{}
		if err := json.Unmarshal(msg, proxyDial); err != nil {
			return err
		}
		return insertEventProxyDial(proxyDial)
	case EventTypeDNSRequest:
		dnsRequest := &EventDnsRequest{}
		if err := json.Unmarshal(msg, dnsRequest); err != nil {
			return err
		}
		return insertEventDnsRequest(dnsRequest)
	default:
		infra.Error("unknown clash tracing event, type: %s\n", basic.Type)
		return nil
	}
}
