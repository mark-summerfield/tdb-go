// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"github.com/mark-summerfield/gset"
	"reflect"
	"time"
)

// These are really constants.
var (
	byteSliceType = reflect.TypeOf([]byte{})
	dateTimeType  = reflect.TypeOf(time.Now())
	reservedWords gset.Set[string]
	emptyBytes    = []byte{}
)

const (
	DateFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02T15:04:05"
	e136str        = "e%d#%d:%s fields don't allow nulls: provide a " +
		"valid %s or change the field's type to %s?"
	e146str = "e%d:can't write null to a not null field: provide " +
		"a valid %s or change the field's type to %s?"
)

const (
	// internal error codes
	// ✔ means tested; ✗ means don't know how to test
	e100 = iota + 100 // ✔
	e101              // ✔
	e102              // ✔
	e103              // ✔
	e104              // ✔
	e105              // ✗
	e106              // ✗
	e107              // ✔
	e108              // ✔
	e109              // ✔
	e110              // ✔
	e111              // ✗
	e112              // ✔
	e113              // ✗
	e114              // ✔
	e115
	e116 // ✔
	e117 // ✔
	e118 // ✔
	e119 // ✗
	e120 // ✔
	e121 // ✔
	e122 // ✔
	e123 // ✔
	e124 // ✔
	e125 // ✔
	e126 // ✔
	e127 // ✔
	e128 // ✗
	e129 // ✔
	e130 // ✔
	e131
	e132
	e133
	e134
	e135
	e136
	e137
	e138
	e139
	e140
	e141
	e142
	e143
	e144
	e145
	e146
)

func init() {
	reservedWords = gset.New("bool", "bytes", "date", "datetime", "int",
		"real", "str")
}
