package models

// Config contains the parameters for setting up a HTTP Connection
type Config struct {
	DisableCompression bool `json:"disable-compression"`
	DisableKeepAlives  bool `json:"disable-keep-alives"`
	MaxConnections     int  `json:"max-connections"`
	Timeout            uint `json:"timeout"`
}

// NewDefaultConfig creates a default config file
func NewDefaultConfig() Config {
	return Config{
		DisableCompression: true,
		Timeout:            30,
	}
}
