// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/mark-summerfield/gong"
	"github.com/mark-summerfield/gset"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var byteSlice []byte
var dateTime time.Time
var byteSliceType = reflect.TypeOf(byteSlice)
var dateTimeType = reflect.TypeOf(dateTime)
var reservedWords gset.Set[string]

func init() {
	reservedWords = gset.New("bool", "bytes", "date", "datetime", "int",
		"real", "str")
}

// Marshal converts the given struct of slices of structs to a string (as
// raw UTF-8-encoded bytes) in Tdb format if possible. For time.Time fields
// use a tag of either `tdb:"date"` or `tdb:"datetime"` to specify the Tdb
// field type; for all other types, the Tdb type is inferred. However, if
// field names in the Tdb text are to be different from the struct
// fieldnames, use tags, with the required name, e.g., `tdb:"MyFieldName"`,
// and for dates and times with the type too, e.g.,
// `tdb:"MyDateField:date"`, etc.
func Marshal(db any) ([]byte, error) {
	var out bytes.Buffer
	dbVal := reflect.ValueOf(db)
	if dbVal.Kind() == reflect.Ptr {
		dbVal = dbVal.Elem()
	}
	if dbVal.Kind() == reflect.Struct {
		for i := 0; i < dbVal.NumField(); i++ {
			field := dbVal.Field(i)
			name := dbVal.Type().Field(i).Name
			if field.Kind() == reflect.Slice {
				if field.Len() > 0 {
					dateIndexes, err := marshalTable(&out, name,
						field.Index(0).Interface())
					if err != nil {
						return nil, err
					}
					for i := 0; i < field.Len(); i++ {
						record := field.Index(i).Interface()
						if err := marshalRecord(&out, record,
							dateIndexes); err != nil {
							return nil, err
						}
					}
					out.WriteString("]\n")
				}
			} else {
				return nil, fmt.Errorf("cannot marshal %T", field)
			}
		}
	} else {
		return nil, fmt.Errorf("cannot marshal %T", dbVal)
	}
	return out.Bytes(), nil
}

func marshalTable(out *bytes.Buffer, name string, table any) (gset.Set[int],
	error) {
	dateIndexes := gset.New[int]()
	tableVal := reflect.ValueOf(table)
	tableType := reflect.TypeOf(table)
	out.WriteByte('[')
	out.WriteString(name)
	for i := 0; i < tableVal.NumField(); i++ {
		field := tableVal.Field(i)
		tag := tableType.Field(i).Tag.Get("tdb")
		name := tableVal.Type().Field(i).Name
		if tag != "" {
			tag, name = parseTag(tag, name)
		}
		out.WriteByte(' ')
		out.WriteString(name)
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
				return dateIndexes, fmt.Errorf("unrecognized field type %T",
					field)
			}
		default:
			x := field.Interface()
			if reflect.TypeOf(x) == dateTimeType {
				if tag == "date" {
					dateIndexes.Add(i)
					out.WriteString(tag)
				} else {
					out.WriteString("datetime")
				}
			} else {
				return dateIndexes, fmt.Errorf("unrecognized field type %T",
					field)
			}
		}
	}
	out.WriteString("\n%\n")
	return dateIndexes, nil
}

func marshalRecord(out *bytes.Buffer, record any,
	dateIndexes gset.Set[int]) error {
	recVal := reflect.ValueOf(record)
	sep := ""
	for i := 0; i < recVal.NumField(); i++ {
		out.WriteString(sep)
		field := recVal.Field(i)
		switch field.Kind() {
		case reflect.Bool:
			if field.Bool() {
				out.WriteByte('T')
			} else {
				out.WriteByte('F')
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Int64:
			i := field.Int()
			if i == IntSentinal {
				out.WriteByte('!')
			} else {
				out.WriteString(strconv.FormatInt(i, 10))
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
			reflect.Uint64:
			out.WriteString(strconv.FormatUint(field.Uint(), 10))
		case reflect.Float32:
			r := field.Float()
			if gong.IsRealClose(r, RealSentinal) {
				out.WriteByte('!')
			} else {
				out.WriteString(strconv.FormatFloat(r, 'f', -1, 32))
			}
		case reflect.Float64:
			r := field.Float()
			if gong.IsRealClose(r, RealSentinal) {
				out.WriteByte('!')
			} else {
				out.WriteString(strconv.FormatFloat(r, 'f', -1, 64))
			}
		case reflect.String:
			s := field.String()
			if s == StrSentinal {
				out.WriteByte('!')
			} else {
				out.WriteString(fmt.Sprintf("<%s>", Escape(s)))
			}
		case reflect.Slice:
			x := field.Interface()
			if reflect.TypeOf(x) == byteSliceType {
				raw := field.Bytes()
				if len(raw) == 1 && raw[0] == BytesSentinal[0] {
					out.WriteByte('!')
				} else {
					out.WriteByte('(')
					out.WriteString(hex.EncodeToString(raw))
					out.WriteByte(')')
				}
			} else {
				return fmt.Errorf("unrecognized field type %T", field)
			}
		default:
			x := field.Interface()
			if d, ok := x.(time.Time); ok {
				if d.Equal(DateSentinal) {
					out.WriteByte('!')
				} else if dateIndexes.Contains(i) {
					out.WriteString(d.Format("2006-01-02"))
				} else {
					out.WriteString(d.Format("2006-01-02T15:04:05"))
				}
			} else {
				return fmt.Errorf("unrecognized field type %T", field)
			}
		}
		sep = " "
	}
	out.WriteByte('\n')
	return nil
}

func parseTag(tag, name string) (string, string) {
	i := strings.IndexByte(tag, ':')
	if i == -1 {
		if reservedWords.Contains(tag) { // `tdb:"type"`
			return tag, name
		}
		return tag, tag // `tdb:"FieldName"`
	} else {
		left := tag[:i]
		right := tag[i+1:]
		if reservedWords.Contains(left) { // `tdb:type:FieldName"`
			return left, right
		}
		if reservedWords.Contains(right) { // `tdb:FieldName:type"`
			return right, left
		}
		return "", tag // `tdb:FieldNameA:FieldNameB"`
	}
}
