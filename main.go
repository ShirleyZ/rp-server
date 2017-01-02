package main

import (
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	"./be"
	"./fe"
	"log"
	"net/http"
)

type Animal struct {
	name    string
	species string
	age     int
	rating  string
}

func init() {
	// Booting onto localhost post
	// log.Println("**********************************")
	// log.Println("*** BOOT UP SEQUENCE INITIATED ***")
	// log.Println("**********************************")

	// // Opening connection to db
	// log.Println("== Opening MDB connection")
	// session, err := mgo.Dial("localhost:27017")
	// // s := be.DataStore{session}
	// c := session.DB("testDB").C("animals")

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// session.SetMode(mgo.Monotonic, true)
	// Deferring connection close
	// defer session.Close()
	// log.Println("**********************************")
	// log.Println("*** BOOT UP SEQUENCE COMPLETED ***")
	// log.Println("**********************************")
}

func main() {
	log.Println("**********************************")
	log.Println("***** LAUNCHING MAIN PROGRAM *****")
	log.Println("**********************************")

	// ==== Front end calls
	log.Println("== Handling front end")
	http.HandleFunc("/", fe.HomeHandler)

	// ==== API calls
	log.Println("== Handling API")
	http.HandleFunc("/api/find/", be.FindHandler)
	http.HandleFunc("/api/add/user/", be.AddUserHandler)
	http.HandleFunc("/api/credits/add/", be.CreditsAddToUserHandler)

	log.Println("== Listen and Serve on port")
	http.ListenAndServe(":8080", nil)
}

// func AddRecord(data, db) {
func addRecord() {
	log.Print("Hi")

}
