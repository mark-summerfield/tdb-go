// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

type (
	ErrorCode   int
	WarningCode int
)

const (
	eWrongType       ErrorCode = iota + 100
	eInvalidLength             // 101
	eInvalidDefault            // 102
	eInvalidRef                // 103
	eInvalidDatabase           // 104

	wInvalidNull      WarningCode = iota + 500
	wWrongType                    // 501
	wValueOutOfRange              // 502
	wValueNotAllowed              // 503
	wLengthOutOfRange             // 504
)
