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

	tracingUrls := map[clash.EventMode]string{
		clash.EventModeTraffic: "ws://192.168.0.1:9090/traffic",
		clash.EventModeTracing: "ws://192.168.0.1:9090/profile/tracing",
	}
	errs := clash.Startup(db, tracingUrls)
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
