// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package main

import (
	"fmt"
	//"github.com/mark-summerfield/tdb"
	"github.com/mark-summerfield/clip"
)

func main() {
	config := getConfig()
	// TODO if config.show use tview to show the tables and the records in
	// the selected table (see tview postgres eg)
	fmt.Println(config)
}

func getConfig() config {
	parser := clip.NewParser()
	parser.LongDesc = "Converts Tdb or CSV input to CSV, JSON, SQLite " +
		"Tdb, UXF, or XML.\n\nUse -- before the positionals if either is -."
	parser.PositionalCount = clip.TwoPositionals
	parser.PositionalHelp = "FILE1 is - for stdin or a .csv, .tdb, or " +
		".tdb.gz file. FILE2 is - for stdout (in Tdb format), or a " +
		".csv, .json, .sqlite, .tdb, .tdb.gz, .uxf, .uxf.gz, or .xml "+
		"file (FILE2 is ignored if -s, --show is used)."
	showOpt := parser.Flag("show",
		"Show the given .tdb file in a text user interface.")
	decimalsOpt := parser.IntInRange("decimals", "How many decimal digits "+
		"to use. Range 1-19 or 0 (few as possible; the default).", 0, 19, 0)
	if err := parser.Parse(); err != nil {
		fmt.Println(err)
	}
	return config{showOpt.Value(), decimalsOpt.Value(),
		parser.Positionals[0], parser.Positionals[1]}
}

type config struct {
	show     bool
	decimals int
	infile   string
	outfile  string
}
