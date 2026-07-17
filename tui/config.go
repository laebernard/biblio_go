package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var API_URL string
var AppConfig *Config
var AppState *State

type Config struct {
	Theme  string `json:"theme"`
	Accent string `json:"accent"`
	APIURL string `json:"api_url"`
}

type State struct {
	LastScreen string `json:"last_screen"`
	LastSearch string `json:"last_search"`
	Token      string `json:"token"`
}

func getConfigDir() string {
	return "config"
}

func ensureConfigDir() {
	os.MkdirAll(getConfigDir(), os.ModePerm)
}

func getConfigPath() string {
	return filepath.Join(getConfigDir(), "config.json")
}

func getStatePath() string {
	return filepath.Join(getConfigDir(), "state.json")
}

func LoadConfig() *Config {
	ensureConfigDir()

	path := getConfigPath()

	// fallback ENV
	envURL := os.Getenv("API_URL")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		cfg := &Config{
			Theme:  "dark",
			Accent: "#FF00FF",
			APIURL: envURL,
		}

		if cfg.APIURL == "" {
			cfg.APIURL = "http://localhost:8080"
		}

		SaveConfig(cfg)
		return cfg
	}

	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var cfg Config
	json.Unmarshal(data, &cfg)

	// si API vide → fallback
	if cfg.APIURL == "" {
		cfg.APIURL = envURL
		if cfg.APIURL == "" {
			cfg.APIURL = "http://localhost:8080"
		}
	}

	return &cfg
}

func SaveConfig(cfg *Config) {
	ensureConfigDir()

	data, _ := json.MarshalIndent(cfg, "", "  ")
	os.WriteFile(getConfigPath(), data, 0644)
}

func LoadState() *State {
	ensureConfigDir()

	path := getStatePath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		st := &State{}
		SaveState(st)
		return st
	}

	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var st State
	json.Unmarshal(data, &st)
	return &st
}

func SaveState(st *State) {
	ensureConfigDir()

	data, _ := json.MarshalIndent(st, "", "  ")
	os.WriteFile(getStatePath(), data, 0644)
}

func init() {
	AppConfig = LoadConfig()
	AppState = LoadState()

	API_URL = AppConfig.APIURL
}
