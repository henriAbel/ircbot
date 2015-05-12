package irc

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	Server,
	ServerPassword,
	Channel,
	ChannelPassword,
	BotName,
	LogFile,
	DBFile,
	YoutubeApiKey string
	Port                   int
	Ssl, AcceptInvalidCert bool
}

var (
	cfg      Config
	location string
)

func GetConfig() *Config {
	return Read(location)
}

func Read(configFile string) *Config {
	if cfg == (Config{}) {
		cfg = Config{}
		location = configFile
		ps := reflect.Indirect(reflect.ValueOf(&cfg))
		reader, _ := os.Open(configFile)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
				continue
			}
			split := strings.Split(line, ":")
			for i := 0; i < ps.NumField(); i++ {
				if strings.ToLower(ps.Type().Field(i).Name) == strings.ToLower(strings.TrimSpace(split[0])) {
					field := ps.FieldByName(ps.Type().Field(i).Name)
					value := strings.TrimSpace(strings.Join(split[1:], ":"))
					if field.Kind() == reflect.String {
						field.SetString(value)
					} else if field.Kind() == reflect.Int {
						intValue, err := strconv.ParseInt(value, 10, 64)
						if err != nil {
							fmt.Sprintln("Cannot parse '%s' integer value '%s'", ps.Type().Field(i).Name, value)
						}
						field.SetInt(intValue)
					} else if field.Kind() == reflect.Bool {
						boolValue, err := strconv.ParseBool(value)
						if err != nil {
							fmt.Sprintln("Cannot parse '%s' boolean value '%s'", ps.Type().Field(i).Name, value)
						}
						field.SetBool(boolValue)
					} else {
						fmt.Println("Unknown field type")
					}
					break
				}
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Cannot open file: ", err)
		}
	}
	return &cfg
}
