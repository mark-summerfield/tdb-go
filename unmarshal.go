// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Unmarshal reads the data from the given string (as raw UTF-8-encoded
// bytes) into a (pointer to a) database struct.
//
// See also [Parse] and [Marshal] and [MarshalDecimals].
func Unmarshal(data []byte, db any) error {
	dbVal, err := getDbValue(data, db)
	if err != nil {
		return err
	}
	tableNames := getTableNames(dbVal)
	metaData := make(metaDataType)
	var metaTable *MetaTableType
	lino := 1
	for len(data) > 0 {
		b := data[0]
		data = data[1:]
		if b == '[' {
			data, metaTable, err = unmarshalTableMetaData(data, metaData,
				dbVal, &lino)
			if err != nil {
				return err
			}
		} else if metaTable != nil {
			if data, err = unmarshalRecords(data, metaTable, dbVal,
				tableNames, &lino); err == nil {
				metaTable = nil
			} else {
				return err
			}
		}
	}
	return nil
}

func getDbValue(data []byte, db any) (reflect.Value, error) {
	var zero reflect.Value
	if len(data) < 10 {
		return zero, fmt.Errorf("e%d#data holds invalid Tdb text", e107)
	}
	dbPtr := reflect.ValueOf(db)
	if dbPtr.Kind() != reflect.Ptr {
		return zero, fmt.Errorf("e%d#target interface must be a pointer",
			e108)
	}
	dbVal := dbPtr.Elem()
	if dbVal.Kind() != reflect.Struct {
		return zero, fmt.Errorf(
			"e%d#target interface must be a pointer to a struct", e109)
	}
	return dbVal, nil
}

func getTableNames(dbVal reflect.Value) map[string]string {
	// key=tableName | tagName value=tableName
	tableNames := make(map[string]string)
	dbType := dbVal.Type()
	for i := 0; i < dbVal.NumField(); i++ {
		tableName := dbVal.Type().Field(i).Name
		tableNames[tableName] = tableName
		tableTagName := dbType.Field(i).Tag.Get("tdb")
		if tableTagName != "" {
			tableNames[tableTagName] = tableName
		}
	}
	return tableNames
}

func unmarshalTableMetaData(data []byte, metaData metaDataType,
	dbVal reflect.Value, lino *int) ([]byte, *MetaTableType, error) {
	end, err := scanToByte(data, '%', lino)
	if err != nil {
		return data, nil, err
	}
	parts := bytes.Fields(bytes.TrimSpace(data[:end]))
	var metaTable *MetaTableType
	var tableName string
	var fieldName string
	for i, part := range parts {
		if i == 0 {
			tableName = string(part)
			metaTable = metaData.addTable(tableName)
		} else if i%2 != 0 {
			fieldName = string(part)
		} else {
			if err := addField(fieldName, string(part), metaTable,
				lino); err != nil {
				return data, nil, err
			}

		}
	}
	*lino++                             // allow for %
	return data[end+1:], metaTable, nil // +1 skips final %
}

func addField(fieldName, typeName string, metaTable *MetaTableType,
	lino *int) error {
	if fieldName == "" {
		return fmt.Errorf("e%d#%d:missing fieldname or type", e111, *lino)
	}
	if ok := metaTable.AddField(fieldName, typeName); !ok {
		return fmt.Errorf("e%d#%d:invalid typename %s", e112, *lino,
			typeName)
	}
	return nil
}

func unmarshalRecords(data []byte, metaTable *MetaTableType,
	dbVal reflect.Value, tableNames map[string]string, lino *int) ([]byte,
	error) {
	var err error
	var table reflect.Value
	var rec reflect.Value
	var recVal reflect.Value
	var field reflect.Value
	var metaField *MetaFieldType
	inRecord := false
	columns := metaTable.Len()
	oldColumn := -1
	column := 0
	for len(data) > 0 {
		if !inRecord {
			data, err = startRecord(data, &inRecord, &oldColumn, &column,
				lino)
			if err != nil {
				return data, err
			}
			table, rec, err = makeRecordType(metaTable.Name, dbVal,
				tableNames)
			if err != nil {
				return data, err
			}
			recVal = reflect.New(rec.Type().Elem()).Elem()
		}
		if column != oldColumn {
			oldColumn = column
			err = checkField(recVal, column, metaTable.Len(), *lino)
			if err != nil {
				return data, err
			}
			field = recVal.Field(column)
			metaField = metaTable.Field(column)
		}
		switch data[0] {
		case '\n': // ignore whitespace separators
			data = data[1:]
			*lino++
		case ' ', '\t', '\r': // ignore whitespace separators
			data = data[1:]
		case '?':
			data, err = unmarshalNull(data, metaField, field, lino)
			column++
		case 'F', 'f', 'N', 'n':
			data, err = unmarshalBool(data, false, metaField, field, lino)
			column++
		case 'T', 't', 'Y', 'y':
			data, err = unmarshalBool(data, true, metaField, field, lino)
			column++
		case '(':
			data, err = unmarshalBytes(data, metaField, field, lino)
			column++
		case '<':
			data, err = unmarshalStr(data, metaField, field, lino)
			column++
		case '-':
			switch metaField.Kind {
			case IntField:
				data, err = unmarshalInt(data, metaField, field, lino)
			case RealField:
				data, err = unmarshalReal(data, metaField, field, lino)
			default:
				err = fmt.Errorf("e%d#%d:got -, expected %s", e118, *lino,
					metaField.Kind)
			}
			column++
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			switch metaField.Kind {
			case BoolField:
				if (data[0] == '0' || data[0] == '1') && len(data) > 1 &&
					bytes.IndexByte([]byte{'.', 'e', 'E', '0', '1', '2',
						'3', '4', '5', '6', '7', '8', '9'}, data[1]) == -1 {
					data, err = unmarshalBool(data, data[0] == '1',
						metaField, field, lino)
				} else {
					err = fmt.Errorf("e%d#%d:got %c%c, expected %s", e130,
						*lino, data[0], data[1], metaField.Kind)
				}
			case IntField:
				data, err = unmarshalInt(data, metaField, field, lino)
			case RealField:
				data, err = unmarshalReal(data, metaField, field, lino)
			case DateField:
				data, err = unmarshalDateTime(data, DateFormat, metaField,
					field, lino)
			case DateTimeField:
				data, err = unmarshalDateTime(data, DateTimeFormat,
					metaField, field, lino)
			default: // Should never happend
				err = fmt.Errorf("e%d#%d:got %c, expected %s", e119, *lino,
					data[0], metaField.Kind)
			}
			column++
		case ']': // end of table
			if column > 0 && column < columns {
				err = fmt.Errorf(
					"e%d#%d:incomplete record %d/%d fields", e120, *lino,
					column+1, columns)
			} else {
				return skipWs(data[1:], lino), nil
			}
		default:
			err = fmt.Errorf("e%d#%d:invalid character %q", e121, *lino,
				rune(data[0]))
		}
		if err != nil {
			return data, err
		}
		if column == columns {
			table.Set(reflect.Append(table, recVal))
			oldColumn = -1
			column = 0
			inRecord = false
		}
	}
	return data, nil
}

func makeRecordType(tableName string, dbVal reflect.Value,
	tableNames map[string]string) (reflect.Value, reflect.Value, error) {
	field := dbVal.FieldByNameFunc(func(name string) bool {
		return name == tableName || name == tableNames[tableName]
	})
	if field.Kind() == reflect.Invalid {
		return field, field, fmt.Errorf("e%d#:invalid record type for %q",
			e128, tableName)
	}
	return field, reflect.New(field.Type().Elem()), nil
}

func startRecord(data []byte, inRecord *bool, oldColumn, column,
	lino *int) ([]byte, error) {
	*inRecord = true
	*oldColumn = -1
	*column = 0
	data = skipWs(data, lino)
	if len(data) == 0 {
		return data, fmt.Errorf("e%d#%d:unexpected end of data", e113,
			*lino)
	}
	return data, nil
}

func checkField(recVal reflect.Value, column, size, lino int) error {
	if !recVal.Type().Field(column).IsExported() {
		return fmt.Errorf(
			"e%d#%d:can't unmarshal to an unexported field: %q",
			e122, lino, recVal.Type().Field(column).Name)
	}
	if column >= size {
		return fmt.Errorf("e%d#%d:missing field name or type", e129, lino)
	}
	return nil
}

func unmarshalNull(data []byte, metaField *MetaFieldType,
	field reflect.Value, lino *int) ([]byte, error) {
	data = data[1:]
	if metaField.AllowNull {
		field.Set(reflect.Zero(field.Type()))
	} else {
		return data, fmt.Errorf("e%d#%d:can't write null to a not null "+
			"field: provide a valid %s or change the field's type to %s?",
			e115, lino, metaField.Kind, metaField.Kind)
	}
	return data, nil
}

func unmarshalBool(data []byte, value bool, metaField *MetaFieldType,
	field reflect.Value, lino *int) ([]byte, error) {
	if metaField.Kind != BoolField {
		return data, fmt.Errorf("e%d#%d:got bool, expected %s", e114, *lino,
			metaField.Kind)
	}
	if field.Kind() == reflect.Ptr {
		pv := reflect.ValueOf(&value)
		field.Set(pv)
	} else {
		field.SetBool(value)
	}
	return data[1:], nil
}

func unmarshalBytes(data []byte, metaField *MetaFieldType,
	field reflect.Value, lino *int) ([]byte, error) {
	data = data[1:] // skip (
	if metaField.Kind != BytesField {
		return data, fmt.Errorf("e%d#%d:got bytes, expected %s", e116,
			*lino, metaField.Kind)
	}
	data, raw, err := readHexBytes(data, lino)
	if err != nil {
		return data, err
	}
	if field.Kind() == reflect.Ptr {
		pv := reflect.ValueOf(&raw)
		field.Set(pv)
	} else {
		field.SetBytes(raw)
	}
	return data, nil
}

func unmarshalStr(data []byte, metaField *MetaFieldType,
	field reflect.Value, lino *int) ([]byte, error) {
	data = data[1:] // skip <
	if metaField.Kind != StrField {
		return data, fmt.Errorf("e%d#%d:got str, expected %s", e117, *lino,
			metaField.Kind)
	}
	data, s, err := readString(data, lino)
	if err != nil {
		return data, err
	}
	if field.Kind() == reflect.Ptr {
		pv := reflect.ValueOf(&s)
		field.Set(pv)
	} else {
		field.SetString(s)
	}
	return data, nil
}

func unmarshalInt(data []byte, metaField *MetaFieldType,
	field reflect.Value, lino *int) ([]byte, error) {
	data, i, err := readInt(data, lino)
	if err != nil {
		return data, err
	}
	if field.Kind() == reflect.Ptr {
		pv := reflect.ValueOf(&i)
		field.Set(pv)
	} else {
		field.SetInt(int64(i))
	}
	return data, nil
}

func unmarshalReal(data []byte, metaField *MetaFieldType,
	field reflect.Value, lino *int) ([]byte, error) {
	data, r, err := readReal(data, lino)
	if err != nil {
		return data, err
	}
	if field.Kind() == reflect.Ptr {
		pv := reflect.ValueOf(&r)
		field.Set(pv)
	} else {
		field.SetFloat(r)
	}
	return data, nil
}

func unmarshalDateTime(data []byte, format string, metaField *MetaFieldType,
	field reflect.Value, lino *int) ([]byte, error) {
	data, d, err := readDateTime(data, format, lino)
	if err != nil {
		return data, err
	}
	if field.Kind() == reflect.Ptr {
		pv := reflect.ValueOf(&d)
		field.Set(pv)
	} else {
		field.Set(reflect.ValueOf(d))
	}
	return data, err
}

func readHexBytes(data []byte, lino *int) ([]byte, []byte, error) {
	end, err := scanToByte(data, ')', lino)
	if err != nil {
		return data, nil, err
	}
	chunks := bytes.Fields(data[:end])
	chunk := bytes.Join(chunks, emptyBytes)
	raw := make([]byte, hex.DecodedLen(len(chunk)))
	_, err = hex.Decode(raw, chunk)
	if err != nil {
		return data, nil, fmt.Errorf("e%d#%d:invalid bytes %q", e123,
			*lino, chunk)
	}
	return data[end+1:], raw, nil // +1 skips final )
}

func readString(data []byte, lino *int) ([]byte, string, error) {
	end, err := scanToByte(data, '>', lino)
	if err != nil {
		return data, "", err
	}
	s := Unescape(string(data[:end]))
	return data[end+1:], s, nil // +1 skips final >
}

func readInt(data []byte, lino *int) ([]byte, int, error) {
	data, raw, err := scan(data, []byte("-+0123456789"), lino)
	if err != nil {
		return data, 0, err
	}
	x, err := strconv.Atoi(string(raw))
	if err != nil {
		return data, 0, fmt.Errorf("e%d#%d:invalid int", e125, *lino)
	}
	return data, int(x), nil
}

func readReal(data []byte, lino *int) ([]byte, float64, error) {
	data, raw, err := scan(data, []byte("-+0123456789.eE"), lino)
	if err != nil {
		return data, 0, err
	}
	x, err := strconv.ParseFloat(string(raw), 64)
	if err != nil {
		return data, 0, fmt.Errorf("e%d#%d:invalid real", e126, *lino)
	}
	return data, x, nil
}

func readDateTime(data []byte, format string,
	lino *int) ([]byte, time.Time, error) {
	data, raw, err := scan(data, []byte("-0123456789T:"), lino)
	if err != nil {
		return data, time.Now(), err
	}
	x, err := time.Parse(format, string(raw))
	if err != nil {
		what := "date"
		if strings.LastIndexByte(format, 'T') != -1 {
			what = "datetime"
		}
		return data, time.Now(), fmt.Errorf("e%d#%d:invalid %s", e127,
			*lino, what)
	}
	return data, x, nil
}

func scan(data, valid []byte, lino *int) ([]byte, []byte, error) {
	data = skipWs(data, lino)
	end := 0
	for end < len(data) {
		b := data[end]
		if bytes.IndexByte(valid, b) == -1 { // end of search
			return data[end:], data[:end], nil
		}
		end++
	}
	return data, emptyBytes, fmt.Errorf("e%d#%d:unexpected end of data",
		e124, *lino)
}

func skipWs(data []byte, lino *int) []byte {
	end := 0
	for end < len(data) {
		b := data[end]
		if b == '\n' {
			*lino++
		}
		if bytes.IndexByte([]byte{' ', '\t', '\n', '\r'}, b) == -1 {
			return data[end:]
		}
		end++
	}
	return data
}

func scanToByte(data []byte, b byte, lino *int) (int, error) {
	end := bytes.IndexByte(data, b)
	if end == -1 {
		return 0, fmt.Errorf("e%d#%d:missing %q", e110, *lino, b)
	}
	*lino += bytes.Count(data[:end], []byte{'\n'})
	return end, nil
}
