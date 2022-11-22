// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import _ "embed"

//go:embed Version.dat
var Version string // This tdb package's version.

const TdbVersion = "1" // The highest Tdb format version this package handles.
