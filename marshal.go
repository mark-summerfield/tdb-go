// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

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
					marshalTable(&out, name, field.Index(0).Interface())
					for i := 0; i < field.Len(); i++ {
						record := field.Index(i).Interface()
						if err := marshalRecord(&out, record); err != nil {
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

func marshalTable(out *bytes.Buffer, name string, table any) error {
	tableVal := reflect.ValueOf(table)
	tableType := reflect.TypeOf(table)
	//fmt.Printf("tableVal %T %v %s %s\n", tableVal, tableVal, tableVal.Kind(), name)
	out.WriteByte('[')
	out.WriteString(name)
	for i := 0; i < tableVal.NumField(); i++ {
		field := tableVal.Field(i)
		tag := tableType.Field(i).Tag.Get("tdb")
		if tag != "" {
			fmt.Println("TAG", tag)
		}
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
		// TODO []byte
		// TODO time.Time as Date or DateTime
		default:
			/* TODO reinstate
			return fmt.Errorf("error#%d:unrecognized field type %T",
				eUnrecognizedFieldType, field)
			*/
		}
	}
	out.WriteString("\n%\n")
	return nil
}

func marshalRecord(out *bytes.Buffer, record any) error {
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
		// TODO []byte
		// TODO time.Time as Date or DateTime
		default:
			return fmt.Errorf("error#%d:unrecognized field type %T",
				eUnrecognizedFieldType, field)
		}
		sep = " "
	}
	out.WriteByte('\n')
	return nil
}
