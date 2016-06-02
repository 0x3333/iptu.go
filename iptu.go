package main

import (
	"flag"

	"iptu.go/importer"
)

func main() {
	// Arguments
	importPtr := flag.String("import", "", "Import IPTU CSV")
	flag.Parse()

	// Importer
	if *importPtr != "" {
		importer.Import(*importPtr)
	}
}
