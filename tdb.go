// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
)

//go:embed Version.dat
var Version string // This tdb package's version.

// TODO Use same API as std. lib. csv, e.g., NewReader() & NewWriter()
// equiv. to Py load & dump

type Tdb struct {
	TableNames []string          // order of reading & writing from/to file
	Tables     map[string]*Table // key is tablename
}

func (me *Tdb) AddTable(table *Table) {
	me.TableNames = append(me.TableNames, table.Name)
	me.Tables[table.Name] = table
}

func (me *Tdb) Write(out io.Writer) error {
	// TODO
	return nil
}

type Table struct {
	MetaTableType // table name and field names and kinds
	Records       []Record
}

type Record []any

func NewTdb(data []byte) (*Tdb, error) {
	db := Tdb{}
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
		}
	}
	return &db, nil
}

func readMeta(data []byte, lino *int) ([]byte, *Table, error) {
	found, data, err := find(data, '%', "expected to find '%'", lino)
	if err != nil {
		return data, nil, err
	}
	table := Table{}
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
	end := bytes.IndexByte(data, what)
	if end == -1 {
		return data, nil, fmt.Errorf("e%d#%d:%s", e131, *lino, message)
	}
	*lino += bytes.Count(data[:end], []byte{'\n'})
	return data[:end], data[end+1:], nil
}

func readRecords(data []byte, table *Table, lino *int) ([]byte, error) {
	var err error
	var record Record = nil
	var fieldMeta *MetaFieldType = nil
	oldColumn := -1
	column := 0
	//columns := table.Len()
	for len(data) > 0 {
		if record == nil {
			record = Record{}
			oldColumn = -1
			column = 0
		}
		if column != oldColumn {
			fieldMeta = table.Fields[column]
		}
		switch data[0] {
		case '\n': // ignore whitespace
			data = data[1:]
			*lino++
		case ' ', '\t', '\r': // ignore whitespace
			data = data[1:]
		case '!':
			if err = handleSentinal(fieldMeta, &record, column,
				lino); err != nil {
				return data, err
			}
			data, column = advance(data, column)
		case 'F', 'f', 'N', 'n':
			if err = handleBool(fieldMeta.Kind, &record, column, lino,
				false); err != nil {
				return data, err
			}
			data, column = advance(data, column)
		case 'T', 't', 'Y', 'y':
			if err = handleBool(fieldMeta.Kind, &record, column, lino,
				true); err != nil {
				return data, err
			}
			data, column = advance(data, column)
		case '(':
			data, err = handleBytes(data[1:], fieldMeta.Kind, &record,
				column, lino)
			if err != nil {
				return data, err
			}
			column++
		}
	}
	return data, nil
}

func advance(data []byte, column int) ([]byte, int) {
	return data[1:], column + 1
}

func handleSentinal(fieldMeta *MetaFieldType, record *Record, column int,
	lino *int) error {
	// TODO
	return nil
}

func handleBool(kind FieldKind, record *Record, column int,
	lino *int, value bool) error {
	// TODO
	return nil
}

func handleBytes(data []byte, kind FieldKind, record *Record, column int,
	lino *int) ([]byte, error) {
	// TODO
	return data, nil
}
