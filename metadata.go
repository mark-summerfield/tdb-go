// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

type MetaData struct {
	Tables       []*MetaTable          // The tables in reading order
	TableForName map[string]*MetaTable // Keys are tablenames
}

func NewMetaData() MetaData {
	return MetaData{Tables: make([]*MetaTable, 0, 1),
		TableForName: make(map[string]*MetaTable)}
}

func (me MetaData) TableNames() []string {
	result := make([]string, 0, len(me.Tables))
	for _, table := range me.Tables {
		result = append(result, table.Name)
	}
	return result
}

func (me MetaData) Table(index int) *MetaTable {
	return me.Tables[index]
}

func (me MetaData) TableByName(tableName string) *MetaTable {
	if table, ok := me.TableForName[tableName]; ok {
		return table
	}
	return nil
}

func (me *MetaData) Add(table MetaTable) {
	me.Tables = append(me.Tables, &table)
	me.TableForName[table.Name] = &table
}
