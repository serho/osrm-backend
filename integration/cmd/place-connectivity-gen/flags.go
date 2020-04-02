package main

import "flag"

var flags struct {
	inputPath    string
	outputFolder string
}

func init() {
	flag.StringVar(&flags.inputPath, "i", "", "path for input file in json format")
	flag.StringVar(&flags.outputFolder, "o", "", "path for output folder")
}
