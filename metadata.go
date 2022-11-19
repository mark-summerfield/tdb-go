// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"strings"
)

type metaDataType map[string]*metaTableType // key is tableName

type metaTableType struct {
	name         string // tableName
	fields       []*metaFieldType
	fieldForName map[string]*metaFieldType
}

func newMetaTable(name string) *metaTableType {
	return &metaTableType{name, make([]*metaFieldType, 0, 1),
		make(map[string]*metaFieldType)}
}

func (me metaTableType) String() string {
	var s strings.Builder
	s.WriteByte('[')
	s.WriteString(me.name)
	for _, field := range me.fields {
		s.WriteByte(' ')
		s.WriteString(field.name)
		s.WriteByte(' ')
		s.WriteString(field.kind.String())
	}
	s.WriteString("%]")
	return s.String()
}

func (me metaTableType) Len() int {
	return len(me.fields)
}

func (me *metaTableType) Add(name, typename string) bool {
	kind, ok := newFieldKind(typename)
	if ok {
		metaField := metaFieldType{name, kind}
		me.fields = append(me.fields, &metaField)
		me.fieldForName[name] = &metaField
	}
	return ok
}

func (me *metaTableType) Field(index int) *metaFieldType {
	return me.fields[index]
}

func (me *metaTableType) FieldByName(fieldName string) *metaFieldType {
	return me.fieldForName[fieldName]
}

type metaFieldType struct {
	name string // fieldName
	kind fieldKind
}

type fieldKind uint8

const (
	boolField fieldKind = 1 << iota
	bytesField
	dateField
	dateTimeField
	intField
	realField
	strField
)

func newFieldKind(typename string) (fieldKind, bool) {
	switch typename {
	case "bool":
		return boolField, true
	case "bytes":
		return bytesField, true
	case "date":
		return dateField, true
	case "datetime":
		return dateTimeField, true
	case "int":
		return intField, true
	case "real":
		return realField, true
	case "str":
		return strField, true
	}
	return boolField, false
}

func (me fieldKind) String() string {
	switch me {
	case boolField:
		return "bool"
	case bytesField:
		return "bytes"
	case dateField:
		return "date"
	case dateTimeField:
		return "datetime"
	case intField:
		return "int"
	case realField:
		return "real"
	case strField:
		return "str"
	}
	panic("invalid fieldKind")
}
