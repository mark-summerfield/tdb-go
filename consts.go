// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"github.com/mark-summerfield/gset"
	"reflect"
	"time"
)

var (
	BytesSentinal    = []byte{ByteSentinal}
	DateSentinal     = time.Date(1808, time.August, 8, 0, 0, 0, 0, time.UTC)
	DateTimeSentinal = time.Date(1808, time.August, 8, 8, 8, 8, 0, time.UTC)
	byteSliceType    = reflect.TypeOf(BytesSentinal)
	dateTimeType     = reflect.TypeOf(DateSentinal)
	reservedWords    gset.Set[string]
	emptyBytes       = []byte{}
)

const (
	ByteSentinal        byte = 0x04
	DateStrSentinal          = "1808-08-08"
	DateTimeStrSentinal      = "1808-08-08T08:08:08"
	IntSentinal              = -1808080808
	RealSentinal             = -1808080808.0808
	StrSentinal              = "\x04"
	DateFormat               = "2006-01-02"
	DateTimeFormat           = "2006-01-02T15:04:05"
)

const (
	// internal error codes
	// ✔ means tested; ✗ means don't know how to test
	e100 = iota + 100 // ✔
	e101              // ✔
	e102              // ✔
	e103              // ✔
	e104              // ✔
	e105              // TODO
	e106              // TODO
	e107              // TODO
	e108              // TODO
	e109              // TODO
	e110              // TODO
	e111              // TODO
	e112              // TODO
	e113              // TODO
	e114              // TODO
	e115              // TODO
	e116              // TODO
	e117              // TODO
	e118              // TODO
	e119              // TODO
	e120              // TODO
	e121              // TODO
	e122              // TODO
	e123              // TODO
	e124              // TODO
	e125              // TODO
	e126              // TODO
	e127              // TODO
	e128              // TODO
)

func init() {
	reservedWords = gset.New("bool", "bytes", "date", "datetime", "int",
		"real", "str")
}
