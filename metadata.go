// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

type MetaData struct {
	tables       []*MetaTable          // The tables in reading order
	tableForName map[string]*MetaTable // Keys are tablenames
}

func NewMetaData() MetaData {
	return MetaData{tables: make([]*MetaTable, 0, 1),
		tableForName: make(map[string]*MetaTable)}
}

func (me MetaData) TableNames() []string {
	result := make([]string, 0, len(me.tables))
	for _, table := range me.tables {
		result = append(result, table.name)
	}
	return result
}

func (me MetaData) Table(index int) *MetaTable {
	return me.tables[index]
}

func (me MetaData) TableByName(tableName string) *MetaTable {
	if table, ok := me.tableForName[tableName]; ok {
		return table
	}
	return nil
}

func (me *MetaData) Add(table MetaTable) {
	me.tables = append(me.tables, &table)
	me.tableForName[table.name] = &table
}
