package config

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"crdx.org/dex/cmd/dexd/env"
	"crdx.org/dex/pkg/util"
	"github.com/imroc/req/v3"
	"github.com/samber/lo"
)

type Config struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key"`
}

var config *Config

func Edit() {
	configFilePath := getConfigFilePath()
	if !util.PathExists(configFilePath) {
		lo.Must0(Save("", ""))
	}

	editor := cmp.Or(os.Getenv("EDITOR"), "vim")
	lo.Must0(util.PassThrough(editor, configFilePath))
}

func Save(key string, url string) error {
	c := Config{
		APIKey:  key,
		BaseURL: url,
	}

	b, err := json.Marshal(&c)
	if err != nil {
		return err
	}

	configFilePath := getConfigFilePath()

	configDir := filepath.Dir(configFilePath)
	if err := os.MkdirAll(configDir, 0o775); err != nil {
		return err
	}

	return os.WriteFile(configFilePath, b, 0o644)
}

func Load() error {
	configFilePath := getConfigFilePath()

	b, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	if len(b) == 0 {
		return errors.New("no configuration")
	}

	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		return err
	}

	config = &c
	return nil
}

func Print() {
	fmt.Printf("BASE_URL=%s\n", config.BaseURL)
	fmt.Printf("API_KEY=%s\n", config.APIKey)
}

func Endpoint(s string) string {
	return lo.Must(url.JoinPath(config.BaseURL, s))
}

func Request() *req.Request {
	c := req.C()
	if env.Debug() {
		c.DevMode()
	}
	return c.R().SetBearerAuthToken(config.APIKey)
}

// —————————————————————————————————————————————————————————————————————————————————————————————————

func getConfigFilePath() string {
	if os.Getenv("DEX_CONFIG_PATH") != "" {
		return os.Getenv("DEX_CONFIG_PATH")
	} else {
		configHome := cmp.Or(os.Getenv("XDG_CONFIG_HOME"), filepath.Join(os.Getenv("HOME"), ".config"))
		return filepath.Join(configHome, "dex", "config.json")
	}
}
