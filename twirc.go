package twirc

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/thoj/go-ircevent"
	"strings"
	"time"
)

var (
	Conf Config
	db   *bolt.DB
)

const DB_STEAM_ID string = "steam_id"

type (
	Config struct {
		// The Twitch IRC server, shouldn't need to change this
		Server string

		// Twitch username
		Name string

		// Twitch oauth IRC key
		// See: http://www.twitchapps.com/tmi/
		Password string

		// Steam API Key
		// See: http://steamcommunity.com/dev/apikey
		ApiKey string

		// Streamers steam ID, used as defaults for some commands
		SteamID string

		// Join these channels automatically
		AutoJoin []string

		// Send error and debugging messages to this channel
		// You probably want to use a different channel from your main
		DebugChannel string

		// Print callback handler debug info to the console
		VerboseCallbackHandler bool

		// Print debug irc info to the console
		Debug bool

		// Database file path to use
		Database string
	}
)

func New(config Config) (*irc.Connection, error) {
	irc_conn := irc.IRC(config.Name, config.Name)
	irc_conn.VerboseCallbackHandler = config.VerboseCallbackHandler
	irc_conn.Debug = config.Debug
	irc_conn.Password = config.Password

	irc_conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		fields := strings.Fields(strings.ToLower(e.Message()))

		if fields[0] == "~halt" {
			irc_conn.Quit()
		}
		if fields[0] == "!steamid" {
			if len(fields) != 2 {
				irc_conn.Privmsg(e.Arguments[0], "[SteamID] Must supply a vanity name to resolve")
				return
			}

			steam_id, err := ResolveVanity(fields[1])
			if err != nil {
				irc_conn.Privmsgf(e.Arguments[0], "[SteamID] Error trying to resolve %s", fields[1])
				return
			}
			time.Sleep(1 * time.Second)
			irc_conn.Privmsgf(e.Arguments[0], "[SteamID] %s => %s", fields[1], steam_id)
		}
		if fields[0] == "!setsteamid" {
			time.Sleep(1 * time.Second)
			if len(fields) != 2 {
				irc_conn.Privmsg(e.Arguments[0], "[SteamID] Command only takes 1 argument, the steamid or vanity name")
				return
			}
			steam_id, err := NewSteamID(fields[1])
			if err != nil {
				irc_conn.Privmsg(e.Arguments[0], "[SteamID] Error resolving steam id")
			} else {
				if err := SetSteamID(e.Nick, steam_id); err != nil {
					irc_conn.Privmsg(e.Arguments[0], "[SteamID] Internal oopsie")
				} else {
					irc_conn.Privmsg(e.Arguments[0], "[SteamID] Set steam id successfully")
				}
			}
			return
		}

		if fields[0] == "!mysteamid" {
			time.Sleep(1 * time.Second)
			steam_id, err := GetSteamID(strings.ToLower(e.Nick))
			if err != nil {
				irc_conn.Privmsg(e.Arguments[0], "[SteamID] Must set steam id with !setsteamid command first")
			} else {
				irc_conn.Privmsgf(e.Arguments[0], "[SteamID] %s => %s", e.Nick, steam_id)
			}
			return
		}

		if fields[0] == "!profile" {
			time.Sleep(1 * time.Second)
			sid, _ := NewSteamID(Conf.SteamID)
			irc_conn.Privmsg(e.Arguments[0], sid.ProfileURL())
			return
		}

		if fields[0] == "!myprofile" {
			time.Sleep(1 * time.Second)
			steam_id, err := GetSteamID(e.Nick)
			if err != nil {
				irc_conn.Privmsg(e.Arguments[0], "Must first set steam id with !setsteamid <steamid>")
			} else {
				irc_conn.Privmsg(e.Arguments[0], steam_id.ProfileURL())
			}
			return
		}

		if fields[0] == "!ip" {
			time.Sleep(1 * time.Second)
			player_info, err := GetPlayerInfo(Conf.ApiKey, SteamID(Conf.SteamID))
			if err != nil {
				irc_conn.Privmsg(e.Arguments[0], "Could not fetch player data")
			} else {
				if player_info.GameServerIP == "" {
					irc_conn.Privmsgf(e.Arguments[0], "[Game] %s Game info n/a or playing unsupported game.", Conf.Name)
				} else {
					irc_conn.Privmsgf(e.Arguments[0], "[Game] %s - %s", player_info.GameExtraInfo, player_info.GameServerIP)
				}
			}
			return
		}

		if fields[0] == "!mvm" || fields[0] == "!mymvm" {
			var steam_id SteamID
			if fields[0] == "!mvm" {
				sid := Conf.SteamID
				if len(fields) >= 2 {
					sid = fields[1]
				}
				var err error
				steam_id, err = NewSteamID(sid)
				log.Println("Using ID:", steam_id)
				if err != nil {
					irc_conn.Privmsg(e.Arguments[0], err.Error())
					return
				}
			} else {
				sid, err := GetSteamID(e.Nick)
				if err != nil {
					irc_conn.Privmsg(e.Arguments[0], "Must first set steam id with !setsteamid <steamid>")
					return
				} else {
					steam_id = sid
				}
			}

			inv, err := FetchInventory(steam_id)
			if err != nil {
				log.Println(err.Error())
			}

			tours := inv.FindMVMData()

			total_tours := uint64(0)
			for _, tour := range tours {
				total_tours += tour.Tours
			}

			var buffer bytes.Buffer
			buffer.WriteString(fmt.Sprintf("[MvM Info] Tours: %d", total_tours))

			for _, t := range tours {
				buffer.WriteString(" | ")
				buffer.WriteString(t.InfoStr())
			}
			time.Sleep(1 * time.Second)
			irc_conn.Privmsg(e.Arguments[0], buffer.String())
		}
	})

	irc_conn.AddCallback("001", func(e *irc.Event) {
		for _, channel := range config.AutoJoin {
			irc_conn.Join(channel)
		}
	})

	err := irc_conn.Connect(config.Server)
	if err != nil {
		return nil, err
	}

	return irc_conn, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Shutdown() {
	db.Close()
}

func init() {
	if _, err := toml.DecodeFile("config.toml", &Conf); err != nil {
		// handle error
		log.Fatalln(err.Error())
	}
	if Conf.ApiKey == "" {
		log.Fatalln("Steam API Key must be set")
	}
	log.Debugln("Using database file:", Conf.Database)
	db_global, err := bolt.Open(Conf.Database, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalln(err.Error())
	}
	db = db_global
	err = db.Update(func(tx *bolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists([]byte(DB_STEAM_ID))
		return e
	})
	if err != nil {
		log.Fatalln(err.Error())
		db.Close()
	}
}
