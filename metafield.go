// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"fmt"
	"github.com/mark-summerfield/gset"
	"strings"
	"time"
)

func BoolField(name string) MetaField {
	return MetaField{name: name, kind: BoolKind, flag: notNullFlag}
}

func BytesField(name string) MetaField {
	return MetaField{name: name, kind: BytesKind, flag: notNullFlag}
}

func DateField(name string) MetaField {
	return MetaField{name: name, kind: DateKind, flag: notNullFlag}
}

func DateTimeField(name string) MetaField {
	return MetaField{name: name, kind: DateTimeKind, flag: notNullFlag}
}

func IntField(name string) MetaField {
	return MetaField{name: name, kind: IntKind, flag: notNullFlag}
}

func RealField(name string) MetaField {
	return MetaField{name: name, kind: RealKind, flag: notNullFlag}
}

func StrField(name string) MetaField {
	return MetaField{name: name, kind: StrKind, flag: notNullFlag}
}

type MetaField struct {
	name       string
	kind       FieldKind
	flag       fieldFlag
	min        any // length for bytes & str; min value otherwise
	max        any // length for bytes & str; max value otherwise
	inInts     gset.Set[int]
	inStrs     gset.Set[string]
	theDefault any
	ref        string // tablename.fieldname or fieldname
}

func (me *MetaField) SetNullable() {
	me.flag = me.flag.with(nullableFlag)
}

// SetUnique sets the field to be unique: this also implies it is not
// nullable.
func (me *MetaField) SetUnique() {
	me.flag = uniqueFlag
}

func (me *MetaField) SetDefault(theDefault any) error {
	switch me.kind {
	case BoolKind:
		if _, ok := theDefault.(bool); ok {
			me.theDefault = theDefault
			return nil
		}
	case BytesKind:
		if _, ok := theDefault.([]byte); ok {
			me.theDefault = theDefault
			return nil
		}
	case DateKind, DateTimeKind:
		if _, ok := theDefault.(time.Time); ok {
			me.theDefault = theDefault
			return nil
		}
	case IntKind:
		if _, ok := theDefault.(int); ok {
			me.theDefault = theDefault
			return nil
		}
	case RealKind:
		if _, ok := theDefault.(float64); ok {
			me.theDefault = theDefault
			return nil
		}
	case StrKind:
		if _, ok := theDefault.(string); ok {
			me.theDefault = theDefault
			return nil
		}
	}
	return fmt.Errorf(
		"error#%d:%v is not a valid default for a field of type %s",
		eInvalidDefault, theDefault, me.kind)
}

func (me *MetaField) SetMin(min any) error {
	if err := validMinOrMax(me.kind, min, "min"); err != nil {
		return err
	}
	me.min = min
	return nil
}

func (me *MetaField) SetMax(max any) error {
	if err := validMinOrMax(me.kind, max, "max"); err != nil {
		return err
	}
	me.max = max
	return nil
}

func (me *MetaField) SetInInts(ints ...int) {
	if me.inInts == nil {
		me.inInts = gset.New[int]()
	} else {
		me.inInts.Clear()
	}
	me.inInts.Add(ints...)
}

func (me *MetaField) SetInStrs(ints ...string) {
	if me.inStrs == nil {
		me.inStrs = gset.New[string]()
	} else {
		me.inStrs.Clear()
	}
	me.inStrs.Add(ints...)
}

func (me *MetaField) SetRef(ref string) error {
	if !strings.Contains(ref, ".") && ref == me.name {
		return fmt.Errorf("error#%d:cannot set a field ref to itself (%s)",
			eInvalidRef, ref)
	}
	// TODO check identifer.identifer or .identifer
	me.ref = ref
	return nil
}

func (me MetaField) String() string {
	var s strings.Builder
	s.WriteString(me.name)
	s.WriteByte(' ')
	s.WriteString(me.kind.String())
	if me.IsNullable() {
		s.WriteString("?")
	}
	if me.theDefault != nil {
		fmt.Fprintf(&s, " default %v", me.theDefault)
	}
	if me.min != nil {
		fmt.Fprintf(&s, " min %v", me.min)
	}
	if me.max != nil {
		fmt.Fprintf(&s, " max %v", me.max)
	}
	if me.IsUnique() {
		s.WriteString(" unique")
	}
	if me.inInts != nil && len(me.inInts) > 0 {
		s.WriteString(" in")
		for _, x := range me.inInts.ToSortedSlice() {
			fmt.Fprintf(&s, " %d", x)
		}
	}
	if me.inStrs != nil && len(me.inStrs) > 0 {
		s.WriteString(" in")
		for _, x := range me.inStrs.ToSortedSlice() {
			s.WriteString(" <")
			s.WriteString(Escape(x))
			s.WriteString(">")
		}
	}
	if me.ref != "" {
		fmt.Fprintf(&s, " ref %s", me.ref)
	}
	return s.String()
}

func (me MetaField) Kind() FieldKind {
	return me.kind
}

func (me MetaField) IsNullable() bool {
	return me.flag.isNullable()
}

func (me MetaField) IsUnique() bool {
	return me.flag.isUnique()
}

func validMinOrMax(kind FieldKind, x any, what string) error {
	switch x := x.(type) {
	case int:
		if kind == BoolKind || kind == DateKind || kind == DateTimeKind {
			return fmt.Errorf("error#%d:cannot set a %s for a %s field",
				eWrongType, what, kind)
		}
		if x < 0 && (kind == BytesKind || kind == StrKind) {
			return fmt.Errorf(
				"error#%d:cannot set a negative %s for a %s field",
				eInvalidLength, what, kind)
		}
	case float64:
		if kind != RealKind {
			return fmt.Errorf(
				"error#%d:cannot set a real %s for a %s field", eWrongType,
				what, kind)
		}
	case time.Time:
		if !(kind == DateKind || kind == DateTimeKind) {
			return fmt.Errorf(
				"error#%d:cannot set a time.Time %s for a %s field",
				eWrongType, what, kind)
		}
	default:
		return fmt.Errorf("error#%d:cannot set a %s for a %s field",
			eWrongType, what, kind)
	}
	return nil
}
