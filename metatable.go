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

func (me MetaTable) Field(index int) *MetaField {
	return me.fields[index]
}

func (me MetaTable) FieldByName(fieldName string) *MetaField {
	if m, ok := me.fieldForName[fieldName]; ok {
		return m
	}
	return nil
}

func (me *MetaTable) Add(field MetaField) {
	me.fields = append(me.fields, &field)
	me.fieldForName[field.name] = &field
}

func (me MetaTable) String() string {
	var s strings.Builder
	s.WriteString(me.name)
	s.WriteByte('\n')
	for _, field := range me.fields {
		s.WriteByte(' ')
		s.WriteString(field.String())
		s.WriteByte('\n')
	}
	return s.String()
}

func (me MetaTable) warnings() []string {
	// TODO iterate over all fields and call their warnings()
	// TODO and for any which are unique, check for uniqueness
	// TODO and for any which are not null, check for not nullness
	return nil
}
