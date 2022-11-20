// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Unmarshal reads the data from the given string (as raw UTF-8-encoded
// bytes) into a (pointer to a) database struct.
func Unmarshal(data []byte, db any) error {
	var err error
	dbPtr, dbVal, err := getDbValues(data, db)
	fmt.Println(dbPtr, dbVal) // TODO delete
	metaData := make(metaDataType)
	var tableName string
	lino := 1
	for len(data) > 0 {
		b := data[0]
		data = data[1:]
		if b == '[' {
			data, tableName, err = readTableMetaData(data, metaData, &lino)
		} else if tableName != "" {
			fmt.Println("Start of Table", tableName) // TODO delete
			if data, err = readRecords(data, metaData[tableName],
				&lino); err == nil {
				fmt.Println("End of Table", tableName) // TODO delete
				tableName = ""
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func getDbValues(data []byte, db any) (reflect.Value, reflect.Value,
	error) {
	zero := reflect.Zero(reflect.TypeOf(false))
	if len(data) < 11 {
		return zero, zero, fmt.Errorf("e%d#data holds invalid Tdb text",
			e107)
	}
	dbPtr := reflect.ValueOf(db)
	if dbPtr.Kind() != reflect.Ptr {
		return zero, zero, fmt.Errorf(
			"e%d#target interface must be a pointer", e108)
	}
	dbVal := dbPtr.Elem()
	if dbVal.Kind() != reflect.Struct {
		return zero, zero, fmt.Errorf(
			"e%d#target interface must be a pointer to a struct", e109)
	}
	return dbPtr, dbVal, nil
}

func readTableMetaData(data []byte, metaData metaDataType,
	lino *int) ([]byte, string, error) {
	end, err := scanToByte(data, '%', lino)
	if err != nil {
		return data, "", err
	}
	parts := bytes.Fields(data[:end])
	var tableName string
	var fieldName string
	for i, part := range parts {
		if i == 0 {
			tableName = string(part)
			metaData[tableName] = newMetaTable(tableName)
		} else if i%2 != 0 {
			fieldName = string(part)
		} else {
			if err := addField(fieldName, string(part), metaData[tableName],
				lino); err != nil {
				return data, "", err
			}

		}
	}
	return data[end+1:], tableName, nil // +1 skips final %
}

func addField(fieldName, typeName string, metaTable *metaTableType,
	lino *int) error {
	if fieldName == "" {
		return fmt.Errorf("e%d#%d:missing fieldname or type", e111, *lino)
	}
	if ok := metaTable.Add(fieldName, typeName); !ok {
		return fmt.Errorf("e%d#%d:invalid typename %s", e112, *lino,
			typeName)
	}
	return nil
}

// TODO take in a reflect.Value for the outer target struct's corresponding
// slice
func readRecords(data []byte, metaTable *metaTableType, lino *int) ([]byte,
	error) {
	var err error
	var metaField *metaFieldType
	var values []reflect.Value // TODO -or- map[string]reflect.Value (key==fieldName)?
	var value *reflect.Value
	fieldCount := metaTable.Len()
	inRecord := false
	oldFieldIndex := -1
	fieldIndex := 0
	for len(data) > 0 {
		if !inRecord {
			data, err = maybeStartRecord(data, &inRecord, &oldFieldIndex,
				&fieldIndex, lino)
			if err != nil {
				return data, err
			}
			values = make([]reflect.Value, 0)
		}
		if fieldIndex != oldFieldIndex {
			oldFieldIndex = fieldIndex
			metaField = metaTable.Field(fieldIndex)
		}
		switch data[0] {
		case '\n': // ignore whitespace separators
			data = data[1:]
			*lino++
		case ' ', '\t', '\r': // ignore whitespace separators
			data = data[1:]
		case '!':
			data, value, err = handleSentinal(data, metaField, lino)
			fieldIndex += 1
		case 'F', 'f', 'N', 'n':
			data, value, err = handleBool(data, false, metaField, lino)
			fieldIndex += 1
		case 'T', 't', 'Y', 'y':
			data, value, err = handleBool(data, true, metaField, lino)
			fieldIndex += 1
		case '(':
			data, value, err = handleBytes(data, metaField, lino)
			fieldIndex += 1
		case '<':
			data, value, err = handleStr(data, metaField, lino)
			fieldIndex += 1
		case '-':
			switch metaField.kind {
			case intField:
				data, value, err = handleInt(data, metaField, lino)
			case realField:
				data, value, err = handleReal(data, metaField, lino)
			default:
				err = fmt.Errorf("e%d#%d:got -, expected %s", e118, *lino,
					metaField.kind)
			}
			fieldIndex += 1
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			switch metaField.kind {
			case intField:
				data, value, err = handleInt(data, metaField, lino)
			case realField:
				data, value, err = handleReal(data, metaField, lino)
			case dateField:
				data, value, err = handleDateTime(data, DateFormat,
					metaField, lino)
			case dateTimeField:
				data, value, err = handleDateTime(data, DateTimeFormat,
					metaField, lino)
			default:
				err = fmt.Errorf("e%d#%d:got -, expected %s", e119, *lino,
					metaField.kind)
			}
			fieldIndex += 1
		case ']': // end of table
			if fieldIndex > 0 && fieldIndex < fieldCount {
				err = fmt.Errorf(
					"e%d#%d:incomplete record %d/%d fields", e120, *lino,
					fieldIndex+1, fieldCount)
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
		if value != nil {
			values = append(values, *value)
		}
		if fieldIndex == fieldCount {
			inRecord = false
			// TODO add record based on values
			fmt.Println("  End of Record", values) // TODO delete
		}
	}
	return data, nil
}

func maybeStartRecord(data []byte, inRecord *bool, oldFieldIndex,
	fieldIndex, lino *int) ([]byte, error) {
	*inRecord = true
	*oldFieldIndex = -1
	*fieldIndex = 0
	data = skipWs(data, lino)
	if len(data) == 0 {
		return data, fmt.Errorf("e%d#%d:unexpected end of data", e113,
			*lino)
	}
	if data[0] != ']' { // TODO delete
		fmt.Println("  Start of Record")
	}
	return data, nil
}

func handleSentinal(data []byte, metaField *metaFieldType,
	lino *int) ([]byte, *reflect.Value, error) {
	if metaField.kind == boolField {
		return data, nil, fmt.Errorf(
			"e%d#%d:the sentinal is invalid for bool fields", e115, lino)
	}
	// TODO add sentinal value for the current field's type to
	// current record
	value := reflect.ValueOf(0) // TODO use the right type!
	fmt.Println("    Got: !")   // TODO delete
	return data[1:], &value, nil
}

func handleBool(data []byte, value bool, metaField *metaFieldType,
	lino *int) ([]byte, *reflect.Value, error) {
	if metaField.kind != boolField {
		return data, nil, fmt.Errorf("e%d#%d:got bool, expected %s", e114,
			*lino, metaField.kind)
	}
	fmt.Printf("    Got: %t\n", value) // TODO delete
	v := reflect.ValueOf(value)
	return data[1:], &v, nil
}

func handleBytes(data []byte, metaField *metaFieldType, lino *int) ([]byte,
	*reflect.Value, error) {
	data = data[1:] // skip (
	if metaField.kind != bytesField {
		return data, nil, fmt.Errorf("e%d#%d:got bytes, expected %s", e116,
			*lino, metaField.kind)
	}
	data, raw, err := readHexBytes(data, lino)
	if err != nil {
		return data, nil, err
	}
	fmt.Printf("    Got: %q\n", raw) // TODO delete
	value := reflect.ValueOf(raw)
	return data, &value, nil
}

func handleStr(data []byte, metaField *metaFieldType, lino *int) ([]byte,
	*reflect.Value, error) {
	data = data[1:] // skip <
	if metaField.kind != strField {
		return data, nil, fmt.Errorf("e%d#%d:got str, expected %s", e117,
			*lino, metaField.kind)
	}

	data, s, err := readString(data, lino)
	if err != nil {
		return data, nil, err
	}
	fmt.Printf("    Got: %q\n", s) // TODO delete
	value := reflect.ValueOf(s)
	return data, &value, nil
}

func handleInt(data []byte, metaField *metaFieldType, lino *int) ([]byte,
	*reflect.Value, error) {
	data, i, err := readInt(data, lino)
	if err != nil {
		return data, nil, err
	}
	fmt.Printf("    Got: %d\n", i) // TODO delete
	value := reflect.ValueOf(i)
	return data, &value, nil
}

func handleReal(data []byte, metaField *metaFieldType, lino *int) ([]byte,
	*reflect.Value, error) {
	data, r, err := readReal(data, lino)
	if err != nil {
		return data, nil, err
	}
	fmt.Printf("    Got: %g\n", r) // TODO delete
	value := reflect.ValueOf(r)
	return data, &value, nil
}

func handleDateTime(data []byte, format string, metaField *metaFieldType,
	lino *int) ([]byte, *reflect.Value, error) {
	data, d, err := readDateTime(data, format, lino)
	if err != nil {
		return data, nil, err
	}
	fmt.Printf("    Got: %s\n", d.Format(format)) // TODO delete
	value := reflect.ValueOf(d)
	return data, &value, err
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
	x, err := strconv.ParseInt(string(raw), 10, 64)
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
	data, raw, err := scan(data, []byte("-0123456789"), lino)
	if err != nil {
		return data, DateSentinal, err
	}
	x, err := time.Parse(format, string(raw))
	if err != nil {
		return data, DateSentinal, fmt.Errorf("e%d#%d:invalid date/time",
			e127, *lino)
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
