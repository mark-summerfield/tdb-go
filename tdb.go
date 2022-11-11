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

// Marshal converts the given Tdb struct to a string (as raw UTF-8-encoded
// bytes) in Tdb format. The Tdb struct must be provided by you, e.g.:
//
//	type MyTdb struct {
//		Custom string // The Tdb file's custom string (often "")
//		Tables []MyStructType // This is where the table data goes
//	}
func Marshal(v any) ([]byte, error) {
	return nil, nil // TODO
}

// Unmarshal reads the data from the given string (as raw UTF-8-encoded
// bytes) into a result struct that you must provide:
//
//	type MyTdb struct {
//		Custom string // The Tdb file's custom string (often "")
//		Tables []MyStructType // This is where the table data goes
//	}
func Unmarshal(data []byte, v any) error {
	// TODO
	return nil // TODO
}
