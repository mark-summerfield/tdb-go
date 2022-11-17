// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
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
		}
		if fieldIndex != oldFieldIndex {
			oldFieldIndex = fieldIndex
			metaField = metaTable.Field(fieldIndex)
		}
		b := data[0]
		data = data[1:]
		switch b {
		case ' ', '\n': // ignore whitespace separators
		case '!':
			fmt.Printf("Got #%d: !\n", fieldIndex) // TODO delete
			// TODO add sentinal value for the current field's type to
			// current record
			fieldIndex += 1
		case 'F':
			fmt.Printf("Got #%d: F\n", fieldIndex) // TODO delete
			if metaField.kind != boolField {
				return emptyBytes, fmt.Errorf("e%d#got bool, expected %s",
					WrongType, metaField.kind)
			}
			// TODO add false to current record
			fieldIndex += 1
		case 'T':
			fmt.Printf("Got #%d: T\n", fieldIndex) // TODO delete
			if metaField.kind != boolField {
				return emptyBytes, fmt.Errorf("e%d#got bool, expected %s",
					WrongType, metaField.kind)
			}
			fieldIndex += 1
			// TODO add true to current record
		case '(':
			data, raw, err = readHexBytes(data)
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
			data, s, err = readString(data)
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
			var i int
			var r float64
			var isInt bool
			data, err = readNegativeNumber(data, &i, &r, &isInt)
			if err != nil {
				return emptyBytes, err
			}
			if isInt {
				if metaField.kind != intField {
					return emptyBytes, fmt.Errorf(
						"e%d#got int, expected %s", WrongType,
						metaField.kind)
				}
				// TODO add int to current record
				fmt.Printf("Got #%d: int %d\n", fieldIndex, i) // TODO delete
			} else {
				if metaField.kind != realField {
					return emptyBytes, fmt.Errorf(
						"e%d#got real, expected %s", WrongType,
						metaField.kind)
				}
				// TODO add float to current record
				fmt.Printf("Got #%d: real %f\n", fieldIndex, r) // TODO delete
			}
			fieldIndex += 1
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			////////////////////////////////////////////////////////////
			// NOTE metaField.kind could be int, real, date, or datetime
			// TODO +ve int or float or date or datetime
			fmt.Printf("TODO #%d: parse +v int|float or date|datetime\n",
				fieldIndex) // TODO delete
			fieldIndex += 1
		case ']': // end of table
			if fieldIndex < fieldCount {
				return emptyBytes, fmt.Errorf(
					"e%d#incomplete record %d/%d fields", IncompleteRecord,
					fieldIndex+1, fieldCount)
			} else {
				fmt.Println("End of Table") // TODO delete
				return data, nil
			}
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
	return emptyBytes, fmt.Errorf("e%d#missing table termiator '['",
		MissingTableTerminator)
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

func readNegativeNumber(data []byte, i *int, r *float64,
	isInt *bool) ([]byte, error) {
	end := 0
loop:
	for end < len(data) {
		b := data[end]
		switch b {
		case '-', '+', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case '.', 'e', 'E':
			*isInt = false
		default:
			break loop
		}
		end++
	}
	if end <= 1 { // nothing (shouldn't happen) or bare - with no digits
		return data, fmt.Errorf("e%d#invalid number", InvalidNumber)
	}
	end++
	n := string(data[:end])
	if *isInt {
		x, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			return data, fmt.Errorf("e%d#invalid int", InvalidInt)
		}
		*i = int(-x)
	} else {
		x, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return data, fmt.Errorf("e%d#invalid real", InvalidReal)
		}
		*r = -x
	}
	return data[end:], nil
}
