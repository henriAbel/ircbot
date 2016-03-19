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
	conf          *irc.Config
	quit          chan bool
	lastReconnect time.Time
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Cannot read config file!")
		os.Exit(3)
	}
	conf = irc.Read(os.Args[1])
	if len(irc.GetConfig().LogFile) > 0 {
		f, err := os.OpenFile(irc.GetConfig().LogFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
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
	CheckAndCreate(conf.DataPath)
	CheckAndCreate(filepath.Join(conf.DataPath, "thumb"))
	CheckAndCreate(filepath.Join(conf.DataPath, "image"))
	CheckAndCreate(filepath.Join(conf.DataPath, "gif"))
	CheckAndCreate(filepath.Join(conf.DataPath, "webm"))

	cfg := client.NewConfig(conf.BotName)
	cfg.SSL = conf.Ssl
	cfg.Server = fmt.Sprintf("%s:%d", conf.Server, conf.Port)
	cfg.Pass = conf.ServerPassword

	c = client.Client(cfg)
	handler := irc.NewHandler(c, conf.Channel)

	var ssl crypt.Config
	ssl.InsecureSkipVerify = conf.AcceptInvalidCert
	cfg.SSLConfig = &ssl

	quit = make(chan bool)
	lastReconnect = time.Now()

	c.HandleFunc("connected",
		func(conn *client.Conn, line *client.Line) {
			log.Infof("Connected to server %s, joining channel %s", conf.Server, conf.Channel)
			conn.Join(fmt.Sprintf("%s %s", conf.Channel, conf.ChannelPassword))
		})

	c.HandleFunc("privmsg",
		func(conn *client.Conn, line *client.Line) {
			handler.Recv(line.Args[1], line.Nick)
		})

	c.HandleFunc("disconnected",
		func(conn *client.Conn, line *client.Line) {
			log.Infof("Disconnected from server %s", conf.Server)
			handleDisconnect()
		})
	c.HandleFunc("KICK",
		func(conn *client.Conn, line *client.Line) {
			log.Infof("Kicked from channel %s", conf.Channel)
			if conf.AutoReJoin {
				log.Infof("Joining channel %s", conf.Channel)
				conn.Join(fmt.Sprintf("%s %s", conf.Channel, conf.ChannelPassword))
			}
		})
	err := connect(conf)
	if err != nil {
		log.Errorf(err.Error())
		handleDisconnect()
	}
	go irc.ImageAction.StartupCheck()
	web.StartWeb(conf)
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
	if conf.AutoReconnect {
		go func() {
			if time.Since(lastReconnect).Seconds() < 10 {
				time.Sleep(time.Duration(10-time.Since(lastReconnect).Seconds()) * time.Second)
			}
			err := connect(conf)
			if nil != err {
				log.Errorf(err.Error())
				handleDisconnect()
			}
		}()
	} else {
		quit <- true
	}
}

func connect(conf *irc.Config) error {
	log.Infof("Connecting to server %s", conf.Server)
	lastReconnect = time.Now()
	return c.Connect()
}
