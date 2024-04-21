package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v3"
)

// Config holds the configuration options for the application.
//
// At the moment, it is quite limited, only supporting the home folder and the
// file name of the metadata.
type Config struct {
	Home string `env:"ZZZ_HOME" yaml:"home"`
	File string `env:"ZZZ_FILE" yaml:"file"`

	DefaultLanguage string `env:"ZZZ_DEFAULT_LANGUAGE" yaml:"default_language"`

	Theme string `env:"ZZZ_THEME" yaml:"theme"`

	PrimaryColor        string `env:"ZZZ_PRIMARY_COLOR" yaml:"primary_color"`
	PrimaryColorSubdued string `env:"ZZZ_PRIMARY_COLOR_SUBDUED" yaml:"primary_color_subdued"`
	BrightGreenColor    string `env:"ZZZ_BRIGHT_GREEN" yaml:"bright_green"`
	GreenColor          string `env:"ZZZ_GREEN" yaml:"green"`
	BrightRedColor      string `env:"ZZZ_BRIGHT_RED" yaml:"bright_red"`
	RedColor            string `env:"ZZZ_RED" yaml:"red"`
	ForegroundColor     string `env:"ZZZ_FOREGROUND" yaml:"foreground"`
	BackgroundColor     string `env:"ZZZ_BACKGROUND" yaml:"background"`
	GrayColor           string `env:"ZZZ_GRAY" yaml:"gray"`
	BlackColor          string `env:"ZZZ_BLACK" yaml:"black"`
	WhiteColor          string `env:"ZZZ_WHITE" yaml:"white"`
}

func newConfig() Config {
	return Config{
		Home:                defaultHome(),
		File:                "snippets.json",
		DefaultLanguage:     defaultLanguage,
		Theme:               "dracula",
		PrimaryColor:        "#AFBEE1",
		PrimaryColorSubdued: "#64708D",
		BrightGreenColor:    "#BCE1AF",
		GreenColor:          "#527251",
		BrightRedColor:      "#E49393",
		RedColor:            "#A46060",
		ForegroundColor:     "15",
		BackgroundColor:     "235",
		GrayColor:           "241",
		BlackColor:          "#373b41",
		WhiteColor:          "#FFFFFF",
	}
}

// default helpers for the configuration.
// We use $XDG_DATA_HOME to avoid cluttering the user's home directory.
func defaultHome() string { return filepath.Join(xdg.DataHome, "zzz") }

// defaultConfig returns the default config path
func defaultConfig() string {
	if c := os.Getenv("ZZZ_CONFIG"); c != "" {
		return c
	}
	cfgPath, err := xdg.ConfigFile("zzz/config.yaml")
	if err != nil {
		return "config.yaml"
	}
	return cfgPath
}

// readConfig returns a configuration read from the environment.
func readConfig() Config {
	config := newConfig()
	fi, err := os.Open(defaultConfig())
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return newConfig()
	}
	if fi != nil {
		defer fi.Close()
		if err := yaml.NewDecoder(fi).Decode(&config); err != nil {
			return newConfig()
		}
	}

	if err := env.Parse(&config); err != nil {
		return newConfig()
	}

	if strings.HasPrefix(config.Home, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			config.Home = filepath.Join(home, config.Home[1:])
		}
	}

	return config
}

// writeConfig returns a configuration read from the environment.
func (config Config) writeConfig() error {
	fi, err := os.Create(defaultConfig())
	if err != nil {
		return err
	}
	if fi != nil {
		defer fi.Close()
		if err := yaml.NewEncoder(fi).Encode(&config); err != nil {
			return err
		}
	}

	return nil
}
