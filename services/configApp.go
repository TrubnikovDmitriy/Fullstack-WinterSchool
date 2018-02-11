package serv

import (
	"os"
	"encoding/json"
	"time"
)

type Config struct {
	Href                   string        `json:"href"`
	MaxMatchesInTourney    int           `json:"max_match_in_tourney"`
	MaxFieldLength         int           `json:"max_field_length"`
	MinPasswordLength      int           `json:"min_password_length"`
	MaxConnectionsToDB     int           `json:"max_connections_to_db"`
	NumberOfShards         int           `json:"number_of_shards"`
	SlaveToMasterReadRatio int32         `json:"slave_to_master_ratio"`
	AccessTokenTTL         time.Duration `json:"access_token_ttl"`
	RefreshTokenTTL        time.Duration `json:"refresh_token_ttl"`
	SignKey                string        `json:"sign_key"`
	Shards                 []Shards      `json:"shards"`
	Redis                  []Redis       `json:"redis"`
}

type Shards struct {
	Master DataBase `json:"master"`
	Slave  DataBase `json:"slave"`
}

type DataBase struct {
	Host   string `json:"host"`
	Port   uint16 `json:"port"`
	DBName string `json:"db_name"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
}

type Redis struct {
	Address string `json:"address"`
	Number  int    `json:"number"`
}


var ApplicationConfig *Config

func init() {
	ApplicationConfig = ReadConfig("application.cfg")
}

func ReadConfig(path string) *Config {

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := new(Config)

	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}

func GetConfig() *Config {
	return ApplicationConfig
}