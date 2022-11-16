// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import "time"

var (
	BytesSentinal    = []byte{ByteSentinal}
	DateSentinal     = time.Date(1808, time.August, 8, 0, 0, 0, 0, time.UTC)
	DateTimeSentinal = time.Date(1808, time.August, 8, 8, 8, 8, 0, time.UTC)
)

const (
	ByteSentinal        byte = 0x04
	DateStrSentinal          = "1808-08-08"
	DateTimeStrSentinal      = "1808-08-08T08:08:08"
	IntSentinal              = -1808080808
	RealSentinal             = -1808080808.0808
	StrSentinal              = "\x04"
)
