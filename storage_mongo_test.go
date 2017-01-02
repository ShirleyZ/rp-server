package main

import (
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	"log"
	"testing"
)

func TestBasic(t *testing.T) {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	log.Println("=== Session created")
	// defer session.Close()
	// db, err = sql.Open("mysql", os.Args[2])
	c := session.DB("testDB").C("animals")
	log.Println("=== Collection selected")

	log.Println("=== Looking for records")
	result := Animal{}
	log.Println("capture container initialised")

	// query := bson.M{"name": "Mary"}
	// query := bson.M{}
	err = c.Find(nil).One(&result)
	log.Printf("Found %v\n", &result)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("=== Adding records")
	// storeMDB.AddRecord()
	// ab := Animal{name: "Abbey", species: "Alpaca"}
	// log.Printf("\n%v\n", ab)

	// err = c.Insert(&Animal{name: "Abbey", species: "Alpaca"})
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
