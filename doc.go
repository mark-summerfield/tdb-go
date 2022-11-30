// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

/*
Tdb provides the [Parse] and [Unmarshal] functions for reading []byte slices
of text in Tdb “Text DataBase” format, and the [Tdb.Write] and [Marshal]
functions for writing to Tdb format.

The [Parse] function creates a [Tdb] object which stores values as type
`any`, so is useful for applications that need to process generic Tdb files.
[Tdb] data is written in Tdb format using the [Tdb.Write] method. However,
if the Tdb file format is known, then it is best to use [Marshal] and
[Unmarshal] since these use the appropriate concrete types (`bool`, `int`,
`string`, and so on).

To use the [Marshal] and [Unmarshal] functions you must provide a populated
(for Marshal) or unpopulated (for Unmarshal) struct. This outer struct
represents a text database. The outer struct must contain one or more public
(inner) fields, each of type slice of struct. Each inner field represents a
database table, and each record is represented by an inner field struct.

# Tdb format

Tdb “Text DataBase” format is a plain text human readable typed database
storage format.

Tdb provides a superior alternative to CSV. In particular, Tdb tables are
named and Tdb fields are strictly typed. Also, there is a clear distinction
between field names and data values, and strings respect whitespace
(including newlines) and have no problems with commas, quotes, etc.
Perhaps best of all, a single Tdb file may contain one—or more—tables.

See README.md at https://github.com/mark-summerfield/tdb-go for more about
the Tdb format.

# Using the tdb package

Import using:

	import tdb "github.com/mark-summerfield/tdb-go"

Types:

	| Tdb Type |  Go Types                  |
	|----------|----------------------------|
	| bool     | bool                       |
	| bytes    | []byte                     |
	| date     | time.Time                  |
	| datetime | time.Time                  |
	| int      | int uint int32 uint32 etc. |
	| real     | float64 float32            |
	| str      | string                     |

Note that for nullable types (e.g., `bool?`, `str?`, etc.) the corresponding
Go type must be a pointer (e.g., `*bool`, `*string`, etc.).

The [Marshal] and [Unmarshal] examples use these structs:

	type classicDatabase struct {
		Employees   []Employee   `tdb:"emp"`
		Departments []Department `tdb:"dept"`
	}

	type Employee struct {
		EID        int       `tdb:"empno"`
		Name       string    `tdb:"ename"`
		Job        string    `tdb:"job"`
		ManagerID  *int      `tdb:"mgr"` // The boss doesn't have a mgr
		HireDate   time.Time `tdb:"hiredate:date"`
		Salary     float64   `tdb:"sal"`
		Commission *float64  `tdb:"comm"` // Most don't get commission
		DeptID     int       `tdb:"deptno"`
	}

	type Department struct {
		DID      int    `tdb:"deptno"`
		Name     string `tdb:"dname"`
		Location string `tdb:"loc"`
	}

Although struct tags are used extensively here, they are only actually
required for two purposes. A tag is needed if a Tdb file's table or field
name is different from the corresponding struct name. And a tag is needed
for time.Time fields if the field is a Tdb `date` field (since the default
is `datetime`). For example, see `db1_test.go` and `csv_test.go` for structs
which work fine despite having few tags.

The order of tables in a Tdb file in relation to the outer struct doesn't
matter. However, the order of fields within a table must match between the
Tdb file's table definition and the corresponding struct.

Naturally, you can use any structs you like that meet tdb's minimum
requirements.
*/
package tdb
