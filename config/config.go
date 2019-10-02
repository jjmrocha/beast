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
	Timeout            int  `json:"timeout"`
}

func Default() *Config {
	return &Config{
		DisableCompression: true,
		DisableKeepAlives:  false,
		MaxConnections:     0,
		Timeout:            30,
	}
}

func Read(fileName string) *Config {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error reading file %s: %v\n", fileName, err)
	}

	config := Default()
	json.Unmarshal(data, config)
	checkConfig(config)
	return config
}

func checkConfig(config *Config) {
	if config.MaxConnections < 0 {
		log.Fatalln("Invalid config, 'max-connections' must be zero or positive")
	}

	if config.Timeout < 0 {
		log.Fatalln("Invalid config, 'timeout' must be zero or positive")
	}
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
