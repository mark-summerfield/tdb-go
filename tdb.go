// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import _ "embed"

//go:embed Version.dat
var Version string // This tdb package's version.

// DecimalPlaces for Marshal: -1 (or 0) signifies use minimum number of
// places to preserve value, e.g, 5.0 → 5 (this is the default). 1-19 means
// use exactly that number; 20+ means use 19.
//
// If you use this we recommend:
//
//	tdb.DecimalPlaces = 6 // or whatever you want
//	defer func() { tdb.DecimalPlaces = -1 }() // restore when done
var DecimalPlaces = -1

const TdbVersion = "1" // The highest Tdb format version this package handles.
