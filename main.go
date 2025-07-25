package main

import (
	"database/sql"
	"github.com/fantasticmao/go-duckpg/duckpg"
	"github.com/fantasticmao/yact/clash"
)

func main() {
	db, err := sql.Open("duckdb", "yact.db")
	checkError(err)

	err = db.Ping()
	checkError(err)

	clashUrls := map[clash.EventType]string{
		clash.EventTraffic: "ws://127.0.0.1:9090/traffic",
		clash.EventTracing: "ws://127.0.0.1:9090/profile/tracing",
	}
	errs := clash.StartupTracing(db, clashUrls)
	checkError(errs...)

	err = duckpg.Startup(db, ":5432")
	checkError(err)
}

func checkError(errs ...error) {
	for _, err := range errs {
		if err != nil {
			panic(err)
		}
	}
}
