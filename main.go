package main

import (
	irc "./irc"
	web "./web"
	crypt "crypto/tls"
	"fmt"
	client "github.com/fluffle/goirc/client"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Cannot read config file!")
		os.Exit(3)
	}
	conf := irc.Read(os.Args[1])

	CheckAndCreate(conf.DataPath)
	CheckAndCreate(conf.DataPath + "/thumb")
	CheckAndCreate(conf.DataPath + "/image")
	CheckAndCreate(conf.DataPath + "/gif")
	CheckAndCreate(conf.DataPath + "/webm")

	cfg := client.NewConfig(conf.BotName)
	cfg.SSL = conf.Ssl
	cfg.Server = fmt.Sprintf("%s:%d", conf.Server, conf.Port)
	cfg.Pass = conf.ServerPassword

	c := client.Client(cfg)
	handler := irc.NewHandler(c, conf.Channel)

	var ssl crypt.Config
	ssl.InsecureSkipVerify = conf.AcceptInvalidCert
	cfg.SSLConfig = &ssl
	quit := make(chan bool)

	c.HandleFunc("connected",
		func(conn *client.Conn, line *client.Line) {
			conn.Join(fmt.Sprintf("%s %s", conf.Channel, conf.ChannelPassword))
		})

	c.HandleFunc("privmsg",
		func(conn *client.Conn, line *client.Line) {
			handler.Recv(line.Args[1], line.Nick)
		})

	c.HandleFunc("disconnected",
		func(conn *client.Conn, line *client.Line) {
			quit <- true
		})

	err := c.Connect()
	if err != nil {
		fmt.Printf("Connection error: %s\n", err)
	}
	go irc.ImageAction.StartupCheck()
	web.StartWeb(conf)
	<-quit
}

func CheckAndCreate(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}
