package irc

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/spf13/viper"
)

func InitConf(configFile string) {
	cfgData, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println(err)
	}
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(cfgData))
}

func GetServerAddress() string {
	return viper.GetString("irc.server.address")
}

func GetServerPort() int {
	return viper.GetInt("irc.server.port")
}

func GetServerPassword() string {
	return viper.GetString("irc.server.serverPassword")
}

func GetChannel() string {
	return viper.GetString("irc.server.channel")
}

func GetChannelPassword() string {
	return viper.GetString("irc.server.channelPassword")
}

func GetBotName() string {
	return viper.GetString("irc.server.name")
}

func GetUseSsl() bool {
	return viper.GetBool("irc.server.ssl")
}

func GetAcceptInvalidCert() bool {
	return viper.GetBool("irc.server.acceptInvalidCert")
}

func GetAutoReJoin() bool {
	return viper.GetBool("irc.server.autoReJoin")
}

func GetAutoReconnect() bool {
	return viper.GetBool("irc.server.autoReoconnect")
}

func GetDbFile() string {
	return viper.GetString("irc.system.dbFile")
}

func GetDataPath() string {
	return viper.GetString("irc.system.dataPath")
}

func GetYoutubeApiKey() string {
	return viper.GetString("irc.system.youtubeApiKey")
}

func GetWebPort() int {
	return viper.GetInt("web.port")
}

func GetWebPassword() string {
	return viper.GetString("web.password")
}

func GetLogFile() string {
	return viper.GetString("irc.system.logFile")
}

func GetCertFile() string {
	return viper.GetString("web.tls.certificate")
}

func GetKeyFile() string {
	return viper.GetString("web.tls.keyfile")
}
