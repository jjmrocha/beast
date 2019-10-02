package cmd

import (
	"encoding/json"
	"os"

	"github.com/jjmrocha/beast/models"
)

func WriteDefaultConfig() error {
	f, err := os.Create("config.json")
	if err != nil {
		return err
	}
	return json.NewEncoder(f).Encode(models.NewDefaultConfig())
}
