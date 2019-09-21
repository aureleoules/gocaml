package db

import (
	"log"

	"gopkg.in/mgo.v2"
)

// GOCAML Database Access Object
var GOCAML *mgo.Database

var (
	// UsersCollection contains users
	UsersCollection = "users"
)

// Connect to MongoDB
func Connect(server string, database string) {
	session, err := mgo.Dial(server)

	if err != nil {
		log.Fatal(err)
	}

	GOCAML = session.DB(database)
	log.Println("Database authenticated")
}
