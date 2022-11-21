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
	e100 = iota + 100 // internal error codes
	e101
	e102
	e103
	e104
	e105
	e106
	e107
	e108
	e109
	e110
	e111
	e112
	e113
	e114
	e115
	e116
	e117
	e118
	e119
	e120
	e121
	e122
	e123
	e124
	e125
	e126
	e127
)

func init() {
	reservedWords = gset.New("bool", "bytes", "date", "datetime", "int",
		"real", "str")
}
