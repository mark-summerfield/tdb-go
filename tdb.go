// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import _ "embed"

//go:embed Version.dat
var Version string // This tdb package's version.

type Tdb struct {
	TableNames []string          // order of reading & writing from/to file
	Tables     map[string]*Table // key is tablename
}

func NewTdb() Tdb {
	return Tdb{make([]string, 0), make(map[string]*Table)}
}

func (me *Tdb) AddTable(table *Table) {
	me.TableNames = append(me.TableNames, table.Name)
	me.Tables[table.Name] = table
}

type Table struct {
	MetaTableType // table name and field names and kinds
	Records       []Record
}

func NewTable() Table {
	return Table{MetaTableType{Fields: make([]*MetaFieldType, 0)},
		make([]Record, 0)}
}

type Record []any

func newRecord(columns int) Record {
	return make([]any, columns)
}
