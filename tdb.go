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

// Marshal converts the given database struct to a string (as raw
// UTF-8-encoded bytes) in Tdb format. The database struct should have a
// tdb.MetaData field and one or more slices of structs (each slice holding
// a table's records, each struct a record's fields).
func Marshal(db any) ([]byte, error) {
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
// bytes) into a database struct.
func Unmarshal(data []byte, db any) error {
	/* NOTE look at encoding/{csv,json,xml} to see how they do it, e.g.,
	NewReader, NewWriter ?
	*/
	// TODO
	return nil // TODO
}

// Check reads the metadata and data from a database struct and returns a
// (possibly empty) list of strings: each string is a warning (e.g.,
// an int out of range).
func Check(db any) []string {
	// TODO Iterate over all the data tables and for every field value call
	// the corresponding metafield's check() method with the value.
	// TODO In addition, within each table, check for uniqueness where
	// required.
	// TODO In addition, check ref values to see that their referrent
	// exists, is not to themselves, and has the same int or str type as
	// themselves.
	return nil
}
