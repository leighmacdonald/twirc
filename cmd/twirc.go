package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/leighmacdonald/twirc"
)

func main() {
	defer twirc.Shutdown()
	irc_conn, err := twirc.New(twirc.Conf)
	if err != nil {
		irc_conn.Quit()
		log.Fatalln(err.Error())
	}
	irc_conn.Loop()
}
