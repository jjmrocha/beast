package cmd

import (
	"fmt"

	"github.com/jjmrocha/beast/config"
)

func Config(fileName string) {
	defaults := config.Default()
	config.Write(fileName, defaults)
	fmt.Printf("File %s was created with default configuration\n", fileName)
}
