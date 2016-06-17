package main

import (
	"flag"

	"bitbucket.org/terciofilho/iptu.go/db"
	"bitbucket.org/terciofilho/iptu.go/importer"
	"bitbucket.org/terciofilho/iptu.go/log"
	"bitbucket.org/terciofilho/iptu.go/server"
	"bitbucket.org/terciofilho/iptu.go/sitemap"
)

func main() {
	// Arguments
	importPtr := flag.String("import", "", "Import IPTU CSV")
	serverPtr := flag.Bool("server", false, "Start a WebServer to handle requests and serve static resources")
	dryRunPtr := flag.Bool("dryrun", false, "Dry run usage, doesn't alter the database")
	sitemapPtr := flag.Bool("sitemap", false, "Generate Sitemap files")
	flag.Parse()

	// Importer
	if *serverPtr {
		log.Info.Println("Starting as a Server...")
		db.ConnectDb()
		server.StartServer()
	} else if *importPtr != "" {
		log.Info.Println("Starting as a Importer...")
		db.ConnectDb()
		importer.RunImport(*importPtr, *dryRunPtr)
	} else if *sitemapPtr {
		log.Info.Println("Starting as Sitemap...")
		db.ConnectDb()
		sitemap.Generate()
	} else {
		flag.Usage()
	}
}
