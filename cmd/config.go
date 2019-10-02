package cmd

import (
	"fmt"

	"github.com/jjmrocha/beast/config"
)

func Config(fileName string) {
	defaults := &config.Config{
		DisableCompression: true,
		DisableKeepAlives:  false,
		MaxConnections:     0,
		Timeout:            30,
	}
	config.Write(fileName, defaults)
	fmt.Printf("File %s was created with default configuration\n", fileName)
}
