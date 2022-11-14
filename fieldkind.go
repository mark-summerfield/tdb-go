// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

type FieldKind uint8

const (
	BoolKind = iota
	BytesKind
	DateKind
	DateTimeKind
	IntKind
	RealKind
	StrKind
)

func (me FieldKind) String() string {
	switch me {
	case BoolKind:
		return "bool"
	case BytesKind:
		return "bytes"
	case DateKind:
		return "date"
	case DateTimeKind:
		return "datetime"
	case IntKind:
		return "int"
	case RealKind:
		return "real"
	case StrKind:
		return "str"
	}
	panic("unhandled FieldKind")
}
