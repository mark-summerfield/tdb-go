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
		return fmt.Errorf("e%d#data holds invalid Tdb text", InvalidTdb)
	}
	dbPtr := reflect.ValueOf(db)
	if dbPtr.Kind() != reflect.Ptr {
		return fmt.Errorf("e%d#target interface must be a pointer",
			InvalidInterface)
	}
	dbVal := dbPtr.Elem()
	if dbVal.Kind() != reflect.Struct {
		return fmt.Errorf(
			"e%d#target interface must be a pointer to a struct",
			InvalidPointerTarget)
	}
	metaData := make(metaDataType)
	var err error
	var tableName string
	for len(data) > 0 { // TODO refactor?
		b := data[0]
		data = data[1:]
		if b == '[' {
			if data, tableName, err = readTableMetaData(data,
				metaData); err != nil {
				fmt.Println("@@@@", metaData) // TODO delete
				return err
			}
		} else {
			fmt.Println("Start of Table", tableName) // TODO delete
			if data, err = readRecords(data,
				metaData[tableName]); err != nil {
				fmt.Println("$$$$", metaData) // TODO delete
				return err
			}
			fmt.Println("End of Table", tableName) // TODO delete
		}
	}
	fmt.Println("####", metaData) // TODO delete
	return nil                    // TODO
}

func readTableMetaData(data []byte,
	metaData metaDataType) ([]byte, string, error) {
	end := bytes.IndexByte(data, '%')
	if end == -1 {
		return data, "", fmt.Errorf(
			"e%d#invalid table definition (missing %%)", InvalidTableDef)
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
			if fieldName == "" {
				return data, "", fmt.Errorf("e%d#missing fieldname or type",
					MissingFieldNameOrType)
			}
			typename := string(part)
			metaTable := metaData[tableName]
			if err := metaTable.Add(fieldName, typename); err != nil {
				return data, "", fmt.Errorf(
					"e%d#invalid typename %s", InvalidTypeName, typename)
			}
		}
	}
	return data[end+1:], tableName, nil // +1 skips final %
}

// TODO take in a reflect.Value for the outer target struct's corresponding
// slice
// TODO refactor
func readRecords(data []byte, metaTable *metaTableType) ([]byte, error) {
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
			fmt.Printf("Start of Record of %d fields\n", fieldCount) // TODO delete
			data = skipWs(data)
			if len(data) == 0 {
				return emptyBytes, fmt.Errorf("e%d#unexpected end of data",
					UnexpectedEndOfData)
			}
		}
		if fieldIndex != oldFieldIndex {
			oldFieldIndex = fieldIndex
			metaField = metaTable.Field(fieldIndex)
			// TODO delete...
			m := len(data)
			if m > 10 {
				m = 10
			}
			fmt.Printf("Expecting %s start=%q\n", metaField.kind, data[:m])
			// TODO end delete
		}
		b := data[0]
		switch b {
		case ' ', '\t', '\n', '\r': // ignore whitespace separators
			data = data[1:]
		case '!':
			fmt.Printf("Got #%d: !\n", fieldIndex) // TODO delete
			// TODO add sentinal value for the current field's type to
			// current record
			fieldIndex += 1
			data = data[1:]
		case 'F', 'f', 'N', 'n':
			fmt.Printf("Got #%d: F\n", fieldIndex) // TODO delete
			if metaField.kind != boolField {
				return emptyBytes, fmt.Errorf("e%d#got bool, expected %s",
					WrongType, metaField.kind)
			}
			// TODO add false to current record
			fieldIndex += 1
			data = data[1:]
		case 'T', 't', 'Y', 'y':
			fmt.Printf("Got #%d: T\n", fieldIndex) // TODO delete
			if metaField.kind != boolField {
				return emptyBytes, fmt.Errorf("e%d#got bool, expected %s",
					WrongType, metaField.kind)
			}
			fieldIndex += 1
			data = data[1:]
			// TODO add true to current record
		case '(':
			data, raw, err = readHexBytes(data[1:])
			if err != nil {
				return emptyBytes, err
			}
			if metaField.kind != bytesField {
				return emptyBytes, fmt.Errorf("e%d#got bytes, expected %s",
					WrongType, metaField.kind)
			}
			fmt.Printf("Got #%d: %v\n", fieldIndex, raw) // TODO delete
			// TODO add raw to current record
			fieldIndex += 1
		case '<':
			data, s, err = readString(data[1:])
			if err != nil {
				return emptyBytes, err
			}
			if metaField.kind != strField {
				return emptyBytes, fmt.Errorf("e%d#got str, expected %s",
					WrongType, metaField.kind)
			}
			fmt.Printf("Got #%d: %q\n", fieldIndex, s) // TODO delete
			// TODO add string to current record
			fieldIndex += 1
		case '-':
			// TODO surely we know whether to expect an int or real based on
			// the metaData?
			switch metaField.kind {
			case intField:
				var i int
				data, i, err = readInt(data)
				if err != nil {
					return emptyBytes, err
				}
				i = -i
				fmt.Printf("Got #%d: %d\n", fieldIndex, i) // TODO delete
			// TODO add int to current record
			case realField:
				var r float64
				data, r, err = readReal(data)
				if err != nil {
					return emptyBytes, err
				}
				r = -r
				fmt.Printf("Got #%d: %f\n", fieldIndex, r) // TODO delete
			// TODO add real to current record
			default:
				return emptyBytes, fmt.Errorf(
					"e%d#got -, expected %s", WrongType,
					metaField.kind)
			}
			fieldIndex += 1
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			switch metaField.kind {
			case intField:
				var i int
				data, i, err = readInt(data)
				if err != nil {
					return emptyBytes, err
				}
				fmt.Printf("Got #%d: %d\n", fieldIndex, i) // TODO delete
			// TODO add int to current record
			case realField:
				var r float64
				data, r, err = readReal(data)
				if err != nil {
					return emptyBytes, err
				}
				fmt.Printf("Got #%d: %f\n", fieldIndex, r) // TODO delete
				// TODO add real to current record
			case dateField:
				var d time.Time
				data, d, err = readDateTime(data, DateFormat)
				if err != nil {
					return emptyBytes, err
				}
				fmt.Printf("Got #%d: %s\n", fieldIndex, d.Format(DateFormat)) // TODO delete
			case dateTimeField:
				var d time.Time
				data, d, err = readDateTime(data, DateTimeFormat)
				if err != nil {
					return emptyBytes, err
				}
				fmt.Printf("Got #%d: %s\n", fieldIndex, d.Format(DateTimeFormat)) // TODO delete
			default:
				return emptyBytes, fmt.Errorf(
					"e%d#got -, expected %s", WrongType,
					metaField.kind)
			}
			fieldIndex += 1
		case ']': // end of table
			if fieldIndex > 0 && fieldIndex < fieldCount {
				return emptyBytes, fmt.Errorf(
					"e%d#incomplete record %d/%d fields", IncompleteRecord,
					fieldIndex+1, fieldCount)
			}
			fmt.Println("End of Table") // TODO delete
			return skipWs(data[1:]), nil
		default:
			return emptyBytes, fmt.Errorf("e%d#invalid character %q",
				InvalidCharacter, rune(b))
		}
		if fieldIndex == fieldCount {
			inRecord = false
			// TODO add field
			fmt.Println("End of Record") // TODO delete
		}
	}
	return data, nil
}

func readHexBytes(data []byte) ([]byte, []byte, error) {
	end := bytes.IndexByte(data, ')')
	if end == -1 {
		return data, nil, fmt.Errorf("e%d#missing bytes terminator ')'",
			MissingBytesTerminator)
	}
	chunks := bytes.Fields(data[:end])
	chunk := bytes.Join(chunks, emptyBytes)
	raw := make([]byte, hex.DecodedLen(len(chunk)))
	_, err := hex.Decode(raw, chunk)
	if err != nil {
		return emptyBytes, nil, fmt.Errorf("e%d#invalid bytes %q",
			InvalidBytes, chunk)
	}
	return data[end+1:], raw, nil // +1 skips final )
}

func readString(data []byte) ([]byte, string, error) {
	end := bytes.IndexByte(data, '>')
	if end == -1 {
		return data, "", fmt.Errorf("e%d#missing string terminator '>'",
			MissingStringTerminator)
	}
	s := Unescape(string(data[:end]))
	return data[end+1:], s, nil // +1 skips final >
}

func readInt(data []byte) ([]byte, int, error) {
	data, raw, err := scan(data, []byte("-+0123456789"))
	if err != nil {
		return data, 0, err
	}
	x, err := strconv.ParseInt(string(raw), 10, 64)
	if err != nil {
		return data, 0, fmt.Errorf("e%d#invalid int", InvalidInt)
	}
	return data, int(x), nil
}

func readReal(data []byte) ([]byte, float64, error) {
	data, raw, err := scan(data, []byte("-+0123456789.eE"))
	if err != nil {
		return data, 0, err
	}
	x, err := strconv.ParseFloat(string(raw), 64)
	if err != nil {
		return data, 0, fmt.Errorf("e%d#invalid real", InvalidReal)
	}
	return data, x, nil
}

func readDateTime(data []byte, format string) ([]byte, time.Time, error) {
	data, raw, err := scan(data, []byte("-0123456789"))
	if err != nil {
		return data, DateSentinal, err
	}
	x, err := time.Parse(format, string(raw))
	if err != nil {
		return data, DateSentinal, fmt.Errorf("e%d#invalid date/time",
			InvalidDate)
	}
	return data, x, nil
}

func scan(data, valid []byte) ([]byte, []byte, error) {
	data = skipWs(data)
	end := 0
	for end < len(data) {
		b := data[end]
		if bytes.IndexByte(valid, b) == -1 { // end of search
			return data[end:], data[:end], nil
		}
		end++
	}
	return emptyBytes, emptyBytes, fmt.Errorf("e%d#unexpected end of data",
		UnexpectedEndOfData)
}

func skipWs(data []byte) []byte {
	end := 0
	for end < len(data) {
		b := data[end]
		if bytes.IndexByte([]byte{' ', '\t', '\n', '\r'}, b) == -1 {
			return data[end:]
		}
		end++
	}
	return data
}
