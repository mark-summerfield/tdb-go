// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"bytes"
	"fmt"
)

func Parse(data []byte) (*Tdb, error) {
	db := NewTdb()
	var err error
	var table *Table
	lino := 1
	for len(data) > 0 {
		b := data[0]
		if b == '\n' {
			lino++
			data = data[1:]
		} else if b == '[' {
			data, table, err = readMeta(data[1:], &lino)
			if err != nil {
				return nil, err
			}
			db.AddTable(table)
		} else { // read records into the current table
			data, err = readRecords(data, table, &lino)
			if err != nil {
				return nil, err
			}
			table = nil
		}
	}
	return &db, nil
}

func readMeta(data []byte, lino *int) ([]byte, *Table, error) {
	data, found, err := find(data, '%', "expected to find '%'", lino)
	if err != nil {
		return data, nil, err
	}
	table := NewTable()
	var fieldName string
	for i, part := range bytes.Fields(bytes.TrimSpace(found)) {
		text := string(part)
		if i == 0 {
			table.Name = text
		} else if i%2 != 0 {
			fieldName = text
		} else {
			if !table.AddField(fieldName, text) {
				return data, nil, fmt.Errorf("e%d#%d:invalid typename %q",
					e131, *lino, text)
			}
			fieldName = ""
		}
	}
	return data, &table, nil
}

func find(data []byte, what byte, message string, lino *int) ([]byte,
	[]byte, error) {
	end, err := scanToByte(data, what, lino)
	if err != nil {
		return data, nil, err
	}
	return data[end+1:], data[:end], nil
}

func readRecords(data []byte, table *Table, lino *int) ([]byte, error) {
	var err error
	var record Record = nil
	var kind FieldKind
	oldColumn := -1
	column := 0
	columns := table.Len()
	for len(data) > 0 {
		if record == nil {
			record = newRecord(columns)
			oldColumn = -1
			column = 0
		}
		if column != oldColumn {
			kind = table.Fields[column].Kind
		}
		switch data[0] {
		case '\n': // ignore whitespace
			data = data[1:]
			*lino++
		case ' ', '\t', '\r': // ignore whitespace
			data = data[1:]
		case '!':
			if err = handleSentinal(kind, record, column,
				lino); err != nil {
				return data, err
			}
			data, column = advance(data, column)
		case 'F', 'f', 'N', 'n':
			if err = handleBool(kind, record, column, lino,
				false); err != nil {
				return data, err
			}
			data, column = advance(data, column)
		case 'T', 't', 'Y', 'y':
			if err = handleBool(kind, record, column, lino,
				true); err != nil {
				return data, err
			}
			data, column = advance(data, column)
		case '(':
			data, err = handleBytes(data[1:], kind, record, column, lino)
			if err != nil {
				return data, err
			}
			column++
		case '<':
			data, err = handleStr(data[1:], kind, record, column, lino)
			if err != nil {
				return data, err
			}
			column++
		case '-':
			data, err = handleNegativeNumber(data, record, column, lino,
				kind)
			if err != nil {
				return data, err
			}
			column++
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			data, err = handleBoolNumberDateTime(data, record, column, lino,
				kind)
			if err != nil {
				return data, err
			}
			column++
		case ']': // end of table
			if 0 < column && column < columns {
				return data, fmt.Errorf("e%d#%d:incomplete record %d/%d",
					e134, *lino, column+1, columns)
			}
			return skipWs(data[1:], lino), nil
		default:
			return data, fmt.Errorf("e%d#%d:invalid character %q", e135,
				*lino, data[0])
		}
		if column == columns {
			table.Records = append(table.Records, record)
			record = nil
		}
	}
	return data, nil
}

func advance(data []byte, column int) ([]byte, int) {
	return data[1:], column + 1
}

func handleSentinal(kind FieldKind, record Record, column int,
	lino *int) error {
	switch kind {
	case DateField:
		record[column] = DateSentinal
	case DateTimeField:
		record[column] = DateTimeSentinal
	case IntField:
		record[column] = IntSentinal
	case RealField:
		record[column] = RealSentinal
	default:
		return fmt.Errorf("e%d#%d:%s fields don't have a sentinal", e136,
			*lino, kind)
	}
	return nil
}

func handleBool(kind FieldKind, record Record, column int,
	lino *int, value bool) error {
	if kind != BoolField {
		return fmt.Errorf("e%d#%d:expected %q, got a bool", e137, *lino,
			kind)
	}
	record[column] = value
	return nil
}

func handleBytes(data []byte, kind FieldKind, record Record, column int,
	lino *int) ([]byte, error) {
	if kind != BytesField {
		return data, fmt.Errorf("e%d#%d:expected %q, got a bytes", e138,
			*lino, kind)
	}
	data, raw, err := readHexBytes(data, lino)
	if err != nil {
		return data, err
	}
	record[column] = raw
	return data, nil
}

func handleStr(data []byte, kind FieldKind, record Record, column int,
	lino *int) ([]byte, error) {
	if kind != StrField {
		return data, fmt.Errorf("e%d#%d:expected %q, got a str", e139,
			*lino, kind)
	}
	data, s, err := readString(data, lino)
	if err != nil {
		return data, err
	}
	record[column] = s
	return data, nil
}

func handleNegativeNumber(data []byte, record Record, column int,
	lino *int, kind FieldKind) ([]byte, error) {
	var err error
	switch kind {
	case IntField:
		data, err = handleInt(data, record, column, lino)
	case RealField:
		data, err = handleReal(data, record, column, lino)
	default:
		err = fmt.Errorf("e%d#%d:expected %q", e132, *lino, kind)
	}
	return data, err
}

func handleBoolNumberDateTime(data []byte, record Record, column int,
	lino *int, kind FieldKind) ([]byte, error) {
	var err error
	switch kind {
	case BoolField:
		if (data[0] == '0' || data[0] == '1') && len(data) > 1 &&
			bytes.IndexByte([]byte{'.', 'e', 'E', '0', '1', '2',
				'3', '4', '5', '6', '7', '8', '9'}, data[1]) == -1 {
		} else {
			err = fmt.Errorf("e%d#%d:got %c%c expected a %s", e133,
				*lino, data[0], data[1], kind)
		}
	case DateField:
		data, err = handleDate(data, record, column, lino)
	case DateTimeField:
		data, err = handleDateTime(data, record, column, lino)
	case IntField:
		data, err = handleInt(data, record, column, lino)
	case RealField:
		data, err = handleReal(data, record, column, lino)
	default:
		err = fmt.Errorf("e%d#%d:expected %q", e132, *lino, kind)
	}
	return data, err
}

func handleInt(data []byte, record Record, column int, lino *int) ([]byte,
	error) {
	data, i, err := readInt(data, lino)
	if err != nil {
		return data, err
	}
	record[column] = i
	return data, nil
}

func handleReal(data []byte, record Record, column int, lino *int) ([]byte,
	error) {
	data, r, err := readReal(data, lino)
	if err != nil {
		return data, err
	}
	record[column] = r
	return data, nil
}

func handleDate(data []byte, record Record, column int, lino *int) ([]byte,
	error) {
	data, d, err := readDateTime(data, DateFormat, lino)
	if err != nil {
		return data, err
	}
	record[column] = d
	return data, nil
}

func handleDateTime(data []byte, record Record, column int,
	lino *int) ([]byte,
	error) {
	data, d, err := readDateTime(data, DateTimeFormat, lino)
	if err != nil {
		return data, err
	}
	record[column] = d
	return data, nil
}
