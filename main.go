package main

import (
	"database/sql"
	"fmt"
	"github.com/fantasticmao/yact/pgwire"
	_ "github.com/marcboeker/go-duckdb/v2"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", ":5432")
	checkError(err)

	db, err := sql.Open("duckdb", "yact.db")
	checkError(err)

	err = db.Ping()
	checkError(err)

	for {
		conn, err := listener.Accept()
		checkError(err)
		fmt.Printf("connection accepted, remote addr: %s\n", conn.RemoteAddr())

		s := pgwire.NewServer(conn, db)
		go func(_conn net.Conn) {
			err := s.Run()
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err.Error())
			}

			_conn.Close()
			fmt.Printf("connection closed, remote addr: %s\n", conn.RemoteAddr())

		}(conn)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
