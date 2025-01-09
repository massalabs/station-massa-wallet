package config

import "sync"

type RuleType string

const (
	RuleTypeDisablePasswordPrompt RuleType = "disable_password_prompt"
	RuleTypeAutoSign              RuleType = "auto_sign"
)

type SignRule struct {
	Name     string   `koanf:"name"`
	ID       string   `koanf:"id"`
	Contract string   `koanf:"contract"`
	RuleType RuleType `koanf:"ruleType"`
	Enabled  bool     `koanf:"enabled"`
}

type AccountCfg struct {
	SignRules []SignRule `koanf:"signRules"`
}

type Config struct {
	Accounts map[string]AccountCfg `koanf:"accounts"`
}
type ConfigManager struct {
	mu     sync.RWMutex
	Config *Config
}
