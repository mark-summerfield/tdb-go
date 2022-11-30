// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import "bytes"

// Escape returns an XML-escaped string, i.e., where runes are replaced as
// follows: & → &amp;, < → &lt;, > → &gt;.
// See also [Unescape].
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

// Unescape accepts an XML-escaped string and returns a plain text string
// with no escapes, i.e., where substrings are replaced with runes as
// follows: &amp; → &, &lt; → <, &gt; → >.
// See also [Escape].
func Unescape(s string) string {
	result := make([]byte, 0, len(s))
	raw := []byte(s)
	for len(raw) > 0 {
		b := raw[0]
		if b == '&' && len(raw) > 3 {
			raw = raw[1:]
			end := bytes.IndexByte(raw, ';')
			if end > -1 {
				found := string(raw[:end])
				if found == "lt" {
					result = append(result, '<')
				} else if found == "gt" {
					result = append(result, '>')
				} else if found == "amp" {
					result = append(result, '&')
				}
				raw = raw[end+1:]
			}
		} else {
			result = append(result, b)
			raw = raw[1:]
		}
	}
	return string(result)
}
