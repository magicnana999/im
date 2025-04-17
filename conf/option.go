package conf

import (
	"fmt"
	"github.com/magicnana999/im/pkg/ip"
	"github.com/magicnana999/im/pkg/str"
	"os"
	"strings"
)

var (
	Global Option
)

type Option struct {
	Name    string  `yaml:"name"`
	Logger  Logger  `yaml:"logger"`
	Mysql   MySQL   `yaml:"mysql"`
	Redis   Redis   `yaml:"cache"`
	Kafka   Kafka   `yaml:"kafka"`
	Broker  Broker  `yaml:"broker"`
	Service Service `yaml:"cmd_service"`
}

type MySQL struct {
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Schema    string `yaml:"schema"`
	Charset   string `yaml:"charset"`
	ParseTime string `yaml:"parseTime"`
	Location  string `yaml:"location"`
}

func (m MySQL) String() string {
	format := "%s:%s@tcp(%s:%s)/%s"
	first := fmt.Sprintf(format, m.User, m.Password, m.Host, m.Port, m.Schema)

	s := make([]string, 0)
	if str.IsNotBlank(m.Charset) {
		s = append(s, fmt.Sprintf("?charset=%s", m.Charset))
	}
	if str.IsNotBlank(m.ParseTime) {
		s = append(s, fmt.Sprintf("parseTime=%s", m.ParseTime))
	}
	if str.IsNotBlank(m.Location) {
		s = append(s, fmt.Sprintf("loc=%s", m.Location))
	}
	second := strings.Join(s, "&")
	return first + second
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

func (r Redis) String() string {
	return r.Host + ":" + r.Port
}

type Kafka struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func (k Kafka) String() string {
	return fmt.Sprintf("%s:%s", k.Host, k.Port)
}

type Broker struct {
	Addr              string `yaml:"addr"`
	ServerInterval    int    `yaml:"serverInterval"`
	HeartbeatInterval int    `yaml:"heartbeatInterval"`
	LoggerLevel       string `yaml:"loggerLevel"`
}

type Service struct {
	Addr string `yaml:"addr"`
}

type Logger struct {
	Level int8 `yaml:"level"`
}

func LoadConfig(path string) error {

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var config Option
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	Global = config

	if config.Broker.Addr == "" {
		i, e := ip.GetLocalIP()
		if e != nil {
			fmt.Errorf("can not get local ip %v", e)
		}

		config.Broker.Addr = fmt.Sprintf("%s:7539", i)
	}

	return nil
}
