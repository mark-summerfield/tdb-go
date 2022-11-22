// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"strings"
)

type metaDataType map[string]*metaTableType // key is tableName or tagName

func (me metaDataType) addTable(tableName, tagName string) *metaTableType {
	table := metaTableType{tableName, tagName, make([]*metaFieldType, 0, 1),
		make(map[string]*metaFieldType)}
	me[tableName] = &table
	if tagName != "" {
		me[tagName] = &table
	}
	return &table
}

type metaTableType struct {
	tableName    string
	tagName      string
	fields       []*metaFieldType
	fieldForName map[string]*metaFieldType
}

func (me metaTableType) name() string {
	if me.tagName != "" {
		return me.tagName
	}
	return me.tableName
}

func (me metaTableType) String() string {
	var s strings.Builder
	s.WriteByte('[')
	s.WriteString(me.name())
	for _, field := range me.fields {
		s.WriteByte(' ')
		s.WriteString(field.name())
		s.WriteByte(' ')
		s.WriteString(field.kind.String())
	}
	s.WriteString("%]")
	return s.String()
}

func (me metaTableType) Len() int {
	return len(me.fields)
}

func (me *metaTableType) addField(fieldName, tagName, typeName string) bool {
	kind, ok := newFieldKind(typeName)
	if ok {
		metaField := metaFieldType{fieldName, tagName, kind}
		me.fields = append(me.fields, &metaField)
		me.fieldForName[fieldName] = &metaField
		if tagName != "" {
			me.fieldForName[tagName] = &metaField
		}
	}
	return ok
}

func (me *metaTableType) field(index int) *metaFieldType {
	return me.fields[index]
}

type metaFieldType struct {
	fieldName string
	tagName   string
	kind      fieldKind
}

func (me *metaFieldType) name() string {
	if me.tagName != "" {
		return me.tagName
	}
	return me.fieldName
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
