// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	_ "embed"
)

//go:embed Version.dat
var Version string // This tdb package's version.

// TODO Use same API as std. lib. csv, e.g., NewReader() & NewWriter()
// equiv. to Py load & dump

type Tdb struct {
	TableNames []string          // order of reading & writing from/to file
	Tables     map[string]*Table // key is tablename
}

type Table struct {
	MetaTableType // table name and field names and kinds
	Records       []Record
}

type Record []any
