// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/mark-summerfield/gset"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Marshal converts the given struct of slices of structs to a string (as
// raw UTF-8-encoded bytes) in Tdb format if possible.
//
// Each tablename is taken from the outer struct's fieldname, but this can
// be overridden using a tag, e.g., `tdb:"MyTableName"`.
// For time.Time fields use a tag of either `tdb:"date"` or `tdb:"datetime"`
// to specify the Tdb field type; for all other types, the Tdb type is
// inferred. However, if fieldnames in the Tdb text are to be different from
// the struct fieldnames, use tags, with the required name, e.g.,
// `tdb:"MyFieldName"`, and for dates and datetimes with the type too, e.g.,
// `tdb:"MyDateField:date"`, etc.
//
// See also [Tdb.Write] and [MarshalDecimals] and [Unmarshal].
func Marshal(db any) ([]byte, error) {
	return MarshalDecimals(db, -1)
}

// MarshalDecimals is a refinement of the Marshal function.
//
// By default for real numbers the [Marshal] function outputs them using the
// fewest number of decimal digits necessary. In particular this means that
// numbers whose fractional part is 0 are output like ints (e.g., 2.0 → 2).
//
// To control decimal output use this function. Pass a decimals value of
// 1-19 to use exactly that number of decimal digits; any other value means
// use the minimum number of decimal digits necessary (which may be none for
// numbers whose fractional part is 0).
//
// See also [Marshal] and [Unmarshal].
func MarshalDecimals(db any, decimals int) ([]byte, error) {
	var out bytes.Buffer
	dbVal := reflect.ValueOf(db)
	if dbVal.Kind() == reflect.Ptr {
		dbVal = dbVal.Elem()
	}
	if dbVal.Kind() == reflect.Struct {
		dp := sanitizedDecimals(decimals)
		dbType := dbVal.Type()
		for i := 0; i < dbVal.NumField(); i++ {
			field := dbVal.Field(i)
			tableName := dbType.Field(i).Tag.Get("tdb")
			if tableName == "" {
				tableName = dbVal.Type().Field(i).Name
			}
			if field.Kind() == reflect.Slice {
				if field.Len() > 0 {
					if err := marshalTable(&out, field, tableName,
						dp); err != nil {
						return nil, err
					}
				}
			} else {
				return nil, fmt.Errorf(
					"e%d#%s: cannot marshal outer struct field %T",
					e100, tableName, field)
			}
		}
	} else {
		return nil, fmt.Errorf("e%d#cannot marshal %T", e101, dbVal)
	}
	if out.Len() == 0 {
		return nil, fmt.Errorf("e%d#cannot marshal empty data", e102)
	}
	return out.Bytes(), nil
}

func marshalTable(out *bytes.Buffer, field reflect.Value, tableName string,
	dp int) error {
	if field.Len() > 0 {
		dateIndexes, fieldNameForIndex, err := marshalMetaData(
			out, tableName, field.Index(0).Interface())
		if err != nil {
			return err
		}
		for i := 0; i < field.Len(); i++ {
			record := field.Index(i).Interface()
			if err := marshalRecord(out, record, dateIndexes, tableName,
				fieldNameForIndex, dp); err != nil {
				return err
			}
		}
		out.WriteString("]\n")
	}
	return nil
}

func marshalMetaData(out *bytes.Buffer, tableName string,
	table any) (gset.Set[int], map[int]string, error) {
	dateIndexes := gset.New[int]()
	fieldNameForIndex := make(map[int]string)
	tableVal := reflect.ValueOf(table)
	tableType := reflect.TypeOf(table)
	out.WriteByte('[')
	out.WriteString(tableName)
	for i := 0; i < tableVal.NumField(); i++ {
		field := tableVal.Field(i)
		fieldName := tableVal.Type().Field(i).Name
		tag := tableType.Field(i).Tag.Get("tdb")
		var typeName string
		if tag != "" {
			fieldName, typeName = readTag(fieldName, tag)
		}
		fieldNameForIndex[i] = fieldName
		isDate, err := marshalTableMetaData(out, field, typeName, tableName,
			fieldName)
		if err != nil {
			return dateIndexes, fieldNameForIndex, err
		}
		if isDate {
			dateIndexes.Add(i)
		}
	}
	out.WriteString("\n%\n")
	return dateIndexes, fieldNameForIndex, nil
}

func marshalTableMetaData(out *bytes.Buffer, field reflect.Value, typeName,
	tableName, fieldName string) (bool, error) {
	isDate := false
	out.WriteByte(' ')
	out.WriteString(fieldName)
	out.WriteByte(' ')
	switch field.Kind() {
	case reflect.Bool:
		out.WriteString("bool")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		out.WriteString("int")
	case reflect.Float32, reflect.Float64:
		out.WriteString("real")
	case reflect.String:
		out.WriteString("str")
	case reflect.Slice:
		x := field.Interface()
		if reflect.TypeOf(x) == byteSliceType {
			out.WriteString("bytes")
		} else {
			return isDate, fmt.Errorf(
				"e%d#%s.%s:unrecognized field slice type %T", e103,
				tableName, fieldName, field)
		}
	default:
		x := field.Interface()
		if reflect.TypeOf(x) == dateTimeType {
			if typeName == "date" {
				isDate = true
				out.WriteString(typeName)
			} else {
				out.WriteString("datetime")
			}
		} else {
			return isDate, fmt.Errorf(
				"e%d#%s.%s:unrecognized field type %T", e104, tableName,
				fieldName, x)
		}
	}
	return isDate, nil
}

func marshalRecord(out *bytes.Buffer, record any, dateIndexes gset.Set[int],
	tableName string, fieldNameForIndex map[int]string, dp int) error {
	recVal := reflect.ValueOf(record)
	sep := ""
	for i := 0; i < recVal.NumField(); i++ {
		out.WriteString(sep)
		field := recVal.Field(i)
		// TODO handle nulls
		switch field.Kind() {
		case reflect.Bool:
			if field.Bool() {
				out.WriteByte('T')
			} else {
				out.WriteByte('F')
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Int64:
			out.WriteString(strconv.FormatInt(field.Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
			reflect.Uint64:
			out.WriteString(strconv.FormatUint(field.Uint(), 10))
		case reflect.Float32:
			out.WriteString(strconv.FormatFloat(field.Float(), 'f', dp, 32))
		case reflect.Float64:
			out.WriteString(strconv.FormatFloat(field.Float(), 'f', dp, 64))
		case reflect.String:
			out.WriteString(fmt.Sprintf("<%s>", Escape(field.String())))
		case reflect.Slice:
			if err := marshalSliceField(out, field, tableName,
				fieldNameForIndex[i]); err != nil {
				return err
			}
		default:
			if err := marshalDateTimeField(out, field, tableName,
				fieldNameForIndex[i], dateIndexes.Contains(i)); err != nil {
				return err
			}
		}
		sep = " "
	}
	out.WriteByte('\n')
	return nil
}

func marshalSliceField(out *bytes.Buffer, field reflect.Value, tableName,
	fieldName string) error {
	x := field.Interface()
	if reflect.TypeOf(x) == byteSliceType {
		raw := field.Bytes()
		out.WriteByte('(')
		out.WriteString(hex.EncodeToString(raw))
		out.WriteByte(')')
	} else {
		return fmt.Errorf("e%d#%s.%s:unrecognized slice's field type %T",
			e105, tableName, fieldName, x)
	}
	return nil
}

func marshalDateTimeField(out *bytes.Buffer, field reflect.Value, tableName,
	fieldName string, isDate bool) error {
	x := field.Interface()
	if d, ok := x.(time.Time); ok {
		var s string
		if isDate {
			s = d.Format(DateFormat)
		} else {
			s = d.Format(DateTimeFormat)
		}
		out.WriteString(s)
	} else {
		return fmt.Errorf(
			"e%d#%s.%s:unrecognized field type (expected time.Time) %T",
			e106, tableName, fieldName, field)
	}
	return nil
}

func readTag(name, tag string) (string, string) {
	i := strings.IndexByte(tag, ':')
	if i == -1 {
		if reservedWords.Contains(tag) { // `tdb:"type"`
			return name, tag
		}
		return tag, "" // `tdb:"FieldName"`
	} else {
		left := tag[:i]
		right := tag[i+1:]
		if reservedWords.Contains(left) { // `tdb:type:FieldName"`
			return right, left // invalid format, but understandable
		}
		if reservedWords.Contains(right) { // `tdb:FieldName:type"`
			return left, right
		}
		return name, "" // `tdb:FieldNameA:FieldNameB"`
	}
}

func sanitizedDecimals(decimals int) int {
	if 1 <= decimals && decimals <= 19 {
		return decimals
	}
	return -1
}
