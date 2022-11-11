// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	_ "embed"
	"fmt"
)

//go:embed Version.dat
var Version string

func Hello() string {
	return fmt.Sprintf("Hello tdb v%s", Version)
}
