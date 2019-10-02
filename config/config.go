package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	DisableCompression bool `json:"disable-compression"`
	DisableKeepAlives  bool `json:"disable-keep-alives"`
	MaxConnections     int  `json:"max-connections"`
	Timeout            uint `json:"timeout"`
}

func Write(fileName string, config *Config) {
	data, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		log.Printf("Error encoding config %v to JSON: %v\n", config, err)
	}

	err = ioutil.WriteFile(fileName, data, 0666)
	if err != nil {
		log.Printf("Error writing to file %s: %v\n", fileName, err)
	}
}
