// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"strings"
)

type MetaTable struct {
	name         string                // The tablename
	fields       []*MetaField          // The fields in reading order
	fieldForName map[string]*MetaField // Keys are fieldnames
}

func NewMetaTable(tablename string) MetaTable {
	return MetaTable{name: tablename, fields: make([]*MetaField, 0, 1),
		fieldForName: make(map[string]*MetaField)}
}

func (me MetaTable) Name() string {
	return me.name
}

func (me MetaTable) FieldNames() []string {
	result := make([]string, 0, len(me.fields))
	for _, field := range me.fields {
		result = append(result, field.name)
	}
	return result
}

func (me MetaTable) Field(index int) *MetaField {
	return me.fields[index]
}

func (me MetaTable) FieldByName(fieldName string) *MetaField {
	if field, ok := me.fieldForName[fieldName]; ok {
		return field
	}
	return nil
}

func (me *MetaTable) Add(field MetaField) {
	me.fields = append(me.fields, &field)
	me.fieldForName[field.name] = &field
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
	s.WriteString(me.name)
	for _, field := range me.fields {
		s.WriteByte(' ')
		s.WriteString(field.String())
	}
	return s.String()
}
