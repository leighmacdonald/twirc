# twirc
A very simple stand-alone Twitch bot that does some basic interactions with the Steam API.


## Commands 

**!mvm** returns the current total and individual tour counts as well as total completed missions for 
any individual tour.

    < !mvm
    > [MvM Info] Tours: 470 | Mecha(4): 1/3 | TwoCities(446): 2/4 | GearGrinder(20): 0/3
    < !mvm manofsnow
    > [MvM Info] Tours: 102 | TwoCities(102): 0/4

**!steamid** returns the numeric steam id for a given vanity name.

    < !steamid b4nny
    > Steam ID: b4nny => 76561197970669109


**!setsteamid** Allows a user to associate their steamid with their twitch username

	< !setsteamid 76561197970669109
    > Set steam id successfully

**!viewers** Fetch and return the current active viewer counts for the channel.

    < !viewers
    > [Viewers] Currently 0 viewers online.

**!ip** Returns the current game info including server ip. Only works for some games using steam api.

	< !ip
    > [Game] Team Fortress 2 - 192.69.96.156:27021

**!scm** Returns the most recent steam community market data for the requested item.

	< !scm Strange Rocket Launcher
    > [Market] Strange Rocket Launcher Lowest: $0.55 Volume: 133
    
    
**!startip** [Owner only] Start the background monitoring service that will announce whenever you connect or 
change game servers.

    < !startip
    > [Game] Started monitoring game state
    
**!stopip** [Owner only] Stop the background game info updater.
    
    < !stopip
    > [Game] Stopped monitoring game state


**!quit** [Owner only] Quit the server and shutdown the bot cleanly.

    < !quit
    > [Death Scene] Twas a scratch!
    

## Configuration

Create a config.toml file from the config_dist.toml example file. Edit the fields as 
you see fit. See below for a description of all the available fields.

	// The Twitch IRC server, shouldn't need to change this
	Server   string
	
    // Twitch username
	Name     string

	// Twitch oauth IRC key
	// See: http://www.twitchapps.com/tmi/
	Password string

	// Steam API Key
	// See: http://steamcommunity.com/dev/apikey
	ApiKey   string

    // Streamers steam ID, used as defaults for some commands
	SteamID string

	// Join these extra channels automatically, your own channel is joined by default already
	AutoJoin []string

	// Send error and debugging messages to this channel
	// You probably want to use a different channel from your main
	DebugChannel string

	// Print callback handler debug info to the console. Prints event handler debug
	// info to the console. You should never need this unless developing yourself.
	VerboseCallbackHandler bool

	// Print debug irc info to the console. Shows all raw irc traffic.
	Debug bool

	// Path to database file to use as permanent storage
	Database = "twirc.db"