package be

import (
	// "errors"
	// "github.com/bwmarrin/discordgo"
	"log"
	"net/http"
	// "strings"
	"encoding/json"
	"fmt"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	"io/ioutil"
)

var features = [4]string{"profile", "rpinv", "roll", "rpg"}

const ERR_NOGUILDSETTINGS = "There are no settings for this guild"

type BotConfig struct {
	GuildSettings map[string]GuildSettings
}

type GuildSettings struct {
	GuildId string
	// Features        map[string]bool
	FeatureSettings map[string]map[string]string
}

type FeatureSettings struct {
	Enabled bool
}

func InitBotConfig(w http.ResponseWriter, r *http.Request) {
	botConf := BotConfig{}
	botConf.GuildSettings = make(map[string]GuildSettings)

	c, s := SetupConn("guildsettings")
	// See if it exists in db
	Result := []GuildSettings{}
	err := c.Find(nil).All(&Result)
	if err != nil {
		log.Printf("%v", err.Error())
	}

	if len(Result) > 0 {
		// We have retrieved settings and should load them
		for _, value := range Result {
			guildId := value.GuildId
			botConf.GuildSettings[guildId] = value
		}
	} else {
		// There are no saved settings and we should init them
	}
	log.Printf("%+v", botConf)

	botConfJson, err := json.Marshal(botConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", botConfJson)

	defer s.Close()
}

func GetAllGuildConfig(w http.ResponseWriter, r *http.Request) {
	c, s := SetupConn("guildsettings")

	Result := []GuildSettings{}
	err := c.Find(nil).All(&Result)

	// If no results
	if err != nil {
		log.Printf("%v", err.Error())
		fmt.Fprint(w, "{}")
	} else {
		botConfJson, err := json.Marshal(Result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s", botConfJson)
	}

	defer s.Close()
}

// func InitGuildConfig(botConf *BotConfig, guildId string) {
// 	log.Printf("= botadmin:InitGuildConfig: GuildId %s", guildId)
// 	newConfig := GuildSettings{guildId, make(map[string]map[string]string)}
// 	for _, value := range features {
// 		log.Println(value)
// 		newConfig.FeatureSettings[value]["enabled"] = "true"
// 	}
// 	// Default profile settings
// 	newConfig.FeatureSettings["profile"]["creditOneAlias"] = "Splots"
// 	newConfig.FeatureSettings["profile"]["creditTwoAlias"] = "Ink Pots"
// 	newConfig.FeatureSettings["profile"]["creditThreeAlias"] = "Paintbrushes"

// 	botConf.GuildSettings[guildId] = newConfig
// 	SaveGuildConfig(newConfig)
// }

// func SetConfig(botConf BotConfig, guildId string, field string, setTo bool) error {
//  log.Println("= botadmin:SetConfig: Setting %s to %t for guild %s", field, setTo, guildId)
//  // Finds the guild setting object
//  guildSettings, err := findSettingByGuildId(guildId, botConf)
//  if err != nil {
//     if (err == )
//    log.Printf("%v", err)
//    return errors.New("Could not retrieve guild settings: %v", err)
//  }

//  // Sets the field to the value
// }

// func FindSettingsByGuildId(guildId string, botConf *BotConfig) (GuildSettings, error) {
// 	log.Println("= botadmin:FindSettingsByGuildId: GuildId %s", guildId)
// 	if len(botConf.GuildSettings) <= 0 {
// 		return GuildSettings{}, errors.New(ERR_NOGUILDSETTINGS)
// 	}

// 	guild := botConf.GuildSettings[guildId]
// 	log.Println(guild)
// 	if guild.GuildId == "" {
// 		InitGuildConfig(botConf, guildId)
// 	}

// 	return guild, nil
// }

// func GetConfigForGuild(guildId string, field string, botConf *BotConfig) (map[string]string, error) {
// 	log.Println("= botadmin:GetConfigForGuild: GuildId %s", guildId)
// 	guild, err := FindSettingsByGuildId(guildId, botConf)
// 	if err != nil {
// 		return make(map[string]string), errors.New("Couldn't find guild settings")
// 	}

// 	return guild.FeatureSettings[field], nil
// }

// func SaveGuildConfig(guildSettings GuildSettings) error {
func SaveGuildConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("= botadmin:SaveGuildConfig:")
	// Get info from body
	// save the thing
	// response

	thing, _ := ioutil.ReadAll(r.Body)
	log.Printf("%+v", string(thing))

	// c, s := SetupConn("guildsettings")
	// var err error
	// if guildSettings.GuildId != "" {
	// 	err = c.Update(bson.M{"id": guildSettings.GuildId}, guildSettings)
	// } else {
	// 	err = errors.New("No guild id given")
	// }
	// defer s.Close()
}
