package be

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

type UserData struct {
	Username string
	Credits  int
	Profile  string
	Title    string
}

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("= Accessing add handler")
	c, s := setupConn("users")

	log.Printf("U is: %v \n", r.URL)

	urlParams := r.URL.Query()
	findName := urlParams.Get("name")
	log.Printf("New user will be: %s", findName)

	if findName != "" {
		newUser := &UserData{findName, 50, "A fresh-faced young adventurer", "Newbie"}
		// err := c.Insert(&UserData{findName, 50, "", ""})
		err := c.Insert(newUser)
		if err != nil {
			log.Fatal(err)
		}
		jsonResult, err := json.Marshal(newUser)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s", jsonResult)

	} else {
		log.Print("User not created. No name given.")
	}
	defer s.Close()
}

func setupConn(table string) (c *mgo.Collection, s *mgo.Session) {
	// Creating conenction each time because, i'm a god damn noob that doesn't know better
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	// defer session.Close() // if i do it here, problems
	c = session.DB("rphelper").C(table)

	// Let this misery end pls
	return c, session
}

// Just want to access db so i can find/insert stuff, but it's not being passed through
func FindHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("= Accessing find handler")
	c, s := setupConn("users")

	log.Printf("U is: %v \n", r.URL)
	urlParams := r.URL.Query()
	findName := urlParams.Get("name")
	log.Printf("Name is: %s", findName)

	if findName != "" {
		// Result := []UserData{}
		Result := UserData{}
		// Result := bson.M{}
		err := c.Find(bson.M{"username": findName}).One(&Result)
		if err != nil {
			log.Println("Find error")
			fmt.Fprint(w, "No user with that name found")

		} else {
			log.Println("Result")
			log.Printf("\n%+v\n", Result)

			jsonResult, err := json.Marshal(Result)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, "%s", jsonResult)
		}
		defer s.Close()

	}
}

func CreditsAddToUserHandler(w http.ResponseWriter, r *http.Request) {

}
