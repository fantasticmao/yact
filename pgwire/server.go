package pgwire

import (
	"database/sql"
	"fmt"
	pg "github.com/jackc/pgx/v5/pgproto3"
	"net"
	"reflect"
)

type Server struct {
	conn    net.Conn
	backend *pg.Backend
	db      *sql.DB
}

func NewServer(conn net.Conn, db *sql.DB) *Server {
	initMsgHandlerMap(db)

	backend := pg.NewBackend(conn, conn)
	return &Server{
		conn:    conn,
		backend: backend,
		db:      db,
	}
}

func (p *Server) Close() error {
	return p.conn.Close()
}

func (p *Server) Run() error {
	defer p.Close()

	err := p.handleStartup()
	if err != nil {
		return err
	}

	for {
		msg, err := p.backend.Receive()
		if err != nil {
			return fmt.Errorf("receive message error: %w", err)
		}

		msgType := reflect.TypeOf(msg)
		msgHandler, ok := msgHandlers[msgType]
		if !ok {
			return fmt.Errorf("unsupported message: %#v", msg)
		}

		buf, err := msgHandler.handle(msg)
		if err != nil {
			return fmt.Errorf("handle message failed: %w", err)
		}
		if buf == nil {
			// Terminate
			return nil
		}

		_, err = p.conn.Write(buf)
		if err != nil {
			return fmt.Errorf("write response buffer error: %w", err)
		}
	}
}

func (p *Server) handleStartup() error {
	msg, err := p.backend.ReceiveStartupMessage()
	if err != nil {
		return fmt.Errorf("receive startup message error: %w", err)
	}

	msgType := reflect.TypeOf(msg)
	msgHandler, ok := msgHandlers[msgType]
	if !ok {
		return fmt.Errorf("unknown startup message: %#v", msg)
	}

	buf, err := msgHandler.handle(msg)
	if err != nil || buf == nil {
		return fmt.Errorf("handle startup message failed: %w", err)
	}

	_, err = p.conn.Write(buf)
	if err != nil {
		return fmt.Errorf("write response buffer error: %w", err)
	}

	if _, ok := msg.(*pg.SSLRequest); ok {
		return p.handleStartup()
	} else {
		return nil
	}
}
