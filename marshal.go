// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"bytes"
	"fmt"
	"github.com/mark-summerfield/gset"
	"reflect"
	"strconv"
	"time"
)

var byteSlice []byte
var dateTime time.Time
var byteSliceType = reflect.TypeOf(byteSlice)
var dateTimeType = reflect.TypeOf(dateTime)

// Marshal converts the given struct of slices of structs to a string (as
// raw UTF-8-encoded bytes) in Tdb format if possible. For time.Time fields
// use a tag of either `tdb:"date"` or `tdb:"datetime"` to specify the Tdb
// field type; for all other types, the Tdb type is inferred.
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
							// return nil, err // TODO reinstate
						}
					}
					out.WriteString("]\n")
				}
			} else {
				// return nil, fmt.Errorf("cannot marshal %T", field) // TODO reinstate
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
	//fmt.Printf("tableVal %T %v %s %s\n", tableVal, tableVal, tableVal.Kind(), name)
	out.WriteByte('[')
	out.WriteString(name)
	for i := 0; i < tableVal.NumField(); i++ {
		field := tableVal.Field(i)
		tag := tableType.Field(i).Tag.Get("tdb")
		out.WriteByte(' ')
		out.WriteString(tableVal.Type().Field(i).Name)
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
	//fmt.Printf("recVal %T %v %s\n", recVal, recVal, recVal.Kind())
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
			out.WriteString(strconv.FormatInt(field.Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
			reflect.Uint64:
			out.WriteString(strconv.FormatUint(field.Uint(), 10))
		case reflect.Float32:
			out.WriteString(strconv.FormatFloat(field.Float(), 'f', -1, 32))
		case reflect.Float64:
			out.WriteString(strconv.FormatFloat(field.Float(), 'f', -1, 64))
		case reflect.String:
			out.WriteString(fmt.Sprintf("<%s>", Escape(field.String())))
		case reflect.Slice:
			x := field.Interface()
			if reflect.TypeOf(x) == byteSliceType {
				out.WriteByte('(')
				// TODO write bytes
				out.WriteByte(')')
			} else {
				return fmt.Errorf("unrecognized field type %T", field)
			}
		default:
			x := field.Interface()
			if reflect.TypeOf(x) == dateTimeType {
				if dateIndexes.Contains(i) {
					out.WriteString("date")
				} else {
					out.WriteString("datetime")
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
