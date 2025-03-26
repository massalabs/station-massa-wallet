package config

import "sync"

type EventType string

const (
	PromptResultEvent  string = "PROMPT_RESULT"
	PromptDataEvent    string = "PROMPT_DATA"
	PromptRequestEvent string = "PROMPT_REQUEST"
)

var EventTypes = []struct {
	Value  EventType
	TSName string
}{
	{EventType(PromptResultEvent), "promptResult"},
	{EventType(PromptDataEvent), "promptData"},
	{EventType(PromptRequestEvent), "promptRequest"},
}

type RuleType string

const (
	RuleTypeDisablePasswordPrompt RuleType = "DISABLE_PASSWORD_PROMPT"
	RuleTypeAutoSign              RuleType = "AUTO_SIGN"
)

var RuleTypes = []struct {
	Value  RuleType
	TSName string
}{
	{RuleTypeDisablePasswordPrompt, "disable_password_prompt"},
	{RuleTypeAutoSign, "auto_sign"},
}

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
