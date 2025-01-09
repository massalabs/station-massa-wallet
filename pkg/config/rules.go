package config

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func (c *Config) AddSignRule(accountName string, rule SignRule) (string, error) {
	if err := validateRule(rule); err != nil {
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
	_, err := c.getSignRuleUnsafe(accountName, rule.ID)
	if err == nil {
		return "", fmt.Errorf("rule %s already exists", rule.ID)
	}

	// Add the new rule
	account.SignRules = append(account.SignRules, rule)
	c.Accounts[accountName] = account

	err = saveConfigUnsafe(c)
	if err != nil {
		return "", fmt.Errorf("error saving configuration: %v", err)
	}

	return rule.ID, nil
}

func (c *Config) DeleteSignRule(accountName, ruleID string) error {
	configManager.mu.Lock()
	defer configManager.mu.Unlock()

	// Check if the account exists
	account, exists := c.Accounts[accountName]
	if !exists {
		return fmt.Errorf("account not found: %s", accountName)
	}

	// Find the index of the rule to delete
	index := -1

	for i, rule := range account.SignRules {
		if rule.ID == ruleID {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("rule not found: %s", ruleID)
	}

	// Remove the rule from the slice
	account.SignRules = append(account.SignRules[:index], account.SignRules[index+1:]...)
	c.Accounts[accountName] = account

	err := saveConfigUnsafe(c)
	if err != nil {
		return fmt.Errorf("error saving configuration: %v", err)
	}

	return nil
}

func (c *Config) UpdateSignRule(accountName, ruleID string, newRule SignRule) (string, error) {
	configManager.mu.Lock()
	defer configManager.mu.Unlock()

	// Check if the account exists
	account, exists := c.Accounts[accountName]
	if !exists {
		return "", fmt.Errorf("account not found: %s", accountName)
	}

	// Find the index of the rule to update
	index := -1

	for i, rule := range account.SignRules {
		if rule.ID == ruleID {
			index = i
			break
		}
	}

	if index == -1 {
		return "", fmt.Errorf("rule not found: %s", ruleID)
	}

	// Delete the existing rule
	account.SignRules = append(account.SignRules[:index], account.SignRules[index+1:]...)
	c.Accounts[accountName] = account

	// Add the new rule
	newRule.ID = generateRuleID(accountName, newRule)
	account.SignRules = append(account.SignRules, newRule)
	c.Accounts[accountName] = account

	err := saveConfigUnsafe(c)
	if err != nil {
		return "", fmt.Errorf("error saving configuration: %v", err)
	}

	return newRule.ID, nil
}

func validateRule(rule SignRule) error {
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

func (c *Config) getSignRuleUnsafe(accountName, ruleID string) (*SignRule, error) {
	// Check if the account exists
	account, exists := c.Accounts[accountName]
	if !exists {
		return nil, fmt.Errorf("account not found: %s", accountName)
	}

	// Search for the rule by ID
	for _, rule := range account.SignRules {
		if rule.ID == ruleID {
			return &rule, nil
		}
	}

	return nil, fmt.Errorf("rule not found: %s", ruleID)
}

func (c *Config) GetSignRule(accountName, ruleID string) (*SignRule, error) {
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

	// Check if there is any enabled rule that applies to the contract
	for _, rule := range account.SignRules {
		if rule.Enabled {
			switch rule.RuleType {
			case RuleTypeAutoSign:
				if contract != nil && rule.Contract == *contract {
					return &rule.RuleType
				}

			case RuleTypeDisablePasswordPrompt:
				if rule.Contract == "*" || (contract != nil && rule.Contract == *contract) {
					return &rule.RuleType
				}
			}
		}
	}

	return nil
}