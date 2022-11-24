// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

/*
Tdb provides [Marshal] and [Unmarshal] functions for writing and reading
[]byte slices to or from Tdb “Text DataBase” format.

To use these functions you must provide a populated (for Marshal) or
unpopulated (for Unmarshal) struct. This outer struct represents a text
database. The outer struct must contain one or more public (inner) fields,
each of type slice of struct. Each inner field represents a database table, and each record is represented by an inner field struct.

# Tdb format

Tdb provides a superior alternative to CSV. In particular, Tdb tables are
named and Tdb fields are strictly typed. Also, there is a clear distinction
between field names and data values, and strings respect whitespace
(including newlines) and have no problems with commas, quotes, etc. Perhaps
best of all, a single Tdb file may contain one—or more—tables. See README.md
at https://github.com/mark-summerfield/tdb-go for more about the Tdb format.

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

The tdb package provides constants for each type's sentinal value (except
for “bool“s for which there is no sentinal value).

The [Marshal] and [Unmarshal] examples use these structs:

	type classicDatabase struct {
		Employees   []Employee   `tdb:"emp"`
		Departments []Department `tdb:"dept"`
	}

	type Employee struct {
		EID        int       `tdb:"empno"`
		Name       string    `tdb:"ename"`
		Job        string    `tdb:"job"`
		ManagerID  int       `tdb:"mgr"`
		HireDate   time.Time `tdb:"hiredate:date"`
		Salary     float64   `tdb:"sal"`
		Commission float64   `tdb:"comm"`
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
