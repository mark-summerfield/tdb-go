// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"fmt"
	"time"
	"unicode/utf8"
)

func (me MetaField) check(value any) []string {
	result := make([]string, 0)
	switch value := value.(type) {
	case nil:
		if !me.IsNullable() {
			result = append(result, fmt.Sprintf(
				"warning#%d:%s is a not null field but has a null value",
				wInvalidNull, me.name))
		}
	case bool:
		if me.kind != BoolKind {
			result = append(result, fmt.Sprintf(
				"warning#%d:%s should be type %s", wWrongType, me.name,
				BoolKind))
			return result
		}
	case []byte:
		warnings := me.checkBytes(value)
		if len(warnings) > 0 {
			result = append(result, warnings...)
		}
	case time.Time:
		warnings := me.checkDate(value)
		if len(warnings) > 0 {
			result = append(result, warnings...)
		}
	case int:
		warnings := me.checkInt(value)
		if len(warnings) > 0 {
			result = append(result, warnings...)
		}
	case float64:
		warnings := me.checkReal(value)
		if len(warnings) > 0 {
			result = append(result, warnings...)
		}
	case string:
		warnings := me.checkStr(value)
		if len(warnings) > 0 {
			result = append(result, warnings...)
		}
	}
	return result
}

func (me MetaField) checkBytes(value []byte) []string {
	result := make([]string, 0)
	if me.kind != BytesKind {
		result = append(result, fmt.Sprintf(
			"warning#%d:%s should be type %s", wWrongType, me.name,
			BytesKind))
	}
	if me.min != nil {
		if m, ok := me.min.(int); ok {
			if len(value) <= m {
				result = append(result, fmt.Sprintf(
					"warning#%d:%s len is %d; must be at least %d",
					wLengthOutOfRange, me.name, value, m))
			}
		}
	}
	if me.max != nil {
		if m, ok := me.max.(int); ok {
			if len(value) >= m {
				result = append(result, fmt.Sprintf(
					"warning#%d:%s len is %d; must be at most %d",
					wLengthOutOfRange, me.name, value, m))
			}
		}
	}
	return result
}

func (me MetaField) checkDate(value time.Time) []string {
	result := make([]string, 0)
	if !(me.kind == DateKind || me.kind == DateTimeKind) {
		result = append(result, fmt.Sprintf(
			"warning#%d:%s should be a %s or %s", wWrongType, me.name,
			DateKind, DateTimeKind))
	}
	if me.min != nil {
		if m, ok := me.min.(time.Time); ok {
			if value.Add(-time.Second).Before(m) {
				result = append(result, fmt.Sprintf(
					"warning#%d: %s %s is not at or before %s",
					wValueOutOfRange, me.name, value, m))
			}
		}
	}
	if me.max != nil {
		if m, ok := me.max.(time.Time); ok {
			if value.Add(time.Second).After(m) {
				result = append(result, fmt.Sprintf(
					"warning#%d:%s %s is not at or after %s",
					wValueOutOfRange, me.name, value, m))
			}
		}
	}
	return result
}

func (me MetaField) checkInt(value int) []string {
	result := make([]string, 0)
	if me.kind != IntKind {
		result = append(result, fmt.Sprintf(
			"warning#%d:%s should be type %s", wWrongType, me.name,
			IntKind))
	}
	if me.inInts != nil {
		if me.inInts.Contains(value) {
			result = append(result, fmt.Sprintf(
				"warning#%d:%s should be in %v", wValueNotAllowed, me.name,
				me.inInts.ToSortedSlice()))
		}
	}
	if me.min != nil {
		if m, ok := me.min.(int); ok {
			if value < m {
				result = append(result, fmt.Sprintf(
					"warning#%d:%s %d is not >= %d", wValueOutOfRange,
					me.name, value, m))
			}
		}
	}
	if me.max != nil {
		if m, ok := me.max.(int); ok {
			if value > m {
				result = append(result, fmt.Sprintf(
					"warning#%d:%s %d is not <= %d", wValueOutOfRange,
					me.name, value, m))
			}
		}
	}
	return result
}

func (me MetaField) checkReal(value float64) []string {
	result := make([]string, 0)
	if me.kind != RealKind {
		result = append(result, fmt.Sprintf(
			"warning#%d:%s should be type %s", wWrongType, me.name,
			RealKind))
	}
	if me.min != nil {
		if m, ok := me.min.(float64); ok {
			if value < m {
				result = append(result, fmt.Sprintf(
					"warning#%d:%s %g is not >= %g", wValueOutOfRange,
					me.name, value, m))
			}
		}
	}
	if me.max != nil {
		if m, ok := me.max.(float64); ok {
			if value > m {
				result = append(result, fmt.Sprintf(
					"warning#%d:%s %g is not <= %g", wValueOutOfRange,
					me.name, value, m))
			}
		}
	}
	return result
}

func (me MetaField) checkStr(value string) []string {
	result := make([]string, 0)
	if me.kind != StrKind {
		result = append(result, fmt.Sprintf(
			"warning#%d:%s should be type %s", wWrongType, me.name,
			StrKind))
	}
	if me.inStrs != nil {
		if !me.inStrs.Contains(value) {
			result = append(result, fmt.Sprintf(
				"warning#%d:%s should be in %v", wValueNotAllowed, me.name,
				me.inStrs))
		}
	}
	size := utf8.RuneCountInString(value)
	if me.min != nil {
		if m, ok := me.min.(int); ok {
			if size < m {
				result = append(result, fmt.Sprintf(
					"warning#%d:%s %q len is %d; must be at least %d",
					wLengthOutOfRange, me.name, value, size, m))
			}
		}
	}
	if me.max != nil {
		if m, ok := me.max.(int); ok {
			if size > m {
				result = append(result, fmt.Sprintf(
					"warning#%d: %s %q len is %d; must be at most %d",
					wLengthOutOfRange, me.name, value, size, m))
			}
		}
	}
	return result
}
