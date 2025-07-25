package clash

import (
	"database/sql"
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

func saveEventTraffic(traffic *Traffic) error {
	query := "INSERT INTO clash_traffic (up, down) VALUES (?, ?)"
	result, err := database.Exec(query, traffic.Up, traffic.Down)

	rowsAffected, err := result.RowsAffected()
	infra.Info("rows affected: %d\n", rowsAffected)
	return err
}
