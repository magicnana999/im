package global

import (
	"fmt"
	"github.com/magicnana999/im/pkg/ip"
	"github.com/magicnana999/im/pkg/logger"
	"github.com/magicnana999/im/pkg/str"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"sync"
)

var (
	C     Config
	cOnce sync.Once
)

type MicroService struct {
	EtcdAddr   []string `yaml:"etcdAddr"`
	BrokerName string   `yaml:"brokerName"`
	BrokerAddr string   `yaml:"brokerAddr"`
	RouterName string   `yaml:"routerName"`
	RouterAddr string   `yaml:"routerAddr"`
	ServerName string   `yaml:"serverName"`
	ServerAddr string   `yaml:"serverAddr"`
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
	Name              string `yaml:"name"`
	Addr              string `yaml:"addr"`
	TickerInterval    string `yaml:"tickerInterval"`
	HeartbeatInterval int    `yaml:"heartbeatInterval"`
}

type Config struct {
	Broker       Broker       `yaml:"broker"`
	MicroService MicroService `yaml:"microService"`
	Mysql        MySQL        `yaml:"mysql"`
	Redis        Redis        `yaml:"redis"`
	Kafka        Kafka        `yaml:"kafka"`
}

func Load(path string) (Config, error) {

	cOnce.Do(func() {
		data, err := os.ReadFile(path)
		if err != nil {
			logger.Fatalf("load config failed: %v", err)
		}

		err = yaml.Unmarshal(data, &C)
		if err != nil {
			logger.Fatalf("load config failed: %v", err)
		}

		if C.Broker.Addr == "" {
			i, e := ip.GetLocalIP()
			if e != nil {
				fmt.Errorf("can not get local ip %v", e)
			}

			C.Broker.Addr = fmt.Sprintf("%s:7539", i)
		}
	})

	return C, nil
}

func GetMicroService() MicroService {
	return C.MicroService
}

func GetMysql() MySQL {
	return C.Mysql
}

func GetRedis() Redis {
	return C.Redis
}

func GetKafka() Kafka {
	return C.Kafka
}

func GetBroker() Broker {
	return C.Broker
}
