package rpcmd

import (
	// "encoding/json"
	// "fmt"
	"../be"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"strings"
	// "strconv"
	// "net/url"
)

func GiveHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("=== Accessing GiveHandler")

	err := r.ParseForm()
	if err != nil {
		log.Println("GiveHandler: failed to pass post form")
		log.Fatal(err)
	}

	userId := r.Form.Get("userid")
	itemParamsString := r.Form.Get("itemparams")

	log.Printf("\nUserId: %+v\n", userId)
	log.Printf("\nItemParams: %+v\n", itemParamsString)

	itemParams := strings.Split(itemParamsString, "-")

	// Split into a map so that i can toss it into db
	var itemMap = make(map[string]string)
	prevKey := ""
	itemValue := ""
	for index, element := range itemParams {
		// check if there's a :
		splitParam := strings.Split(element, ":")
		if splitParam == nil {
			if index == 0 {
				log.Println("Error: No key in first itemParam")
				log.Printf("\nitemParams: \n%+v\n", itemParams)
			} else {
				// Add onto the previous
				itemMap[prevKey] = itemMap[prevKey] + element
			}
		} else {
			itemValue = splitParam[1]
			if len(splitParam) > 2 {
				for i := 2; i < len(splitParam); i++ {
					itemValue = itemValue + splitParam[i]
				}
			}
			itemMap[splitParam[0]] = strings.Trim(itemValue, " ")
			prevKey = splitParam[0]
		}
	}
	log.Printf("itemMap:\n%+v\n", itemMap)

	// Start tossing stuff into the db
	c, s := setupConn("users")
	defer s.Close()
}
