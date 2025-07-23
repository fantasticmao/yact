package app

import (
	"database/sql"
	"github.com/fantasticmao/go-duckpg/duckpg"
)

func StartUpDuckPg(address string, db *sql.DB) error {
	server := duckpg.NewServer(address, db)
	return server.StartUp()
}
