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
	DateSentinal     = time.Date(1808, time.August, 8, 0, 0, 0, 0, time.UTC)
	DateTimeSentinal = time.Date(1808, time.August, 8, 8, 8, 8, 0, time.UTC)
	byteSliceType    = reflect.TypeOf([]byte{})
	dateTimeType     = reflect.TypeOf(DateSentinal)
	reservedWords    gset.Set[string]
	emptyBytes       = []byte{}
)

const (
	DateStrSentinal     = "1808-08-08"
	DateTimeStrSentinal = "1808-08-08T08:08:08"
	IntSentinal         = -1808080808
	RealSentinal        = -1808080808.0808
	DateFormat          = "2006-01-02"
	DateTimeFormat      = "2006-01-02T15:04:05"
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
	e115              // ✔
	e116              // ✔
	e117              // ✔
	e118              // ✔
	e119              // ✗
	e120              // ✔
	e121              // ✔
	e122              // ✔
	e123              // ✔
	e124              // ✔
	e125              // ✔
	e126              // ✔
	e127              // ✔
	e128              // ✗
	e129              // ✔
	e130              // ✔
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
)

func init() {
	reservedWords = gset.New("bool", "bytes", "date", "datetime", "int",
		"real", "str")
}
