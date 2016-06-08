package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // Just to initialize MySQL Driver

	"bitbucket.org/terciofilho/iptu.go/log"
)

var (
	// Instance is a reference to the Database
	Instance *sql.DB
)

//ConnectDb connects to the Database
func ConnectDb() {
	log.Info.Println("Connecting to Database... ")
	instance, err := sql.Open("mysql", "iptu:iptu@/iptu?autocommit=false")
	if err != nil {
		panic(err.Error())
	}
	err = instance.Ping()
	if err != nil {
		panic(err.Error())
	}
	Instance = instance
	log.Info.Println("Connected to Database!")
}
