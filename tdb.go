// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	_ "embed"
	"fmt"
	"reflect"
)

//go:embed Version.dat
var Version string

const Header = "Tdb1"

// Marshal maps all of the structs in a slice of structs to a string (as raw
// UTF-8-encoded bytes) in Tdb1 format.
func Marshal(v any) ([]byte, error) {
	return nil, nil // TODO
}

// Unmarshal maps the data from given string (as raw UTF-8-encoded bytes)
// into a slice of structs.
func Unmarshal(data []byte, v any) error {
	sliceValuePtr := reflect.ValueOf(v)
	if sliceValuePtr.Kind() != reflect.Ptr {
		return fmt.Errorf(
			"#%d: target isn't a pointer to a slice of structs",
			eNotAPointer)
	}
	sliceValue := sliceValPtr.Elem()
	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("#%d: target isn't a slice",
			eNotASlice)
	}
	structType := sliceVal.Type().Elem()
	if structType.Kind() != reflect.Struct {
		return fmt.Errorf("#%d: target isn't a slice of structs",
			eNotASliceOfStructs)
	}
	size := len(Header)
	if len(data) > size {
		if data[:size] == []byte(Header) { // Strip header if present
			data = data[size:]
		}
	}
	return nil // TODO
}
