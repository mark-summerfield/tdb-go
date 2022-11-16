// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	_ "embed"
)

//go:embed Version.dat
var Version string

const TdbVersion = "1"

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
	/* NOTE look at encoding/{csv,json,xml} to see how they do it, e.g.,
	NewReader, NewWriter ?
	*/
	return nil, nil // TODO
}

// Unmarshal reads the data from the given string (as raw UTF-8-encoded
// bytes) into a (pointer to a) database struct.
func Unmarshal(data []byte, db any) error {
	/* NOTE look at encoding/{csv,json,xml} to see how they do it, e.g.,
	NewReader, NewWriter ?
	*/
	// TODO
	return nil // TODO
}
