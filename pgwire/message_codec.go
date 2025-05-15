package pgwire

import (
	"database/sql"
	pg "github.com/jackc/pgx/v5/pgproto3"
	"github.com/jackc/pgx/v5/pgtype"
	"strconv"
)

const (
	rowDescTableOID             = uint32(0)
	rowDescTableAttributeNumber = uint16(0)
	rowDescTypeModifier         = int32(-1)
	rowDescFormat               = int16(0)
)

var pgTypeMap = pgtype.NewMap()

// Message format of RowDescription: https://www.postgresql.org/docs/current/protocol-message-formats.html#PROTOCOL-MESSAGE-FORMATS-ROWDESCRIPTION
func encodeRowDescription(buf []byte, rows *sql.Rows) ([]byte, int, error) {
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, -1, err
	}

	colLen := len(colTypes)
	rowDesc := make([]pg.FieldDescription, colLen)
	for i, colType := range colTypes {
		// FIXME
		//dt, ok := pgTypeMap.TypeForValue(colType.ScanType())
		//if !ok {
		//	return nil, -1, fmt.Errorf("unknown type %T", colType.ScanType())
		//}
		colName := colType.Name()

		rowDesc[i] = pg.FieldDescription{
			Name:                 []byte(colName),
			TableOID:             rowDescTableOID,
			TableAttributeNumber: rowDescTableAttributeNumber,
			DataTypeOID:          25,        // FIXME
			DataTypeSize:         int16(-1), // FIXME
			TypeModifier:         rowDescTypeModifier,
			Format:               rowDescFormat,
		}
	}
	buf, err = (&pg.RowDescription{Fields: rowDesc}).Encode(buf)
	return buf, colLen, err
}

// Message format of DataRow: https://www.postgresql.org/docs/current/protocol-message-formats.html#PROTOCOL-MESSAGE-FORMATS-DATAROW
func encodeDataRow(buf []byte, rows *sql.Rows, colLen int) ([]byte, int, error) {
	rowCnt := 0
	for ; rows.Next(); rowCnt++ {
		data := make([][]byte, colLen)
		for i := range data {
			// FIXME column data type size
			data[i] = make([]byte, 8)
		}

		dataScan := make([]any, colLen)
		for i := range data {
			dataScan[i] = &data[i]
		}

		err := rows.Scan(dataScan...)
		if err != nil {
			return buf, rowCnt, err
		}

		buf, err = (&pg.DataRow{Values: data}).Encode(buf)
		if err != nil {
			return buf, rowCnt, err
		}
	}
	return buf, rowCnt, nil
}

// Message format of CommandComplete: https://www.postgresql.org/docs/current/protocol-message-formats.html#PROTOCOL-MESSAGE-FORMATS-COMMANDCOMPLETE
func encodeCommandComplete(buf []byte, rowCnt int) ([]byte, error) {
	buf, err := (&pg.CommandComplete{CommandTag: []byte("SELECT " + strconv.Itoa(rowCnt))}).Encode(buf)
	if err != nil {
		return buf, err
	}
	return buf, nil
}

// Message format of ReadyForQuery: https://www.postgresql.org/docs/current/protocol-message-formats.html#PROTOCOL-MESSAGE-FORMATS-READYFORQUERY
func encodeReadyForQuery(buf []byte) ([]byte, error) {
	buf, err := (&pg.ReadyForQuery{TxStatus: 'I'}).Encode(buf)
	if err != nil {
		return buf, err
	}
	return buf, nil
}
