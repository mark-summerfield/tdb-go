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
	if len(data) < 11 {
		return fmt.Errorf("e%d#data holds invalid Tdb text", e107)
	}
	dbPtr := reflect.ValueOf(db)
	if dbPtr.Kind() != reflect.Ptr {
		return fmt.Errorf("e%d#target interface must be a pointer", e108)
	}
	dbVal := dbPtr.Elem()
	if dbVal.Kind() != reflect.Struct {
		return fmt.Errorf(
			"e%d#target interface must be a pointer to a struct", e109)
	}
	metaData := make(metaDataType)
	var err error
	var tableName string
	lino := 1
	for len(data) > 0 {
		b := data[0]
		data = data[1:]
		if b == '[' {
			if data, tableName, err = readTableMetaData(data, metaData,
				&lino); err != nil {
				return err
			}
		} else if tableName != "" {
			fmt.Println("Start of Table", tableName) // TODO delete
			if data, err = readRecords(data, metaData[tableName],
				&lino); err != nil {
				return err
			} else {
				fmt.Println("End of Table", tableName) // TODO delete
				tableName = ""
			}
		}
	}
	return nil
}

func readTableMetaData(data []byte, metaData metaDataType,
	lino *int) ([]byte, string, error) {
	end := bytes.IndexByte(data, '%')
	if end == -1 {
		return data, "", fmt.Errorf(
			"e%d#%d:invalid table definition (missing %%)", e110, *lino)
	}
	*lino += bytes.Count(data[:end], []byte{'\n'})
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
			if fieldName == "" {
				return data, "", fmt.Errorf(
					"e%d#%d:missing fieldname or type", e111, *lino)
			}
			typename := string(part)
			metaTable := metaData[tableName]
			if err := metaTable.Add(fieldName, typename); err != nil {
				return data, "", fmt.Errorf(
					"e%d#%d:invalid typename %s", e112, *lino, typename)
			}
		}
	}
	return data[end+1:], tableName, nil // +1 skips final %
}

// TODO take in a reflect.Value for the outer target struct's corresponding
// slice
// TODO refactor
func readRecords(data []byte, metaTable *metaTableType, lino *int) ([]byte,
	error) {
	var err error
	var raw []byte
	var s string
	var metaField *metaFieldType
	fieldCount := metaTable.Len()
	oldFieldIndex := -1
	fieldIndex := 0
	inRecord := false
	for len(data) > 0 {
		if !inRecord {
			inRecord = true
			oldFieldIndex = -1
			fieldIndex = 0
			data = skipWs(data, lino)
			if len(data) == 0 {
				return emptyBytes, fmt.Errorf(
					"e%d#%d:unexpected end of data", e113, *lino)
			}
			if data[0] != ']' { // TODO delete
				fmt.Printf("  Start of Record of %d fields\n", fieldCount)
			}
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
			fmt.Printf("    Got #%d: !\n", fieldIndex) // TODO delete
			// TODO add sentinal value for the current field's type to
			// current record
			fieldIndex += 1
			data = data[1:]
		case 'F', 'f', 'N', 'n':
			fmt.Printf("    Got #%d: F\n", fieldIndex) // TODO delete
			if metaField.kind != boolField {
				err = fmt.Errorf("e%d#%d:got bool, expected %s", e114,
					*lino, metaField.kind)
			} else {
				// TODO add false to current record
				fieldIndex += 1
				data = data[1:]
			}
		case 'T', 't', 'Y', 'y':
			fmt.Printf("    Got #%d: T\n", fieldIndex) // TODO delete
			if metaField.kind != boolField {
				err = fmt.Errorf("e%d#%d:got bool, expected %s", e115,
					*lino, metaField.kind)
			} else {
				// TODO add true to current record
				fieldIndex += 1
				data = data[1:]
			}
		case '(':
			data, raw, err = readHexBytes(data[1:], lino)
			if err == nil && metaField.kind != bytesField {
				err = fmt.Errorf("e%d#%d:got bytes, expected %s", e116,
					*lino, metaField.kind)
			} else {
				fmt.Printf("    Got #%d: %q\n", fieldIndex, raw) // TODO delete
				// TODO add raw to current record
				fieldIndex += 1
			}
		case '<':
			data, s, err = readString(data[1:], lino)
			if err == nil && metaField.kind != strField {
				err = fmt.Errorf("e%d#%d:got str, expected %s", e117, *lino,
					metaField.kind)
			} else {
				fmt.Printf("    Got #%d: %q\n", fieldIndex, s) // TODO delete
				// TODO add string to current record
				fieldIndex += 1
			}
		case '-':
			switch metaField.kind {
			case intField:
				var i int
				data, i, err = readInt(data, lino)
				if err == nil {
					i = -i
					// TODO add int to current record
					fmt.Printf("    Got #%d: %d\n", fieldIndex, i) // TODO delete
				}
			case realField:
				var r float64
				data, r, err = readReal(data, lino)
				if err == nil {
					r = -r
					// TODO add real to current record
					fmt.Printf("    Got #%d: %f\n", fieldIndex, r) // TODO delete
				}
			default:
				err = fmt.Errorf("e%d#%d:got -, expected %s", e118, *lino,
					metaField.kind)
			}
			fieldIndex += 1
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			switch metaField.kind {
			case intField:
				var i int
				data, i, err = readInt(data, lino)
				if err == nil {
					fmt.Printf("    Got #%d: %d\n", fieldIndex, i) // TODO delete
					// TODO add int to current record
				}
			case realField:
				var r float64
				data, r, err = readReal(data, lino)
				if err == nil {
					fmt.Printf("    Got #%d: %f\n", fieldIndex, r) // TODO delete
					// TODO add real to current record
				}
			case dateField:
				var d time.Time
				data, d, err = readDateTime(data, DateFormat, lino)
				if err == nil {
					fmt.Printf("    Got #%d: %s\n", fieldIndex, d.Format(DateFormat)) // TODO delete
					// TODO add date to current record
				}
			case dateTimeField:
				var d time.Time
				data, d, err = readDateTime(data, DateTimeFormat, lino)
				if err == nil {
					fmt.Printf("    Got #%d: %s\n", fieldIndex, d.Format(DateTimeFormat)) // TODO delete
					// TODO add datetime to current record
				}
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
			return emptyBytes, err
		}
		if fieldIndex == fieldCount {
			inRecord = false
			// TODO add record (if we haven't added fields as we go?)
			fmt.Println("  End of Record") // TODO delete
		}
	}
	return data, nil
}

func readHexBytes(data []byte, lino *int) ([]byte, []byte, error) {
	end := bytes.IndexByte(data, ')')
	if end == -1 {
		return data, nil, fmt.Errorf("e%d#%d:missing bytes terminator ')'",
			e122, *lino)
	}
	*lino += bytes.Count(data[:end], []byte{'\n'})
	chunks := bytes.Fields(data[:end])
	chunk := bytes.Join(chunks, emptyBytes)
	raw := make([]byte, hex.DecodedLen(len(chunk)))
	_, err := hex.Decode(raw, chunk)
	if err != nil {
		return emptyBytes, nil, fmt.Errorf("e%d#%d:invalid bytes %q", e123,
			*lino, chunk)
	}
	return data[end+1:], raw, nil // +1 skips final )
}

func readString(data []byte, lino *int) ([]byte, string, error) {
	end := bytes.IndexByte(data, '>')
	if end == -1 {
		return data, "", fmt.Errorf("e%d#%d:missing string terminator '>'",
			e124, *lino)
	}
	*lino += bytes.Count(data[:end], []byte{'\n'})
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
	return emptyBytes, emptyBytes, fmt.Errorf(
		"e%d#%d:unexpected end of data", e128, *lino)
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
