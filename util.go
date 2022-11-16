// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"github.com/mark-summerfield/gong"
	"time"
)

func Escape(s string) string {
	result := make([]rune, 0, len(s)/2)
	for _, c := range s {
		switch c {
		case '&':
			result = append(result, []rune{'&', 'a', 'm', 'p', ';'}...)
		case '<':
			result = append(result, []rune{'&', 'l', 't', ';'}...)
		case '>':
			result = append(result, []rune{'&', 'g', 't', ';'}...)
		default:
			result = append(result, c)
		}
	}
	return string(result)
}

func IsSentinal(value any) bool {
	switch value := value.(type) {
	case bool:
		return !value
	case []byte:
		return len(value) == 1 && value[0] == ByteSentinal
	case time.Time:
		return value.Equal(DateSentinal) || value.Equal(DateTimeSentinal)
	case int8, int16, uint16, uint32, uint64, uint:
		return false
	case int32:
		return int(value) == IntSentinal
	case int64:
		return int(value) == IntSentinal
	case int:
		return value == IntSentinal
	case float32:
		return gong.IsRealClose(float64(value), RealSentinal)
	case float64:
		return gong.IsRealClose(value, RealSentinal)
	case string:
		return value == StrSentinal
	}
	return false
}
