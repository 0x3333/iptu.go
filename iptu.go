package main

import (
	"flag"

	"bitbucket.org/terciofilho/iptu.go/importer"
)

func main() {
	// Arguments
	importPtr := flag.String("import", "", "Import IPTU CSV")
	dryRunPtr := flag.Bool("dryrun", false, "Dry run usage, doesn't alter the database")

	flag.Parse()

	// Importer
	if *importPtr != "" {
		importer.Import(*importPtr, *dryRunPtr)
	} else {
		flag.Usage()
	}
}
