package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"regexp"
)

type Bot struct {
	server        string
	port          string
	nick          string
	user          string
	channel       string
	pass          string
	pread, pwrite chan string
	conn          net.Conn
}

func NewBot() *Bot {
	return &Bot{
		server:  "irc.freenode.net",
		port:    "6667",
		nick:    "steward",
		channel: "#hackedu",
		pass:    "",
		conn:    nil,
		user:    "steward"}
}

func (bot *Bot) Connect() (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", bot.server+":"+bot.port)
	if err != nil {
		log.Fatal("Unable to connect to IRC server ", err)
	}
	bot.conn = conn
	log.Printf("Connected to IRC server %s (%s)\n", bot.server,
		bot.conn.RemoteAddr())
	return bot.conn, nil
}

func (bot *Bot) Cmd(strfmt string, args ...interface{}) {
	fmt.Fprintf(bot.conn, strfmt+"\r\n", args...)
}

func main() {
	var msg = regexp.MustCompile(`PRIVMSG`)

	bot := NewBot()
	conn, _ := bot.Connect()
	bot.Cmd("USER %s 8 * :%s", bot.nick, bot.nick)
	bot.Cmd("NICK %s", bot.nick)
	bot.Cmd("JOIN %s", bot.channel)
	defer conn.Close()

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)
	for {
		line, err := tp.ReadLine()
		if err != nil {
			break
		}

		if msg.MatchString(line) {
			bot.Cmd("PRIVMSG %s :%s", bot.channel, line)
		}
		fmt.Printf("%s\n", line)
	}
}
