package main

import (
	crypt "crypto/tls"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	irc "./irc"
	web "./web"
	log "github.com/Sirupsen/logrus"
	client "github.com/fluffle/goirc/client"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	c             *client.Conn
	quit          chan bool
	lastReconnect time.Time
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Cannot read config file!")
		os.Exit(3)
	}
	irc.InitConf(os.Args[1])
	if len(irc.GetLogFile()) > 0 {
		f, err := os.OpenFile(irc.GetLogFile(), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			fmt.Printf("Error opening log file: %v", err)
			os.Exit(3)
		}
		defer f.Close()
		log.SetOutput(f)
	} else {
		if runtime.GOOS == "windows" {
			var format = new(prefixed.TextFormatter)
			format.ForceColors = true
			log.SetFormatter(format)
		}
	}

	log.SetLevel(log.DebugLevel)
	log.Infof("<<<<<----- Application started %s ----->>>>>", time.Now())
	irc.InitImageAction()
	CheckAndCreate(irc.GetDataPath())
	CheckAndCreate(filepath.Join(irc.GetDataPath(), "thumb"))
	CheckAndCreate(filepath.Join(irc.GetDataPath(), "image"))
	CheckAndCreate(filepath.Join(irc.GetDataPath(), "gif"))
	CheckAndCreate(filepath.Join(irc.GetDataPath(), "webm"))

	cfg := client.NewConfig(irc.GetBotName())
	cfg.SSL = irc.GetUseSsl()
	cfg.Server = fmt.Sprintf("%s:%d", irc.GetServerAddress(), irc.GetServerPort())
	cfg.Pass = irc.GetServerPassword()

	c = client.Client(cfg)
	handler := irc.NewHandler(c, irc.GetChannel())

	var ssl crypt.Config
	ssl.InsecureSkipVerify = irc.GetAcceptInvalidCert()
	cfg.SSLConfig = &ssl

	quit = make(chan bool)
	lastReconnect = time.Now()

	c.HandleFunc("connected",
		func(conn *client.Conn, line *client.Line) {
			log.Infof("Connected to server %s, joining channel %s", irc.GetServerAddress(), irc.GetChannel())
			conn.Join(fmt.Sprintf("%s %s", irc.GetChannel(), irc.GetChannelPassword()))
		})

	c.HandleFunc("privmsg",
		func(conn *client.Conn, line *client.Line) {
			handler.Recv(line.Args[1], line.Nick)
		})

	c.HandleFunc("disconnected",
		func(conn *client.Conn, line *client.Line) {
			log.Infof("Disconnected from server %s", irc.GetServerAddress())
			handleDisconnect()
		})
	c.HandleFunc("KICK",
		func(conn *client.Conn, line *client.Line) {
			log.Infof("Kicked from channel %s", irc.GetChannel())
			if irc.GetAutoReJoin() {
				log.Infof("Joining channel %s", irc.GetChannel())
				conn.Join(fmt.Sprintf("%s %s", irc.GetChannel(), irc.GetChannelPassword()))
			}
		})
	err := connect()
	if err != nil {
		log.Errorf(err.Error())
		handleDisconnect()
	}
	go irc.ImageAction.StartupCheck()
	web.StartWeb()
	<-quit
}

// CheckAndCreate checks if directory exisrs, if not then creats
func CheckAndCreate(path string) {
	log.Debugf("Checking directory: %s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Debugf("Directory %s don't exist, creating", path)
		os.Mkdir(path, os.ModePerm)
	}
}

func handleDisconnect() {
	if irc.GetAutoReconnect() {
		go func() {
			if time.Since(lastReconnect).Seconds() < 10 {
				time.Sleep(time.Duration(10-time.Since(lastReconnect).Seconds()) * time.Second)
			}
			err := connect()
			if nil != err {
				log.Errorf(err.Error())
				handleDisconnect()
			}
		}()
	} else {
		quit <- true
	}
}

func connect() error {
	log.Infof("Connecting to server %s", irc.GetServerAddress())
	lastReconnect = time.Now()
	return c.Connect()
}
