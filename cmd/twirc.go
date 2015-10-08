package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/leighmacdonald/twirc"
)


func main() {
	twirc.LoadConfig()
	irc_conn, err := twirc.NewIRCClient(twirc.Conf)
	if err != nil {
		irc_conn.Quit()
		log.Fatalln(err.Error())
	}
	irc_conn.Loop()
}