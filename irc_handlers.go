package twirc

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/thoj/go-ircevent"
	"log"
	"strings"
)

var (
	InvalidArguments = errors.New("Incomplete arguments supplied")
	InvalidResponse  = errors.New("Invalid response from remote")
	InternalError    = errors.New("Internal error processing request")
)

type IRCHandlerFunc func([]string, *irc.Event) (string, error)

func HandleViewers(fields []string, e *irc.Event) (string, error) {
	chatters, err := Chatters(fmt.Sprintf("#%s", e.Arguments[0]))
	if err != nil {
		return "Error finding channel viewers", err
	} else {
		return fmt.Sprintf("[Viewers] Currently %d viewers online.", chatters.ChatterCount), nil
	}
}

func HandleGetSteamID(fields []string, e *irc.Event) (string, error) {
	if len(fields) != 2 {
		return "[SteamID] Must supply a vanity name to resolve", InvalidArguments
	}

	steam_id, err := ResolveVanity(fields[1])
	if err != nil {
		return fmt.Sprintf("[SteamID] Error trying to resolve %s", fields[1]), InvalidResponse

	}
	return fmt.Sprintf("[SteamID] %s => %s", fields[1], steam_id), nil
}

func HandleSetSteamID(fields []string, e *irc.Event) (string, error) {
	if len(fields) != 2 {
		return "[SteamID] Command only takes 1 argument, the steamid or vanity name", InvalidArguments
	}
	steam_id, err := NewSteamID(fields[1])
	if err != nil {
		return "[SteamID] Error resolving steam id", InvalidResponse
	} else {
		if err := SetSteamID(e.Nick, steam_id); err != nil {
			return "[SteamID] Internal oopsie", InternalError
		} else {
			return "[SteamID] Set steam id successfully", nil
		}
	}
}

func HandleMySteamID(fields []string, e *irc.Event) (string, error) {
	steam_id, err := GetSteamID(strings.ToLower(e.Nick))
	if err != nil {
		return "[SteamID] Must set steam id with !setsteamid command first", InvalidArguments
	} else {
		return fmt.Sprintf("[SteamID] %s => %s", e.Nick, steam_id), nil
	}
}

func HandleProfile(fields []string, e *irc.Event) (string, error) {
	sid, err := NewSteamID(Conf.SteamID)
	if err != nil {
		return "Streamer SteamID not set or invalid", InvalidArguments
	}
	return fmt.Sprintf("[Steam Profile] %s", sid.ProfileURL()), nil
}

func HandleMyProfile(fields []string, e *irc.Event) (string, error) {
	steam_id, err := GetSteamID(e.Nick)
	if err != nil {
		return "Must first set steam id with !setsteamid <steamid>", InvalidArguments
	} else {
		return fmt.Sprintf("[Steam Profile] %s", steam_id.ProfileURL()), nil
	}

}

func HandleCommands(fields []string, e *irc.Event) (string, error) {
	return "[Commands] !scm, !setsteamid, !mvm, !mvmlobby, !mysteamid, !ip, !viewers", nil
}

func HandleMVM(fields []string, e *irc.Event) (string, error) {
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
			return "Failed to find valid SteamID", err
		}
	} else {
		sid, err := GetSteamID(e.Nick)
		if err != nil {
			return "Must first set steam id with !setsteamid <steamid>", InvalidArguments
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

	return buffer.String(), nil
}

func HandleSCM(fields []string, e *irc.Event) (string, error) {
	if len(fields) == 1 {
		return "[Market] Must supply item name to lookup", InvalidArguments
	} else {
		price, err := GetPrice(strings.Join(fields[1:], " "))
		if err != nil {
			return "[Market] Failed to fetch price", InvalidResponse
		} else {
			return fmt.Sprintf("[Market] %s Lowest: %s Volume: %d", price.Name, price.LowestPrice, price.Volume), nil
		}
	}
}

func HandleIP(fields []string, e *irc.Event) (string, error) {
	player_info, err := GetPlayerInfo(Conf.ApiKey, SteamID(Conf.SteamID))
	if err != nil {
		return "Could not fetch player data", InvalidResponse
	} else {
		if player_info.GameServerIP == "" {
			return fmt.Sprintf("[Game] %s Game info n/a or playing unsupported game.", Conf.Name), nil
		} else {
			return fmt.Sprintf("[Game] %s - steam://connect/%s", player_info.GameExtraInfo, player_info.GameServerIP), nil
		}
	}
}

func HandleMVMLobby(fields []string, e *irc.Event) (string, error) {
	//var steam_id SteamID
	sid := Conf.SteamID
	if len(fields) >= 2 {
		sid = fields[1]
	}
	steam_id, err := NewSteamID(sid)
	log.Println("Using ID:", steam_id)
	if err != nil {
		return "Failed to determine steam id, use !setsteamid to set it", err
	} else {
		return fmt.Sprintf("[MVMLobby] %s", steam_id.MVMLobbyURL()), nil
	}
}

func HandleStartIP(fields []string, e *irc.Event) (string, error) {
	if isOwner(e.Nick) {
		UpdateGameData = true
		return "[Game] Started monitoring game state", nil
	} else {
		return "[Game] Unauthorize, must be bot owner", nil
	}
}

func HandleStopIP(fields []string, e *irc.Event) (string, error) {
	if isOwner(e.Nick) {
		UpdateGameData = true
		return "[Game] Stopped monitoring game state", nil
	} else {
		return "[Game] Unauthorize, must be bot owner", nil
	}
}

func HandleQuit(fields []string, e *irc.Event) (string, error) {
	return "[Death Scene] Twas a scratch!", nil
}

func isOwner(nick string) bool {
	return strings.ToLower(nick) == strings.ToLower(Conf.Name)
}
