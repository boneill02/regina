package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"crypto/tls"
	"github.com/mattn/go-xmpp"
)

var server = flag.String("server", "localhost:5223", "server")
var username = flag.String("username", "me", "username")
var password = flag.String("password", "password", "password")
var status = flag.String("status", "xa", "status")
var statusMessage = flag.String("status-msg", "i am on xmpp.", "status message")
var notls = flag.Bool("notls", false, "No TLS")
var debug = flag.Bool("debug", false, "debug output")
var session = flag.Bool("session", false, "use server session")

func serverName(host string) string {
	return strings.Split(host, ":")[0]
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: sxmpp [options]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()

	xmpp.DefaultConfig = tls.Config{
			ServerName:         serverName(*server),
			InsecureSkipVerify: false,
	}

	options := xmpp.Options{
		Host:          *server,
		User:          *username,
		Password:      *password,
		NoTLS:         *notls,
		Debug:         *debug,
		Session:       *session,
		Status:        *status,
		StatusMessage: *statusMessage,
	}

	client, err := options.NewClient()
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			chat, err := client.Recv()
			if err != nil {
				panic(err)
			}

			switch v := chat.(type) {
			case xmpp.Chat:
				fmt.Println(v.Remote, v.Text)
			case xmpp.Presence:
				fmt.Println(v.From, v.Show)
			}
		}
	}()

	for {
		in := bufio.NewReader(os.Stdin)

		line, err := in.ReadString('\n')
		if err != nil {
			continue
		}
		line = strings.TrimRight(line, "\n")
		tokens := strings.SplitN(line, " ", 2)
		if len(tokens) == 2 {
			client.Send(xmpp.Chat{Remote: tokens[0], Type: "chat", Text: tokens[1]})
		}
	}
}
