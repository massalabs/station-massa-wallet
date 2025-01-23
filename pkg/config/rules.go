package config

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/massalabs/station-massa-wallet/pkg/utils"
)

func (c *Config) AddSignRule(accountName string, rule SignRule) (string, error) {
	if err := ValidateRule(rule); err != nil {
		return "", fmt.Errorf("invalid rule: %v", err)
	}

	rule.ID = generateRuleID(accountName, rule)

	configManager.mu.Lock()
	defer configManager.mu.Unlock()

	// Check if the account exists
	account, exists := c.Accounts[accountName]
	if !exists {
		account = AccountCfg{SignRules: []SignRule{}}
		c.Accounts[accountName] = account
	}

	// check if the rule already exists
	existingRule := c.getSignRuleUnsafe(accountName, rule.ID)
	if existingRule != nil {
		return "", fmt.Errorf("rule %s already exists", rule.ID)
	}

	// Add the new rule
	account.SignRules = append(account.SignRules, rule)
	c.Accounts[accountName] = account

	if err := saveConfigUnsafe(c); err != nil {
		return "", fmt.Errorf("error saving configuration: %v", err)
	}

	return rule.ID, nil
}

func (c *Config) findRuleIndex(accountConfig AccountCfg, ruleID string) (int, error) {
	for i, rule := range accountConfig.SignRules {
		if rule.ID == ruleID {
			return i, nil
		}
	}

	return -1, fmt.Errorf("rule not found: %s", ruleID)
}

func (c *Config) DeleteSignRule(accountName, ruleID string) error {
	configManager.mu.Lock()
	defer configManager.mu.Unlock()

	// Check if the account exists
	accountConfig, exists := c.Accounts[accountName]
	if !exists {
		return fmt.Errorf("account not found: %s", accountName)
	}

	index, err := c.findRuleIndex(accountConfig, ruleID)
	if err != nil {
		return fmt.Errorf("deleting sign rule: %w", err)
	}

	// Remove the rule from the slice
	accountConfig.SignRules = append(accountConfig.SignRules[:index], accountConfig.SignRules[index+1:]...)
	c.Accounts[accountName] = accountConfig

	err = saveConfigUnsafe(c)
	if err != nil {
		return fmt.Errorf("error saving configuration: %v", err)
	}

	return nil
}

func (c *Config) UpdateSignRule(accountName, ruleID string, newRule SignRule) (string, error) {
	if err := ValidateRule(newRule); err != nil {
		return "", fmt.Errorf("invalid rule: %v", err)
	}

	configManager.mu.Lock()
	defer configManager.mu.Unlock()

	// Check if the account exists
	accountConfig, exists := c.Accounts[accountName]
	if !exists {
		return "", fmt.Errorf("account not found: %s", accountName)
	}

	index, err := c.findRuleIndex(accountConfig, ruleID)
	if err != nil {
		return "", fmt.Errorf("updating sign rule: %w", err)
	}

	// Delete the existing rule
	accountConfig.SignRules = append(accountConfig.SignRules[:index], accountConfig.SignRules[index+1:]...)
	c.Accounts[accountName] = accountConfig

	// Add the new rule
	newRule.ID = generateRuleID(accountName, newRule)
	accountConfig.SignRules = append(accountConfig.SignRules, newRule)
	c.Accounts[accountName] = accountConfig

	err = saveConfigUnsafe(c)
	if err != nil {
		return "", fmt.Errorf("error saving configuration: %v", err)
	}

	return newRule.ID, nil
}

func ValidateRule(rule SignRule) error {
	if rule.Contract != "*" && !utils.IsValidContract(rule.Contract) {
		return fmt.Errorf("invalid contract address: %s", rule.Contract)
	}

	switch rule.RuleType {
	case RuleTypeDisablePasswordPrompt:
		return nil

	case RuleTypeAutoSign:
		if rule.Contract == "*" {
			return fmt.Errorf("RuleTypeAutoSign cannot have a wildcard contract")
		}

		return nil

	default:
		return fmt.Errorf("invalid ruleType: %s", rule.RuleType)
	}
}

// Two rules with the same SignRule will always generate the same ID.
func generateRuleID(accountName string, rule SignRule) string {
	// Create a string that combines all rule parameters
	ruleString := fmt.Sprintf("%s:%s:%s", accountName, rule.Contract, string(rule.RuleType))

	hash := sha256.Sum256([]byte(ruleString))

	// Convert the first 10 bytes of the hash to a hex string
	return hex.EncodeToString(hash[:5])
}

func (c *Config) ensureUniqueRuleIDs() error {
	idMap := make(map[string]bool)

	for _, account := range c.Accounts {
		for _, rule := range account.SignRules {
			if idMap[rule.ID] {
				return fmt.Errorf("duplicate ID found: %s", rule.ID)
			}

			idMap[rule.ID] = true
		}
	}

	return nil
}

func ValidateRuleID(accountName string, rule SignRule) error {
	// Generate the expected ID based on the rule's parameters
	expectedID := generateRuleID(accountName, rule)

	// Compare the expected ID with the actual ID
	if rule.ID != expectedID {
		return fmt.Errorf("rule ID mismatch: expected %s, got %s", expectedID, rule.ID)
	}

	return nil
}

func (c *Config) validateAllRuleIDs() error {
	for accountName, account := range c.Accounts {
		for i, rule := range account.SignRules {
			if err := ValidateRuleID(accountName, rule); err != nil {
				return fmt.Errorf("invalid rule ID in account '%s', rule index %d: %v", accountName, i, err)
			}
		}
	}

	return nil
}

func (c *Config) getSignRuleUnsafe(accountName, ruleID string) *SignRule {
	// Check if the account exists
	account, exists := c.Accounts[accountName]
	if !exists {
		return nil
	}

	// Search for the rule by ID
	for _, rule := range account.SignRules {
		if rule.ID == ruleID {
			return &rule
		}
	}

	return nil
}

func (c *Config) GetSignRule(accountName, ruleID string) *SignRule {
	configManager.mu.RLock()
	defer configManager.mu.RUnlock()

	return c.getSignRuleUnsafe(accountName, ruleID)
}

func (c *Config) HasEnabledRule(accountName string) bool {
	configManager.mu.RLock()
	defer configManager.mu.RUnlock()

	// Check if the account exists
	account, exists := c.Accounts[accountName]
	if !exists {
		return false
	}

	// Check if there is any enabled rule
	for _, rule := range account.SignRules {
		if rule.Enabled {
			return true
		}
	}

	return false
}

func (c *Config) GetEnabledRuleForContract(accountName string, contract *string) *RuleType {
	configManager.mu.RLock()
	defer configManager.mu.RUnlock()

	// Check if the account exists
	account, exists := c.Accounts[accountName]
	if !exists {
		return nil
	}

	var ruleType *RuleType = nil

	// Check if there is any enabled rule that applies to the contract
	for _, rule := range account.SignRules {
		if rule.Enabled {
			switch rule.RuleType {
			case RuleTypeAutoSign:
				if contract != nil && rule.Contract == *contract {
					ruleType = &rule.RuleType
				}

			case RuleTypeDisablePasswordPrompt:
				if rule.Contract == "*" || (contract != nil && rule.Contract == *contract) {
					if ruleType == nil {
						// If there are multiple rules that apply to the contract, the rule with the highest priority is used (AutoSign)
						ruleType = &rule.RuleType
					}
				}
			}
		}
	}

	return ruleType
}

func (c *Config) IsExistingRule(accountName string, rule SignRule) bool {
	configManager.mu.RLock()
	defer configManager.mu.RUnlock()

	// Generate the expected ID based on the rule's parameters
	expectedID := generateRuleID(accountName, rule)

	rulePtr := c.getSignRuleUnsafe(accountName, expectedID)

	return rulePtr != nil
}
