// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package main

import (
	"fmt"
	//"github.com/mark-summerfield/tdb"
	"github.com/mark-summerfield/clip"
	"strings"
)

func main() {
	config := getConfig()
	var tables Tables
	if config.infile == "-" || strings.HasSuffix(config.infile, ".tdb") ||
		strings.HasSuffix(config.infile, ".tdb.gz") {
		tables = readTdb(config.infile)
	} else if strings.HasSuffix(config.infile, ".csv") {
		tables = readCsv(config.infile)
	}
	if config.outfile == "-" || strings.HasSuffix(config.outfile, "tdb") ||
		strings.HasSuffix(config.outfile, ".tdb.gz") {
		writeTdb(config.outfile, config.decimals, tables)
	} else if strings.HasSuffix(config.outfile, ".json") {
		writeJson(config.outfile, config.decimals, tables)
	} else if strings.HasSuffix(config.outfile, "uxf") ||
		strings.HasSuffix(config.outfile, ".uxf.gz") {
		writeUxf(config.outfile, config.decimals, tables)
	}
}

func getConfig() Config {
	parser := clip.NewParser()
	parser.LongDesc = "Converts Tdb or CSV input to CSV, JSON, SQLite " +
		"Tdb, UXF, or XML."
	parser.PositionalCount = clip.TwoPositionals
	parser.PositionalHelp = "FILE1 is - for stdin or a .csv, .tdb, or " +
		".tdb.gz file. FILE2 is - for stdout (in Tdb format), or a " +
		".json, .tdb, .tdb.gz, .uxf, .uxf.gz, or .xml file."
	decimalsOpt := parser.IntInRange("decimals", "How many decimal digits "+
		"to use. Range 1-19 or 0 (few as possible; the default).", 1, 19, 0)
	if err := parser.Parse(); err != nil {
		fmt.Println(err)
	}
	return Config{decimalsOpt.Value(), parser.Positionals[0],
		parser.Positionals[1]}
}

func readTdb(filename string) Tables {
	fmt.Println("readTdb", filename) // TODO
	return nil
}

func readCsv(filename string) Tables {
	fmt.Println("readCsv", filename) // TODO
	return nil
}

func writeTdb(filename string, dp int, tables Tables) {
	fmt.Println("writeTdb", filename, dp, tables) // TODO
}

func writeJson(filename string, dp int, tables Tables) {
	fmt.Println("writeJson", filename, dp, tables) // TODO
}

func writeUxf(filename string, dp int, tables Tables) {
	fmt.Println("writeUxf", filename, dp, tables) // TODO
}

type Tables [][]any

type Config struct {
	decimals int
	infile   string
	outfile  string
}

func (me Config) String() string {
	return fmt.Sprintf("{decimals: %d, infile: %q, outfile: %q}",
		me.decimals, me.infile, me.outfile)
}
