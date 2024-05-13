package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ContentEnc    string `json:"content_enc"`
	TotpMasterEnc string `json:"totp_master_enc"`
}

// NewByPath
func NewByPath(path string) (*Config, error) {
	if inf, err := os.Stat(path); err != nil {
		return nil, err
	} else if inf.IsDir() {
		return nil, fmt.Errorf("%s is a directory", path)
	}

	jsonContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return New(jsonContent)
}

// New
func New(jsonContent []byte) (*Config, error) {
	var config Config
	if err := json.Unmarshal(jsonContent, &config); err != nil {
		return nil, err
	}

	if config.ContentEnc == "" {
		return nil, fmt.Errorf(`you must set "content_enc" field`)
	}

	if config.TotpMasterEnc == "" {
		return nil, fmt.Errorf(`you must set "totp_master_enc" field`)
	}

	return &config, nil
}
