// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"strings"
)

func BoolField(name string) MetaField {
	return MetaField{Name: name, Kind: BoolKind}
}

func BytesField(name string) MetaField {
	return MetaField{Name: name, Kind: BytesKind}
}

func DateField(name string) MetaField {
	return MetaField{Name: name, Kind: DateKind}
}

func DateTimeField(name string) MetaField {
	return MetaField{Name: name, Kind: DateTimeKind}
}

func IntField(name string) MetaField {
	return MetaField{Name: name, Kind: IntKind}
}

func RealField(name string) MetaField {
	return MetaField{Name: name, Kind: RealKind}
}

func StrField(name string) MetaField {
	return MetaField{Name: name, Kind: StrKind}
}

type MetaField struct {
	Name string
	Kind FieldKind
}

func (me MetaField) String() string {
	var s strings.Builder
	s.WriteString(me.Name)
	s.WriteByte(' ')
	s.WriteString(me.Kind.String())
	return s.String()
}
