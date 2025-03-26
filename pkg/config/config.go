package config

import (
	jsonStd "encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
)

var (
	configManager *ConfigManager
	once          sync.Once
	k             = koanf.New(".")
)

const (
	fileName     = "wallet_config.json"
	nanoIDLength = 10 // Length of NanoID
)

var configFileDirOverride string

func Load() *ConfigManager {
	once.Do(func() {
		path := configFilePath()

		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Printf("Config file not found. Creating default config at %s", path)

			if err := saveConfigUnsafe(defaultConfig()); err != nil {
				log.Fatalf("Error saving config: %v", err)
			}
		}

		if err := k.Load(file.Provider(path), json.Parser()); err != nil {
			log.Fatalf("Error loading config: %v", err)
		}

		configManager = &ConfigManager{}

		var cfg Config
		if err := k.Unmarshal("", &cfg); err != nil {
			log.Fatalf("Error unmarshaling config: %v", err)
		}

		configManager.Config = &cfg

		// Validate all rule IDs
		if err := configManager.Config.validateAllRuleIDs(); err != nil {
			log.Fatalf("Invalid rule IDs found: %v", err)
		}

		if err := saveConfigUnsafe(configManager.Config); err != nil {
			log.Fatalf("Error saving config: %v", err)
		}
	})

	return configManager
}

// Used by unit test
func SetConfigFileDirOverride(path string) {
	configFileDirOverride = path
}

func Get() *Config {
	if configManager == nil {
		log.Fatal("Config not loaded. Call Load() first.")
	}

	configManager.mu.RLock()
	defer configManager.mu.RUnlock()

	return configManager.Config
}

func GetAccountConfig(accountName string) (*AccountCfg, error) {
	// Use koanf's in-memory structure to get the account configuration
	key := fmt.Sprintf("accounts.%s", accountName)
	if !k.Exists(key) {
		return nil, fmt.Errorf("account '%s' not found in configuration", accountName)
	}

	var account AccountCfg
	if err := k.Unmarshal(key, &account); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account configuration for '%s': %v", accountName, err)
	}

	return &account, nil
}

func defaultConfig() *Config {
	return &Config{
		Accounts: map[string]AccountCfg{},
	}
}

func configFilePath() string {
	if configFileDirOverride != "" {
		return filepath.Join(configFileDirOverride, fileName)
	}

	configDir, err := wallet.Path()
	if err != nil {
		log.Fatalf("Failed to get wallet config directory: %v", err)
	}

	return filepath.Join(configDir, fileName)
}

// Public method for saving config
func SaveConfig(config *Config) error {
	configManager.mu.Lock()
	defer configManager.mu.Unlock()

	return saveConfigUnsafe(config)
}

// saveConfigUnsafe saves the config without locking.
// It should only be called from methods that already hold the lock.
func saveConfigUnsafe(config *Config) error {
	// Check for duplicate IDs
	if err := config.ensureUniqueRuleIDs(); err != nil {
		return err
	}

	path := configFilePath()

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create or save config file: %v", err)
	}
	defer file.Close()

	encoder := jsonStd.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print JSON

	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to write config to file: %v", err)
	}

	return nil
}
