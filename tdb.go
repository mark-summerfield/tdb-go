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
	// The format to use is:
	// [tablename
	//   fieldname1 constraints1
	//		:
	//   fieldnameN constraintsN
	// %
	// row0field0 ... row0fieldN
	//		:
	// rowMfield0 ... rowMfieldN
	// ]
	/* NOTE look at encoding/{csv,json,xml} to see how they do it, e.g.,
	NewReader, NewWriter ?
	*/
	return nil, nil // TODO
}

// Unmarshal reads the data from the given string (as raw UTF-8-encoded
// bytes) into a struct.
func Unmarshal(data []byte, v any) error {
	/* NOTE look at encoding/{csv,json,xml} to see how they do it, e.g.,
	NewReader, NewWriter ?
	*/
	// TODO
	return nil // TODO
}
