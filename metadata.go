// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"strings"
)

type MetaData struct {
	tables       []*MetaTable          // The tables are in reading order
	tableForName map[string]*MetaTable // Keys are tablenames
}

func NewMetaData() MetaData {
	return MetaData{tables: make([]*MetaTable, 0, 1),
		tableForName: make(map[string]*MetaTable)}
}

func (me *MetaData) Add(table MetaTable) {
	me.tables = append(me.tables, &table)
	me.tableForName[table.name] = &table
}

func (me MetaData) String() string {
	var s strings.Builder
	for _, table := range me.tables {
		s.WriteByte('[')
		s.WriteString(table.String())
	}
	s.WriteString("%\n")
	return s.String()
}
