// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"strings"
)

type MetaTable struct {
	Name         string                // The tablename
	Fields       []*MetaField          // The fields in reading order
	FieldForName map[string]*MetaField // Keys are fieldnames
}

func NewMetaTable(tablename string) MetaTable {
	return MetaTable{Name: tablename, Fields: make([]*MetaField, 0, 1),
		FieldForName: make(map[string]*MetaField)}
}

func (me MetaTable) FieldNames() []string {
	result := make([]string, 0, len(me.Fields))
	for _, field := range me.Fields {
		result = append(result, field.Name)
	}
	return result
}

func (me MetaTable) Field(index int) *MetaField {
	return me.Fields[index]
}

func (me MetaTable) FieldByName(fieldName string) *MetaField {
	if field, ok := me.FieldForName[fieldName]; ok {
		return field
	}
	return nil
}

func (me *MetaTable) Add(field MetaField) {
	me.Fields = append(me.Fields, &field)
	me.FieldForName[field.Name] = &field
}

func (me *MetaTable) AddBool(names ...string) {
	for _, name := range names {
		me.Add(BoolField(name))
	}
}

func (me *MetaTable) AddBytes(names ...string) {
	for _, name := range names {
		me.Add(BytesField(name))
	}
}

func (me *MetaTable) AddDate(names ...string) {
	for _, name := range names {
		me.Add(DateField(name))
	}
}

func (me *MetaTable) AddDateTime(names ...string) {
	for _, name := range names {
		me.Add(DateTimeField(name))
	}
}

func (me *MetaTable) AddInt(names ...string) {
	for _, name := range names {
		me.Add(IntField(name))
	}
}

func (me *MetaTable) AddReal(names ...string) {
	for _, name := range names {
		me.Add(RealField(name))
	}
}

func (me *MetaTable) AddStr(names ...string) {
	for _, name := range names {
		me.Add(StrField(name))
	}
}

func (me MetaTable) String() string {
	var s strings.Builder
	s.WriteByte('[')
	s.WriteString(me.Name)
	for _, field := range me.Fields {
		s.WriteByte(' ')
		s.WriteString(field.String())
	}
	return s.String()
}
