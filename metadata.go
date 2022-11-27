// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import "strings"

type metaDataType map[string]*MetaTableType // key is tableName

func (me metaDataType) addTable(tableName string) *MetaTableType {
	table := MetaTableType{tableName, make([]*MetaFieldType, 0, 1)}
	me[tableName] = &table
	return &table
}

// MetaTableType holds the name of a table and a slice of its fields (names
// and kinds)
type MetaTableType struct {
	Name   string
	Fields []*MetaFieldType
}

func (me MetaTableType) String() string {
	var s strings.Builder
	s.WriteByte('[')
	s.WriteString(me.Name)
	for _, field := range me.Fields {
		s.WriteByte(' ')
		s.WriteString(field.Name)
		s.WriteByte(' ')
		s.WriteString(field.Kind.String())
	}
	s.WriteString("%]")
	return s.String()
}

func (me MetaTableType) Len() int {
	return len(me.Fields)
}

func (me *MetaTableType) AddField(fieldName, typeName string) bool {
	kind, ok := newFieldKind(typeName)
	if ok {
		metaField := MetaFieldType{fieldName, kind}
		me.Fields = append(me.Fields, &metaField)
	}
	return ok
}

func (me *MetaTableType) Field(index int) *MetaFieldType {
	return me.Fields[index]
}

type MetaFieldType struct {
	Name string
	Kind FieldKind
}

type FieldKind uint8

const (
	BoolField FieldKind = 1 << iota
	BytesField
	DateField
	DateTimeField
	IntField
	RealField
	StrField
)

func newFieldKind(typename string) (FieldKind, bool) {
	switch typename {
	case "bool":
		return BoolField, true
	case "bytes":
		return BytesField, true
	case "date":
		return DateField, true
	case "datetime":
		return DateTimeField, true
	case "int":
		return IntField, true
	case "real":
		return RealField, true
	case "str":
		return StrField, true
	}
	return BoolField, false
}

func (me FieldKind) String() string {
	switch me {
	case BoolField:
		return "bool"
	case BytesField:
		return "bytes"
	case DateField:
		return "date"
	case DateTimeField:
		return "datetime"
	case IntField:
		return "int"
	case RealField:
		return "real"
	case StrField:
		return "str"
	}
	panic("invalid FieldKind")
}
