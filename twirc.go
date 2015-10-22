package twirc

import (
	"fmt"
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/thoj/go-ircevent"
	"strings"
	"time"
)

var (
	Conf      Config
	db        *bolt.DB
	stop_chan chan struct{}
	Handlers  map[string]IRCHandlerFunc
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

		// Prefix for IRC commands, "!" default
		Prefix string
	}
)

func (c *Config) ChannelName() string {
	return strings.ToLower(fmt.Sprintf("#%s", c.Name))
}

func _pause() {
	time.Sleep(1 * time.Second)
}

func New(config Config) (*irc.Connection, error) {
	irc_conn := irc.IRC(config.Name, config.Name)
	irc_conn.VerboseCallbackHandler = config.VerboseCallbackHandler
	irc_conn.Debug = config.Debug
	irc_conn.Password = config.Password

	irc_conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		fields := strings.Fields(strings.ToLower(e.Message()))
		if strings.HasPrefix(e.Message(), Conf.Prefix) {
			command := fields[0][1:]
			if fn, ok := Handlers[command]; ok {
				msg, err := fn(fields, e)
				if err != nil {
					log.Println(err.Error())
				}
				_pause()
				irc_conn.Privmsg(e.Arguments[0], msg)
			}
		}
	})
	irc_conn.AddCallback("JOIN", func(e *irc.Event) {
		//log.Println(e.Nick)
	})
	irc_conn.AddCallback("PART", func(e *irc.Event) {
		//log.Println(e.Nick)
	})
	irc_conn.AddCallback("001", func(e *irc.Event) {
		irc_conn.SendRaw("CAP REQ :twitch.tv/membership")
		//		irc_conn.SendRaw("CAP REQ :twitch.tv/tags")
		irc_conn.SendRaw("CAP REQ :twitch.tv/commands")
		_pause()
		irc_conn.Join(fmt.Sprintf("#%s", Conf.Name))
		for _, channel := range config.AutoJoin {
			irc_conn.Join(channel)
		}
	})

	err := irc_conn.Connect(config.Server)
	if err != nil {
		return nil, err
	}

	ticker := time.NewTicker(300 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				if UpdateGameData {
					player_info, err := GetPlayerInfo(Conf.ApiKey, SteamID(Conf.SteamID))
					if err != nil {
						log.Errorln(err.Error())
						continue
					}
					if player_info.GameServerIP != LastGameIP {
						irc_conn.Privmsgf(
							fmt.Sprintf("#%s", Conf.Name),
							"[Game Update] %s - steam://connect/%s", player_info.GameExtraInfo, player_info.GameServerIP,
						)
						LastGameIP = player_info.GameServerIP
					}
				}
			case <-stop_chan:
				ticker.Stop()
				return
			}
		}
	}()

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
	close(stop_chan)
	db.Close()
}

func init() {
	// Maps !{command} to actual functions
	Handlers = map[string]IRCHandlerFunc{
		"viewers":    HandleViewers,
		"steamid":    HandleGetSteamID,
		"setsteamid": HandleSetSteamID,
		"mysteamid":  HandleMySteamID,
		"profile":    HandleProfile,
		"myprofile":  HandleMyProfile,
		"commands":   HandleCommands,
		"mvm":        HandleMVM,
		"mymvm":      HandleMVM,
		"ip":         HandleIP,
		"startip":    HandleStartIP,
		"stopip":     HandleStopIP,
		"quit":       HandleQuit,
		"mvmlobby":   HandleMVMLobby,
		"scm":        HandleSCM,
	}

	stop_chan = make(chan struct{})
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
