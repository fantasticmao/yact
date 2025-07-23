package main

import (
	"database/sql"
	"github.com/fantasticmao/yact/app"
	_ "github.com/marcboeker/go-duckdb/v2"
)

func main() {
	db, err := sql.Open("duckdb", "test.db")
	checkError(err)

	err = db.Ping()
	checkError(err)

	err = app.StartUpDuckPg(":5432", db)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
