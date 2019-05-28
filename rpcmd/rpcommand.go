package rpcmd

import (
	"../be"
	// "gopkg.in/mgo.v2"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"strconv"
	"strings"
	// "net/url"
)

const itemTableName = "rpcmditems"

type UserItemData struct {
	UserId  string              `bson:"userid" json:"userid"`
	GuildId string              `bson:"guildid" json:"guildid"`
	Items   []map[string]string `bson:"items" json:"items"`
}

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("=== rpcmd: Accessing CheckHandler")
	c, s := be.SetupConn(itemTableName)
	defer s.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println("CheckHandler: failed to pass post form")
		log.Printf("Error: %+v", err.Error())
		fmt.Fprint(w, "Error check 1")
		return
	}

	userId := r.Form.Get("userid")
	guildId := r.Form.Get("guildid")
	userItemRecord := UserItemData{}
	if userId == "" {
		log.Println("Error: No userid given or parsed")
		fmt.Fprint(w, "Error check 2")
		return
	}

	log.Printf("Looking for user: %s", userId)
	err = c.Find(bson.M{"userid": userId, "guildid": guildId}).One(&userItemRecord)
	if err != nil && err.Error() != "not found" {
		log.Println("Error: Something went wrong trying to find user in collection")
		log.Printf("%+v", err.Error())
		fmt.Fprint(w, "Error check 3")
		return
	} else if err != nil && err.Error() == "not found" {
		// Create new user with blank items
		itemArray := make([]map[string]string, 0)
		newRecord := &UserItemData{userId, guildId, itemArray}
		c.Insert(&newRecord)
		log.Println("User not found. Creating new user")
		err = c.Find(bson.M{"userid": userId, "guildid": guildId}).One(&userItemRecord)
		if err != nil {
			log.Println("Error: Can't find user that's just been added")
			log.Printf("%+v", err.Error())
			fmt.Fprint(w, "Error check 4")
			return
		}
	}

	jsonResult, err := json.Marshal(userItemRecord)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Sending this: %s", jsonResult)
	fmt.Fprintf(w, "%s", jsonResult)
}

func DiscardHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("=== rpcmd: Accessing DiscardHandler")
	err := r.ParseForm()
	if err != nil {
		log.Println("DiscardHandler: failed to pass post form")
		log.Printf("Error: %+v", err.Error())
		fmt.Fprint(w, "Could not parse info")
		return
	}

	userId := r.Form.Get("userid")
	guildId := r.Form.Get("guildid")
	itemIdString := r.Form.Get("itemid")

	log.Printf("\nUserId: %+v\n", userId)
	log.Printf("\nguildId: %+v\n", guildId)
	log.Printf("\nitemId: %+v\n", itemIdString)

	// find user in that guild
	c, s := be.SetupConn(itemTableName)
	defer s.Close()

	userItemRecord := UserItemData{}
	if userId == "" {
		log.Println("Error: No user id")
		fmt.Fprint(w, "No user id")
		return
	}

	err = c.Find(bson.M{"userid": userId, "guildid": guildId}).One(&userItemRecord)
	if err != nil {
		log.Printf("Error: %+v", err.Error())
		fmt.Fprint(w, "Error finding that user")
		return
	}

	// Check if id exceeds length
	itemNumber, _ := strconv.Atoi(itemIdString[2:])
	log.Printf("converting: %v", itemIdString[2:])
	log.Printf("itemNumber: %v", itemNumber)
	if len(userItemRecord.Items) < itemNumber || itemNumber <= 0 {
		log.Println("Error: ItemId is longer than inv length or below 0")
		fmt.Fprint(w, "Itemid does not exist")
		return
	}
	// delete item
	i := 1
	newLength := len(userItemRecord.Items) - 1
	log.Printf("items length: %i new length: %i", len(userItemRecord.Items), newLength)
	newItemArray := make([]map[string]string, 0)
	log.Printf("Items before: %+v", userItemRecord.Items)
	for key, value := range userItemRecord.Items {
		log.Printf("key: %+v", key)
		if value["itemid"] != itemIdString {
			log.Printf("counter: " + strconv.Itoa(i))
			log.Printf("itemid was: " + value["itemid"])

			value["itemid"] = "ID" + strconv.Itoa(i)

			log.Printf("itemid now: " + value["itemid"])
			newItemArray = append(newItemArray, value)
			i += 1
		}
	}
	log.Printf("Items after: %+v", newItemArray)
	userItemRecord.Items = newItemArray
	c.Update(bson.M{"userid": userId, "guildid": guildId}, &userItemRecord)
	fmt.Fprint(w, "OK")

	// shuffle inventory up
}

func GiveHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("=== rpcmd: Accessing GiveHandler")

	err := r.ParseForm()
	if err != nil {
		log.Println("GiveHandler: failed to pass post form")
		log.Printf("Error: %+v", err.Error())
		fmt.Fprint(w, "Error give 1")
		return
	}

	userId := r.Form.Get("userid")
	guildId := r.Form.Get("guildid")
	itemParamsString := r.Form.Get("itemparams")

	log.Printf("\nUserId: %+v\n", userId)
	log.Printf("\nguildId: %+v\n", guildId)
	log.Printf("\nItemParams: %+v\n", itemParamsString)

	itemParams := strings.Split(itemParamsString, "-")

	// Split into a map so that i can toss it into db
	var itemMap = make(map[string]string)
	prevKey := ""
	itemValue := ""
	trimmedKey := ""
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
			trimmedKey = strings.Trim(splitParam[0], " ")
			itemMap[trimmedKey] = strings.Trim(itemValue, " ")
			prevKey = trimmedKey
		}
	}
	log.Printf("itemMap:\n%+v\n", itemMap)

	// Start tossing stuff into the db
	c, s := be.SetupConn(itemTableName)
	defer s.Close()

	userItemRecord := UserItemData{}
	if userId != "" {
		err := c.Find(bson.M{"userid": userId, "guildid": guildId}).One(&userItemRecord)
		if err != nil && err.Error() != "not found" {
			log.Println("There was an error")
			log.Printf("error: %+v", err.Error())
			fmt.Fprint(w, "Error give 2")
			return
		} else if err != nil && err.Error() == "not found" {
			log.Println("User not found. Creating new record")
			itemArray := make([]map[string]string, 1)
			itemMap["itemid"] = "ID1"
			itemArray[0] = itemMap
			newRecord := &UserItemData{userId, guildId, itemArray}
			c.Insert(&newRecord)

		} else if err == nil {
			log.Println("Found record.")
			log.Printf("Record: %+v", userItemRecord)
			itemArray := userItemRecord.Items
			inventoryLength := len(itemArray)
			itemMap["itemid"] = "ID" + strconv.Itoa(inventoryLength+1)
			itemArray = append(itemArray, itemMap)
			userItemRecord.Items = itemArray
			c.Update(bson.M{"userid": userId, "guildid": guildId}, &userItemRecord)
		}
	}

}
