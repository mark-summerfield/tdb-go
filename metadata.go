// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"fmt"
	"github.com/mark-summerfield/gset"
	"golang.org/x/exp/maps"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

type Metadata map[string]Metatable // Keys are tablenames

func (me Metadata) String() string {
	tableNames := maps.Keys(me)
	sort.Strings(tableNames)
	var s strings.Builder
	for _, tableName := range tableNames {
		s.WriteByte('[')
		s.WriteString(tableName)
		s.WriteByte('\n')
		s.WriteString(me[tableName].String())
	}
	return s.String()
}

func (me Metadata) Check() error {
	// TODO iterate over all tables and call their isValid method
	// TODO and for any which have refs, check that their is a corresponding
	// table & field of the same int or str type
	return nil
}

type Metatable struct {
	fields       []*Metadatum          // The fields in order
	fieldForName map[string]*Metadatum // Keys are fieldnames
}

func NewMetatable() Metatable {
	return Metatable{fields: make([]*Metadatum, 0, 1),
		fieldForName: make(map[string]*Metadatum)}
}

func (me Metatable) Field(index int) *Metadatum {
	return me.fields[index]
}

func (me Metatable) FieldByName(fieldName string) *Metadatum {
	if m, ok := me.fieldForName[fieldName]; ok {
		return m
	}
	return nil
}

func (me *Metatable) Add(datum Metadatum) {
	me.fields = append(me.fields, &datum)
	me.fieldForName[datum.name] = &datum
}

func (me Metatable) String() string {
	var s strings.Builder
	for _, datum := range me.fields {
		s.WriteByte(' ')
		s.WriteString(datum.String())
		s.WriteByte('\n')
	}
	return s.String()
}

func (me Metatable) check() error {
	// TODO iterate over all fields and call their check()
	// TODO and for any which are unique, check for uniqueness
	// TODO and for any which are not null, check for not nullness
	return nil
}

func BoolMeta(name string) Metadatum {
	return Metadatum{name: name, kind: BoolKind, flag: notNullFlag}
}

func BytesMeta(name string) Metadatum {
	return Metadatum{name: name, kind: BytesKind, flag: notNullFlag}
}

func DateMeta(name string) Metadatum {
	return Metadatum{name: name, kind: DateKind, flag: notNullFlag}
}

func DateTimeMeta(name string) Metadatum {
	return Metadatum{name: name, kind: DateTimeKind, flag: notNullFlag}
}

func IntMeta(name string) Metadatum {
	return Metadatum{name: name, kind: IntKind, flag: notNullFlag}
}

func RealMeta(name string) Metadatum {
	return Metadatum{name: name, kind: RealKind, flag: notNullFlag}
}

func StrMeta(name string) Metadatum {
	return Metadatum{name: name, kind: StrKind, flag: notNullFlag}
}

type Metadatum struct {
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

func (me *Metadatum) SetNullable() {
	me.flag = me.flag.with(nullableFlag)
}

// SetUnique sets the field to be unique: this also implies it is not
// nullable.
func (me *Metadatum) SetUnique() {
	me.flag = uniqueFlag
}

func (me *Metadatum) SetDefault(theDefault any) error {
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
	return fmt.Errorf("%v is not a valid default for a field of type %s",
		theDefault, me.kind)
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

func (me *Metadatum) SetRef(ref string) error {
	if !strings.Contains(ref, ".") && ref == me.name {
		return fmt.Errorf("cannot set a field ref to itself (%s)", ref)
	}
	// TODO check identifer.identifer or .identifer
	me.ref = ref
	return nil
}

func (me Metadatum) String() string {
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

func (me Metadatum) Kind() FieldKind {
	return me.kind
}

func (me Metadatum) IsNullable() bool {
	return me.flag.isNullable()
}

func (me Metadatum) IsUnique() bool {
	return me.flag.isUnique()
}

func (me Metadatum) check(fieldName string, value any) error {
	switch value := value.(type) {
	case bool:
		if me.kind != BoolKind {
			return fmt.Errorf("%s should be type %s", fieldName, BoolKind)
		}
	case []byte:
		if me.kind != BytesKind {
			return fmt.Errorf("%s should be type %s", fieldName, BytesKind)
		}
		if me.min != nil {
			if m, ok := me.min.(int); ok {
				if len(value) < m {
					return fmt.Errorf("%s %d is not > %d", fieldName, value,
						m)
				}
			}
		}
		if me.max != nil {
			if m, ok := me.max.(int); ok {
				if len(value) > m {
					return fmt.Errorf("%s %d is not < %d", fieldName, value,
						m)
				}
			}
		}
	case time.Time:
		if !(me.kind == DateKind || me.kind == DateTimeKind) {
			return fmt.Errorf("%s should be a %s or %s", fieldName,
				DateKind, DateTimeKind)
		}
		if me.min != nil {
			if m, ok := me.min.(time.Time); ok {
				if value.Before(m) {
					return fmt.Errorf("%s %s is not > %s", fieldName, value,
						m)
				}
			}
		}
		if me.max != nil {
			if m, ok := me.max.(time.Time); ok {
				if value.After(m) {
					return fmt.Errorf("%s %s is not < %s", fieldName, value,
						m)
				}
			}
		}
	case int:
		if me.kind != IntKind {
			return fmt.Errorf("%s should be type %s", fieldName, IntKind)
		}
		if me.inInts != nil {
			if me.inInts.Contains(value) {
				return fmt.Errorf("%s should be in %v", fieldName,
					me.inInts.ToSortedSlice())
			}
		}
		if me.min != nil {
			if m, ok := me.min.(int); ok {
				if value < m {
					return fmt.Errorf("%s %d is not > %d", fieldName, value,
						m)
				}
			}
		}
		if me.max != nil {
			if m, ok := me.max.(int); ok {
				if value > m {
					return fmt.Errorf("%s %d is not < %d", fieldName, value,
						m)
				}
			}
		}
	case float64:
		if me.kind != RealKind {
			return fmt.Errorf("%s should be type %s", fieldName, RealKind)
		}
		if me.min != nil {
			if m, ok := me.min.(float64); ok {
				if value < m {
					return fmt.Errorf("%s %g is not > %g", fieldName, value,
						m)
				}
			}
		}
		if me.max != nil {
			if m, ok := me.max.(float64); ok {
				if value > m {
					return fmt.Errorf("%s %g is not < %g", fieldName, value,
						m)
				}
			}
		}
	case string:
		if me.kind != StrKind {
			return fmt.Errorf("%s should be type %s", fieldName, StrKind)
		}
		if me.inStrs != nil {
			if me.inStrs.Contains(value) {
				return fmt.Errorf("%s should be in %v", fieldName,
					me.inStrs.ToSortedSlice())
			}
		}
		size := utf8.RuneCountInString(value)
		if me.min != nil {
			if m, ok := me.min.(int); ok {
				if size < m {
					return fmt.Errorf("%s %q is not > %d", fieldName, value,
						m)
				}
			}
		}
		if me.max != nil {
			if m, ok := me.max.(int); ok {
				if size > m {
					return fmt.Errorf("%s %q is not < %d", fieldName, value,
						m)
				}
			}
		}
	}
	return nil
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
