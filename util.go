// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

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
