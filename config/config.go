package config

import (
	"time"
)

type serverConfig struct {
	Port            []int
	DefaultLanguage string
	Hostname        time.Time
}

type server stuct {
	Config serverConfig
}

type servers map[string]server

func ReadConfig() { 
	var config servers
	if _, err := Decode(tomlBlob, &config); err != nil {
		log.Fatal(err)
	}

}
