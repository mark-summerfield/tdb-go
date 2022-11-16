// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

// Marshal converts the given database struct to a string (as raw
// UTF-8-encoded bytes) in Tdb format. The database struct should have a
// tdb.MetaData field and one or more slices of structs (each slice holding
// a table's records, each struct a record's fields).
func Marshal(db any) ([]byte, error) {
	// The format to use is:
	// [tablename fieldname1 type1 ... fieldnameN typeN
	// %
	// row0field0 ... row0fieldN
	//		:
	// rowMfield0 ... rowMfieldN
	// ]
	var out bytes.Buffer
	dbVal := reflect.ValueOf(db)
	if dbVal.Kind() == reflect.Ptr {
		dbVal = dbVal.Elem()
	}
	if dbVal.Kind() != reflect.Struct {
		return nil, fmt.Errorf("error#%d: expected a struct", eNotADatabase)
	}
	tableDefs := make(map[string]string)
	for i := 0; i < dbVal.NumField(); i++ {
		field := dbVal.Field(i)
		name := dbVal.Type().Field(i).Name
		switch field.Kind() {
		case reflect.Struct:
			if name != "MetaData" {
				return nil, fmt.Errorf("error#%d: expected metadata",
					eExpectedMetaData)
			}
			meta := field.Interface()
			populateTableDefs(meta, tableDefs)
		case reflect.Slice:
			if tableDef, ok := tableDefs[name]; ok {
				out.WriteString(tableDef)
				for j := 0; j < field.Len(); j++ {
					record := field.Index(j).Interface()
					if err := marshalRecord(&out, record); err != nil {
						// return nil, err // TODO reinstate
					}
				}
			} else {
				return nil, fmt.Errorf("error#%d: no metadata for %s",
					eSliceWithoutMetaData, name)
			}
			out.WriteString("\n]\n")
		default:
			return nil, fmt.Errorf(
				"error#%d: expected metadata or slices of records",
				eUnexpectedContent)
		}
	}
	return out.Bytes(), nil
}

func populateTableDefs(meta any, tableDefs map[string]string) {
	metaVal := reflect.ValueOf(meta)
	tables := metaVal.FieldByName("Tables")
	for i := 0; i < tables.Len(); i++ {
		if table, ok := tables.Index(i).Interface().(*MetaTable); ok {
			tableDefs[table.Name] = fmt.Sprintf("%s\n%%\n", table)
		}
	}
}

func marshalRecord(out *bytes.Buffer, record any) error {
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
