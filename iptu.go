package main

import (
	"database/sql"
	"flag"

	"bitbucket.org/terciofilho/iptu.go/importer"
	"bitbucket.org/terciofilho/iptu.go/server"
)

func main() {
	// Arguments
	importPtr := flag.String("import", "", "Import IPTU CSV")
	serverPtr := flag.Bool("server", false, "Start a WebServer to handle requests and serve static resources")
	dryRunPtr := flag.Bool("dryrun", false, "Dry run usage, doesn't alter the database")
	flag.Parse()

	// Connect to database
	db := connectDb()

	// Importer
	if *serverPtr {
		server.Server(db)
	} else if *importPtr != "" {
		importer.Import(db, *importPtr, *dryRunPtr)
	} else {
		flag.Usage()
	}
}

func connectDb() *sql.DB {
	print("Connecting to DB... ")
	db, err := sql.Open("mysql", "iptu:iptu@/iptu?autocommit=false")
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	println(" Connected to DB!")
	return db
}
