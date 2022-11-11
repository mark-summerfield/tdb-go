// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	_ "embed"
	//"fmt"
	//"reflect"
)

//go:embed Version.dat
var Version string

const TdbVersion = "1"

// Marshal converts the given struct to a string (as raw UTF-8-encoded
// bytes) in Tdb format.
func Marshal(v any) ([]byte, error) {
	// The format to use is: [ TDEF \n % \n (ROW \n)* ] \n
	return nil, nil // TODO
}

// Unmarshal reads the data from the given string (as raw UTF-8-encoded
// bytes) into a struct.
func Unmarshal(data []byte, v any) error {
	// TODO
	return nil // TODO
}
