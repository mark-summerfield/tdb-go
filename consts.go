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
)

const (
	ByteSentinal        byte = 0x04
	DateStrSentinal          = "1808-08-08"
	DateTimeStrSentinal      = "1808-08-08T08:08:08"
	IntSentinal              = -1808080808
	RealSentinal             = -1808080808.0808
	StrSentinal              = "\x04"
)

const (
	CannotMarshal           = iota + 100 // 100
	CannotMarshalOuter                   // 101
	CannotMarshalEmpty                   // 102
	InvalidSliceType                     // 103
	InvalidFieldType                     // 104
	InvalidSliceFieldType                // 105 NOTE not sure how to test this
	InvalidDateTime                      // 106 ditto
	InvalidTdb                           // 107
	InvalidInterface                     // 108
	InvalidPointerTarget                 // 109
	InvalidTableDef                      // 110
	InvalidBytes                         // 111
	InvalidNumber                        // 112
	InvalidInt                           // 113
	InvalidReal                          // 114
	InvalidCharacter                     // 115
	MissingFieldNameOrType               // 116
	MissingBytesTerminator               // 117
	MissingStringTerminator              // 118
	MissingTableTerminator               // 119
	IncompleteRecord                     // 120
)

func init() {
	reservedWords = gset.New("bool", "bytes", "date", "datetime", "int",
		"real", "str")
}
