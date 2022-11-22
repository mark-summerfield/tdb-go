// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

/*
Tdb provides Marshal and Unmarshal functions for writing and reading []byte
slices to or from .tdb format.

To use these functions you must provide a populated (for Marshal) or
unpopulated (for Unmarshal) struct. This outer struct represents a text
database. The outer struct must contain one or more public (inner) fields,
each of type slice of struct. Each inner field represents a database table, and each record is represented by an inner field struct.

Import using:

	import tdb "github.com/mark-summerfield/tdb-go"

The Marshal and Unmarshal examples use these structs:

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

Naturally, you can use any structs you like that meet tdb's minimum
requirements.
*/
package tdb
