/*
 * Copyright 2019 Joaquim Rocha <jrocha@gmailbox.org> and Contributors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package config provides functions to save and read configurations
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Config defines the structure a configuration file
type Config struct {
	DisableCompression      bool `json:"disable-compression"`
	DisableKeepAlives       bool `json:"disable-keep-alives"`
	MaxConnections          int  `json:"max-connections"`
	MaxIdleConnections      int  `json:"max-idle-connections"`
	RequestTimeout          int  `json:"request-timeout"`
	DisableCertificateCheck bool `json:"disable-certificate-check"`
	DisableRedirects        bool `json:"disable-redirects"`
}

// GetMaxIdleConnections if Config.MaxIdleConnections is zero resturns parallelConns else will return Config.MaxIdleConnections
func (c Config) GetMaxIdleConnections(parallelConns int) int {
	if c.MaxIdleConnections == 0 {
		return parallelConns
	}
	return c.MaxIdleConnections
}

// Default return the default configuration
func Default() *Config {
	return &Config{
		DisableCompression:      true,
		DisableKeepAlives:       false,
		MaxConnections:          0,
		MaxIdleConnections:      0,
		RequestTimeout:          30,
		DisableCertificateCheck: false,
		DisableRedirects:        true,
	}
}

// Read reads the configuration from a file
func Read(fileName string) (*Config, error) {
	data, err := os.ReadFile(fileName) //nolint:gosec // reads the config file the user named
	if err != nil {
		return nil, err
	}

	cfg := Default()
	if err := json.Unmarshal(data, cfg); err != nil {
		log.Printf("Invalid JSON config file %s: %v\n", fileName, err)
	}

	if err := checkConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func checkConfig(cfg *Config) error {
	if cfg.MaxConnections < 0 {
		return fmt.Errorf("%w: 'max-connections' must be zero or positive", ErrInvalidConfig)
	}

	if cfg.MaxIdleConnections < 0 {
		return fmt.Errorf("%w: 'max-idle-connections' must be zero or positive", ErrInvalidConfig)
	}

	if cfg.RequestTimeout < 0 {
		return fmt.Errorf("%w: 'request-timeout' must be zero or positive", ErrInvalidConfig)
	}

	return nil
}

// Write writes a configuration to a file
func Write(fileName string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, data, 0o600)
}
