// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import "time"

var (
	BytesSentinal = []byte{0x04}
	StrSentinal   = "\x04"
	DateSentinal  = time.Date(1808, time.August, 8, 8, 8, 8, 0, time.UTC)
)

const (
	BoolSentinal = false
	IntSentinal  = -1808080808
	RealSentinal = -1808080808.0808
)
