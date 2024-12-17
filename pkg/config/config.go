package config

import (
	jsonStd "encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
)

type SignRule struct {
	Contract           string `json:"contract"`
	PasswordPrompt     bool   `json:"passwordPrompt"`
	ConfirmationBypass bool   `json:"confirmationBypass"`
}

type Account struct {
	SignRules []SignRule `json:"signRules"`
}

type Config struct {
	Accounts map[string]Account `json:"accounts"`
}

var (
	cfg  *Config
	once sync.Once
	k    = koanf.New(".")
)

const (
	fileName = "wallet_config.json"
)

func Load() *Config {
	once.Do(func() {
		configDir, err := wallet.Path()
		if err != nil {
			log.Fatalf("Failed to get wallet config directory: %v", err)
		}

		path := filepath.Join(configDir, fileName)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Printf("config file not found. Creating default config at %s", path)
			createDefaultConfigFile(path, defaultConfig())
		}

		if err := k.Load(file.Provider(path), json.Parser()); err != nil {
			log.Fatalf("Error loading config: %v", err)
		}

		cfg = &Config{}
		if err := k.Unmarshal("", cfg); err != nil {
			log.Fatalf("Error unmarshaling config: %v", err)
		}
	})

	return cfg
}

func Get() *Config {
	if cfg == nil {
		log.Fatal("Config not loaded. Call Load() first.")
	}

	return cfg
}

func defaultConfig() *Config {
	return &Config{
		Accounts: map[string]Account{},
	}
}

func createDefaultConfigFile(filePath string, defaultConfig *Config) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("failed to create config file: %v", err)
	}
	defer file.Close()

	encoder := jsonStd.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print JSON

	if err := encoder.Encode(defaultConfig); err != nil {
		log.Fatalf("failed to write default config to file: %v", err)
	}
}
