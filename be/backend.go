package be

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"strconv"
	// "net/url"
)

type UserData struct {
	Id       string
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
	findId := urlParams.Get("id")
	log.Printf("New user will be: %s", findName)
	log.Printf("Id is: %s\n", findId)

	if findId != "" {
		newUser := &UserData{findId, findName, 50, "A fresh-faced young adventurer", "Newbie"}
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
		log.Print("User not created. No id given.")
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
	findId := urlParams.Get("id")
	log.Printf("Name is: %s\n", findName)
	log.Printf("Id is: %s\n", findId)

	if findName != "" {
		Result := UserData{}
		err := c.Find(bson.M{"id": findId}).One(&Result)

		if err != nil {
			log.Println("No user found")
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
	}
	defer s.Close()
}

func CreditsAddToUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("= Accessing add credits handler")

	collection, session := setupConn("users")
	log.Printf("Url is: %v \n", r.URL)
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Body is: %+v \n", r.Form)

	// urlParams := r.URL.Query()
	urlParams := r.Form
	receiverName := urlParams.Get("username")
	receiverId := urlParams.Get("id")
	log.Printf("\nName\n%v\n", receiverName)
	log.Printf("\nId\n%v\n", receiverId)
	log.Printf("\n%v\n", urlParams.Get("amount"))

	var amount int
	if urlParams.Get("amount") != "" {
		log.Println("Converting amount to int")
		amount, err = strconv.Atoi(urlParams.Get("amount"))
		if err != nil {
			log.Println("Unable to read amount as int")
			log.Fatal(err)
		}
	} else {
		log.Println("Unable to read amount as int")
		log.Fatal(err)
	}

	// If missing arguments
	if receiverId == "" {
		fmt.Fprint(w, "Error: Missing argument username")
	} else {
		// Get the user's current amount of money
		result := UserData{}
		err := collection.Find(bson.M{"id": receiverId}).One(&result)

		if err != nil {
			log.Println("Find error")
			fmt.Fprint(w, "No user with that name found")

		} else {
			// Send off update query to update
			log.Printf("\nBefore Addition\n%+v\n", &result)
			result.Credits += amount
			log.Printf("\nAfter Addition\n%+v\n", &result)

			err = collection.Update(bson.M{"id": receiverId}, &result)
			if err != nil {
				log.Println("can't do it boss")
				log.Fatal(err)
			}
			fmt.Fprint(w, "Success")
		}
	}

	defer session.Close()

}

func ProfileEditHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("= Accessing profile edit handler")

	collection, session := setupConn("users")
	log.Printf("Url is: %v \n", r.URL)
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Body is: %+v \n", r.Form)

	// Get the user's current amount of money
	receiverId := r.Form.Get("id")
	result := UserData{}
	err = collection.Find(bson.M{"id": receiverId}).One(&result)

	if err != nil {
		log.Println("Find error")
		fmt.Fprint(w, "No user with that id found")

	} else {
		// Send off update query to update
		log.Printf("\nBefore Change\n%+v\n", &result)
		result.Profile = r.Form.Get("profile")
		log.Printf("\nAfter Change\n%+v\n", &result)

		err = collection.Update(bson.M{"id": receiverId}, &result)
		if err != nil {
			log.Println("can't do it boss")
			log.Fatal(err)
		}
		fmt.Fprint(w, "Success")
	}
	defer session.Close()
}
