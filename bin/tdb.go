// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package main

import (
	"errors"
	"fmt"
	"github.com/mark-summerfield/clip"
	tdb "github.com/mark-summerfield/tdb-go"
	"io"
	"os"
	"strings"
)

func main() {
	config, onError := getConfig()
	inFile, err := os.Open(config.infile)
	if err != nil {
		onError(fmt.Errorf("error #3: failed to open infile %q: %s",
			config.infile, err))
	}
	defer inFile.Close()
	raw, err := io.ReadAll(inFile)
	if err != nil {
		onError(fmt.Errorf("error #4: failed to read infile %q: %s",
			config.infile, err))
	}
	db, err := tdb.Parse(raw)
	if err != nil {
		onError(fmt.Errorf("error #5: failed to parse infile %q: %s",
			config.infile, err))
	}
	var outFile *os.File
	if config.outfile == "-" {
		outFile = os.Stdout
	} else {
		outFile, err = os.OpenFile(config.outfile, os.O_CREATE|
			os.O_WRONLY, 0755)
		if err != nil {
			onError(fmt.Errorf("error #6: failed to open outfile %q: %s",
				config.outfile, err))
		}
		defer outFile.Close()
	}
	err = db.WriteDecimals(outFile, config.decimals)
	if err != nil {
		onError(fmt.Errorf("error #7: failed to write outfile %q: %s",
			config.outfile, err))
	}
}

func getConfig() (config, func(error)) {
	parser := clip.NewParser()
	parser.LongDesc = "Converts Tdb input to Tdb in the standard format."
	parser.PositionalCount = clip.OneOrTwoPositionals
	parser.PositionalHelp = "FILE1 must be a .tdb file. " +
		"If FILE2 is - or not given, output is to stdout; " +
		"otherwise to the given .tdb file."
	decimalsOpt := parser.IntInRange("decimals", "How many decimal digits "+
		"to use. Range 1-19 or 0 (few as possible; the default).", 0, 19, 0)
	if err := parser.Parse(); err != nil {
		fmt.Println(err)
	}
	outfile := "-"
	if len(parser.Positionals) == 2 {
		outfile = parser.Positionals[1]
	}
	config := config{decimalsOpt.Value(), parser.Positionals[0], outfile}
	if !strings.HasSuffix(config.infile, ".tdb") {
		parser.OnError(errors.New(
			"error #1: can only read .tdb files"))
	}
	if !(config.outfile == "-" ||
		strings.HasSuffix(config.outfile, ".tdb")) {
		parser.OnError(errors.New("error #2: can only write Tdb format"))
	}
	return config, parser.OnError
}

type config struct {
	decimals int
	infile   string
	outfile  string
}
