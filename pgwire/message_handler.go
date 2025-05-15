package pgwire

import (
	"database/sql"
	pg "github.com/jackc/pgx/v5/pgproto3"
	"reflect"
	"sync"
)

var (
	msgHandlersInitOnce = sync.Once{}
	msgHandlers         = make(map[reflect.Type]messageHandler)
)

func initMsgHandlerMap(db *sql.DB) {
	msgHandlersInitOnce.Do(func() {
		handlers := []messageHandler{
			startupHandler{},
			sslRequestHandler{},
			queryHandler{db: db},
			terminateHandler{},
		}

		for _, handler := range handlers {
			msgHandlers[handler.msgType()] = handler
		}
	})
}

type messageHandler interface {
	// msgType returns the type of the message that this handler can handle.
	msgType() reflect.Type

	// handle handles the message and returns the response.
	handle(msg pg.FrontendMessage) ([]byte, error)
}

type startupHandler struct {
}

func (h startupHandler) msgType() reflect.Type {
	msg := pg.StartupMessage{}
	return reflect.TypeOf(&msg)
}
func (h startupHandler) handle(msg pg.FrontendMessage) ([]byte, error) {
	buf, err := (&pg.AuthenticationOk{}).Encode(nil)
	if err != nil {
		return nil, err
	}

	buf, err = (&pg.ReadyForQuery{TxStatus: 'I'}).Encode(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

type sslRequestHandler struct {
}

func (h sslRequestHandler) msgType() reflect.Type {
	msg := pg.SSLRequest{}
	return reflect.TypeOf(&msg)
}

func (h sslRequestHandler) handle(msg pg.FrontendMessage) ([]byte, error) {
	buf := []byte("N")
	return buf, nil
}

type queryHandler struct {
	db *sql.DB
}

func (h queryHandler) msgType() reflect.Type {
	msg := pg.Query{}
	return reflect.TypeOf(&msg)
}

func (h queryHandler) handle(msg pg.FrontendMessage) ([]byte, error) {
	query := msg.(*pg.Query).String
	rows, err := h.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	buf, colLen, err := encodeRowDescription(nil, rows)
	if err != nil {
		return buf, err
	}

	buf, rowCnt, err := encodeDataRow(buf, rows, colLen)
	if err != nil {
		return buf, err
	}

	buf, err = encodeCommandComplete(buf, rowCnt)
	if err != nil {
		return buf, err
	}

	buf, err = encodeReadyForQuery(buf)
	if err != nil {
		return buf, err
	}

	return buf, nil
}

type terminateHandler struct {
}

func (h terminateHandler) msgType() reflect.Type {
	msg := pg.Terminate{}
	return reflect.TypeOf(&msg)
}

func (h terminateHandler) handle(msg pg.FrontendMessage) ([]byte, error) {
	return nil, nil
}
