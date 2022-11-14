// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	_ "embed"
	"fmt"
	"github.com/mark-summerfield/gset"
	"time"
	"unicode/utf8"
)

func BoolMeta(theDefault bool) Metadatum {
	return Metadatum{kind: BoolKind, flag: notNullFlag,
		theDefault: theDefault}
}

func BytesMeta(theDefault []byte) Metadatum {
	return Metadatum{kind: BytesKind, flag: notNullFlag,
		theDefault: theDefault}
}

func DateMeta(theDefault time.Time) Metadatum {
	return Metadatum{kind: DateKind, flag: notNullFlag,
		theDefault: theDefault}
}

func DateTimeMeta(theDefault time.Time) Metadatum {
	return Metadatum{kind: DateTimeKind, flag: notNullFlag,
		theDefault: theDefault}
}

func IntMeta(theDefault int) Metadatum {
	return Metadatum{kind: IntKind, flag: notNullFlag,
		theDefault: theDefault}
}

func AutoMeta() Metadatum {
	return Metadatum{kind: IntKind, flag: autoFlag, next: 1}
}

func RealMeta(theDefault float64) Metadatum {
	return Metadatum{kind: RealKind, flag: notNullFlag,
		theDefault: theDefault}
}

func StrMeta(theDefault string) Metadatum {
	return Metadatum{kind: StrKind, flag: notNullFlag,
		theDefault: theDefault}
}

type Metadatum struct {
	kind       FieldKind
	flag       fieldFlag
	min        any // length for bytes & str; min value otherwise
	max        any // length for bytes & str; max value otherwise
	inInts     gset.Set[int]
	inStrs     gset.Set[string]
	theDefault any
	ref        string // tablename.fieldname or fieldname
	next       int    // for auto fields
}

func (me *Metadatum) SetNullable() {
	me.flag = me.flag.with(nullableFlag)
}

func (me *Metadatum) SetMin(min any) error {
	if err := checkNumericKind(me.kind, min, "min"); err != nil {
		return err
	}
	me.min = min
	return nil
}

func (me *Metadatum) SetMax(max any) error {
	if err := checkNumericKind(me.kind, max, "max"); err != nil {
		return err
	}
	me.max = max
	return nil
}

func (me *Metadatum) SetInInts(ints ...int) {
	if me.inInts == nil {
		me.inInts = gset.New[int]()
	} else {
		me.inInts.Clear()
	}
	me.inInts.Add(ints...)
}

func (me *Metadatum) SetInStrs(ints ...string) {
	if me.inStrs == nil {
		me.inStrs = gset.New[string]()
	} else {
		me.inStrs.Clear()
	}
	me.inStrs.Add(ints...)
}

func (me *Metadatum) SetRef(ref string) {
	// TODO check identifer.identifer or .identifer
	me.ref = ref
}

func (me Metadatum) Kind() FieldKind {
	return me.kind
}

func (me Metadatum) IsNullable() bool {
	return me.flag.isNullable()
}

func (me Metadatum) IsUnique() bool {
	return me.flag.isUnique()
}

func (me Metadatum) IsAuto() bool {
	return me.flag.isAuto()
}

func (me *Metadatum) Next() (int, bool) {
	if !me.IsAuto() {
		return 0, false
	}
	if me.min != nil {
		if m, ok := me.min.(int); ok {
			if me.next < m {
				me.next = m + 1
				return m, true
			}
		}
	}
	n := me.next
	me.next++
	return n, true
}

func (me Metadatum) IsInRange(fieldName string, value any) bool {
	switch value := value.(type) {
	case bool:
		if me.kind != BoolKind {
			return false
		}
	case []byte:
		if me.kind != BytesKind {
			return false
		}
		if me.min != nil {
			if m, ok := me.min.(int); ok {
				if len(value) < m {
					return false
				}
			}
		}
		if me.max != nil {
			if m, ok := me.max.(int); ok {
				if len(value) > m {
					return false
				}
			}
		}
	case time.Time:
		if !(me.kind == DateKind || me.kind == DateTimeKind) {
			return false
		}
		if me.min != nil {
			if m, ok := me.min.(time.Time); ok {
				if value.Before(m) {
					return false
				}
			}
		}
		if me.max != nil {
			if m, ok := me.max.(time.Time); ok {
				if value.After(m) {
					return false
				}
			}
		}
	case int:
		if me.kind != IntKind {
			return false
		}
		if me.inInts != nil {
			if me.inInts.Contains(value) {
				return true
			}
		}
		if me.min != nil {
			if m, ok := me.min.(int); ok {
				if value < m {
					return false
				}
			}
		}
		if me.max != nil {
			if m, ok := me.max.(int); ok {
				if value > m {
					return false
				}
			}
		}
	case float64:
		if me.kind != RealKind {
			return false
		}
		if me.min != nil {
			if m, ok := me.min.(float64); ok {
				if value < m {
					return false
				}
			}
		}
		if me.max != nil {
			if m, ok := me.max.(float64); ok {
				if value > m {
					return false
				}
			}
		}
	case string:
		if me.kind != StrKind {
			return false
		}
		if me.inStrs != nil {
			if me.inStrs.Contains(value) {
				return true
			}
		}
		size := utf8.RuneCountInString(value)
		if me.min != nil {
			if m, ok := me.min.(int); ok {
				if size < m {
					return false
				}
			}
		}
		if me.max != nil {
			if m, ok := me.max.(int); ok {
				if size > m {
					return false
				}
			}
		}
	}
	return false
}

func checkNumericKind(kind FieldKind, x any, what string) error {
	switch x := x.(type) {
	case int:
		if kind == BoolKind || kind == DateKind || kind == DateTimeKind {
			return fmt.Errorf("cannot set a %s for a %s field", what, kind)
		}
		if x < 0 && (kind == BytesKind || kind == StrKind) {
			return fmt.Errorf("cannot set a negative %s for a %s field",
				what, kind)
		}
	case float64:
		if kind != RealKind {
			return fmt.Errorf("cannot set a real %s for a %s field", what,
				kind)
		}
	case time.Time:
		if !(kind == DateKind || kind == DateTimeKind) {
			return fmt.Errorf("cannot set a time.Time %s for a %s field",
				what, kind)
		}
	default:
		return fmt.Errorf("cannot set a %s for a %s field", what, kind)
	}
	return nil
}
