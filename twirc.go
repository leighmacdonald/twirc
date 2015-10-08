package twirc

import (
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/thoj/go-ircevent"
	"strings"
	"bytes"
	"fmt"
)

var Conf Config


type (
	Config struct {
		Server   string
		Name     string
		Password string
		AutoJoin []string
	}
)

func NewIrcClient(config Config) (*irc.Connection, error) {
	irc_conn := irc.IRC(config.Name, config.Name)
	irc_conn.VerboseCallbackHandler = true
	irc_conn.Debug = true
	irc_conn.Password = config.Password

	irc_conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		fields := strings.Fields(e.Message())

		if fields[0] == "~halt" {
			irc_conn.Quit()
		}
		if fields[0] == "!mvm" {
			steam_id := "76561198084134025"
			if len(fields) >= 2 {
				steam_id = fields[1]
			}
			inv, err := FetchInventory(steam_id)
			if err != nil {
				irc_conn.Privmsg("#roto_", err.Error())
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
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}