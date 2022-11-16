// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"strings"
)

func BoolField(name string) MetaField {
	return MetaField{name: name, kind: BoolKind}
}

func BytesField(name string) MetaField {
	return MetaField{name: name, kind: BytesKind}
}

func DateField(name string) MetaField {
	return MetaField{name: name, kind: DateKind}
}

func DateTimeField(name string) MetaField {
	return MetaField{name: name, kind: DateTimeKind}
}

func IntField(name string) MetaField {
	return MetaField{name: name, kind: IntKind}
}

func RealField(name string) MetaField {
	return MetaField{name: name, kind: RealKind}
}

func StrField(name string) MetaField {
	return MetaField{name: name, kind: StrKind}
}

type MetaField struct {
	name string
	kind FieldKind
}

func (me MetaField) String() string {
	var s strings.Builder
	s.WriteString(me.name)
	s.WriteByte(' ')
	s.WriteString(me.kind.String())
	return s.String()
}

func (me MetaField) Kind() FieldKind {
	return me.kind
}
