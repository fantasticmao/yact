package clash

import (
	"database/sql"
	"encoding/json"
	"github.com/fantasticmao/yact/infra"
	_ "github.com/marcboeker/go-duckdb/v2"
	"sync"
)

var database *sql.DB

var databaseInitOnce = sync.Once{}

func initDatabase(db *sql.DB) {
	databaseInitOnce.Do(func() {
		database = db
	})
}

func insertEventTraffic(traffic *EventTraffic) error {
	insertSql := "INSERT INTO clash_traffic(up, down) VALUES (?, ?)"

	result, err := database.Exec(insertSql, traffic.Up, traffic.Down)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	infra.Info("insert to clash_traffic, affected rows: %d\n", rowsAffected)
	return err
}

func insertEventRuleMatch(ruleMatch *EventRuleMatch) error {
	b, err := json.Marshal(ruleMatch.Metadata)
	if err != nil {
		return err
	}
	metadata := string(b)
	insertSql := `
INSERT INTO clash_rule_match(id, type, duration, error, proxy, rule, payload, metadata)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`

	result, err := database.Exec(insertSql, ruleMatch.ID, ruleMatch.Type, ruleMatch.Duration,
		ruleMatch.Error, ruleMatch.Proxy, ruleMatch.Rule, ruleMatch.Payload, metadata)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	infra.Info("insert to clash_rule_match, affected rows: %d\n", rowsAffected)
	return err
}

func insertEventProxyDial(proxyDial *EventProxyDial) error {
	b, err := json.Marshal(proxyDial.Chain)
	if err != nil {
		return err
	}
	chain := string(b)
	insertSql := `
INSERT INTO clash_proxy_dial(id, type, duration, error, proxy, chain, address, host)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`

	result, err := database.Exec(insertSql, proxyDial.ID, proxyDial.Type, proxyDial.Duration,
		proxyDial.Error, proxyDial.Proxy, chain, proxyDial.Address, proxyDial.Host)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	infra.Info("insert to clash_proxy_dial, affected rows: %d\n", rowsAffected)
	return nil
}

func insertEventDnsRequest(dnsRequest *EventDnsRequest) error {
	b, err := json.Marshal(dnsRequest.Answer)
	if err != nil {
		return err
	}
	answer := string(b)
	insertSql := `
INSERT INTO clash_dns_request(id, type, duration, error, dnsType, name, qType, answer, source)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`

	result, err := database.Exec(insertSql, dnsRequest.ID, dnsRequest.Type, dnsRequest.Duration,
		dnsRequest.Error, dnsRequest.DndType, dnsRequest.Name, dnsRequest.QType, answer, dnsRequest.Source)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	infra.Info("insert to clash_dns_request, affected rows: %d\n", rowsAffected)
	return nil
}
