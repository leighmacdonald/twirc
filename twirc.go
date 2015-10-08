package twirc

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/thoj/go-ircevent"
	"strings"
	"time"
)

var Conf Config

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
	}
)

func NewIRCClient(config Config) (*irc.Connection, error) {
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
				irc_conn.Privmsg(e.Arguments[0], "> Must supply a vanity name to resolve")
				return
			}

			steam_id, err := ResolveVanity(fields[1])
			if err != nil {
				irc_conn.Privmsg(e.Arguments[0], fmt.Sprintf("> Error trying to resolve %s", fields[1]))
				return
			}
			time.Sleep(1 * time.Second)
			irc_conn.Privmsg(e.Arguments[0], fmt.Sprintf("> Steam ID: %s => %s", fields[1], steam_id))
		}
		if fields[0] == "!mvm" {
			steam_id := Conf.SteamID
			if len(fields) >= 2 {
				steam_id = fields[1]
			}
			id, err := SteamID(steam_id)
			if err != nil {
				irc_conn.Privmsg(e.Arguments[0], err.Error())
				return
			}
			inv, err := FetchInventory(id)
			if err != nil {
				irc_conn.Privmsg(e.Arguments[0], "Invalid ID supplied")
				irc_conn.Privmsg(Conf.DebugChannel, err.Error())
				return
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

func LoadConfig() {
	if _, err := toml.DecodeFile("config.toml", &Conf); err != nil {
		// handle error
		log.Fatalln(err.Error())
	}
	if Conf.ApiKey == "" {
		log.Fatalln("Steam API Key must be set")
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
