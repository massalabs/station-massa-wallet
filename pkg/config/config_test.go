package config

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/massalabs/station/pkg/logger"
	"github.com/stretchr/testify/assert"
)

var cfg *Config

const accountName = "test_account"

func TestMain(m *testing.M) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "*-wallet-dir")
	if err != nil {
		log.Fatalf("while creating temporary wallet directory: %s", err.Error())
	}

	// Load config file with config file path override
	SetConfigFileDirOverride(tempDir)
	manager := Load()
	cfg = manager.Config

	os.Exit(m.Run())
}

func TestAddSignRule(t *testing.T) {
	authorizedOrigin := "http://localhost:3000"
	rule := SignRule{
		Name:             "Test Rule",
		Contract:         "AS12U4TZfNK7qoLyEERBBRDMu8nm5MKoRzPXDXans4v9wdATZedz9",
		RuleType:         RuleTypeAutoSign,
		Enabled:          true,
		AuthorizedOrigin: &authorizedOrigin,
	}

	ruleID, err := cfg.AddSignRule(accountName, rule)
	assert.NoError(t, err)
	assert.NotEmpty(t, ruleID)

	// Verify the rule was added
	account, exists := cfg.Accounts[accountName]
	assert.True(t, exists)
	assert.Equal(t, 1, len(account.SignRules))
	assert.Equal(t, ruleID, account.SignRules[0].ID)

	// Add another rule
	rule2 := SignRule{
		Name:     "Test Rule 2",
		Contract: "AS133eqPPaPttJ6hJnk3sfoG5cjFFqBDi1VGxdo2wzWkq8AfZnan",
		RuleType: RuleTypeDisablePasswordPrompt,
		Enabled:  false,
	}

	ruleID, err = cfg.AddSignRule(accountName, rule2)
	assert.NoError(t, err)
	assert.NotEmpty(t, ruleID)

	// Verify the rule was added
	account, exists = cfg.Accounts[accountName]
	assert.True(t, exists)
	assert.Equal(t, 2, len(account.SignRules))
	assert.Equal(t, ruleID, account.SignRules[1].ID)

	// Add the same rule again and verify it fails (duplicate rule)
	_, err = cfg.AddSignRule(accountName, rule)
	assert.Error(t, err)

	_, err = cfg.AddSignRule(accountName, rule2)
	assert.Error(t, err)

	// check there is no duplicate rule
	err = cfg.ensureUniqueRuleIDs()
	assert.NoError(t, err)
}

func TestDeleteSignRule(t *testing.T) {
	contract := "AS12UMSUxgpRBB6ArZDJ19arHoxNkkpdfofQGekAiAJqsuE6PEFJy"
	authorizedOrigin := "http://localhost:3000"

	rule := SignRule{
		Name:             "Test Rule",
		Contract:         contract,
		RuleType:         RuleTypeAutoSign,
		Enabled:          true,
		AuthorizedOrigin: &authorizedOrigin,
	}

	ruleID, err := cfg.AddSignRule(accountName, rule)
	assert.NoError(t, err)

	err = cfg.DeleteSignRule(accountName, ruleID)
	assert.NoError(t, err)

	// Verify the rule was deleted
	deletedRule := cfg.GetSignRule(accountName, ruleID)
	assert.Nil(t, deletedRule)
}

func TestUpdateSignRule(t *testing.T) {
	accountName := "test_account"
	contract := "AS12UMSUxgpRBB6ArZDJ19arHoxNkkpdfofQGekAiAJqsuE6PEFJy"
	authorizedOrigin := "http://localhost:3000"
	rule := SignRule{
		Name:             "Test Rule",
		Contract:         contract,
		RuleType:         RuleTypeAutoSign,
		Enabled:          true,
		AuthorizedOrigin: &authorizedOrigin,
	}

	ruleID, err := cfg.AddSignRule(accountName, rule)
	assert.NoError(t, err)

	newRule := SignRule{
		Name:     "Updated Rule",
		Contract: "AS125oPLYRTtfVjpWisPZVTLjBhCFfQ1jDsi75XNtRm1NZux54eCj",
		RuleType: RuleTypeDisablePasswordPrompt,
		Enabled:  false,
	}

	newRuleID, err := cfg.UpdateSignRule(accountName, ruleID, newRule)
	assert.NoError(t, err)
	assert.NotEqual(t, ruleID, newRuleID)

	// Verify the rule was updated
	updatedRule := cfg.GetSignRule(accountName, newRuleID)
	assert.NotNil(t, updatedRule)
	assert.Equal(t, newRuleID, updatedRule.ID)
	assert.Equal(t, newRule.Name, updatedRule.Name)
	assert.Equal(t, newRule.Contract, updatedRule.Contract)
	assert.Equal(t, newRule.RuleType, updatedRule.RuleType)
	assert.Equal(t, newRule.Enabled, updatedRule.Enabled)

	// Verify previous rule is deleted
	deletedRule := cfg.GetSignRule(accountName, ruleID)
	assert.Nil(t, deletedRule)
}

func TestValidateRuleID(t *testing.T) {
	rule := SignRule{
		Name:     "Test Rule",
		Contract: "AS1hCJXjndR4c9vekLWsXGnrdigp4AaZ7uYG3UKFzzKnWVsrNLPJ",
		RuleType: RuleTypeAutoSign,
		Enabled:  true,
	}

	rule.ID = generateRuleID(accountName, rule)
	err := ValidateRuleID(accountName, rule)
	assert.NoError(t, err)

	// Test with an invalid ID
	rule.ID = "invalid_id"
	err = ValidateRuleID(accountName, rule)
	assert.Error(t, err)
}

func TestHasEnabledRule(t *testing.T) {
	accountName := "test_account"
	contract := "AS124vf3YfAJCSCQVYKczzuWWpXrximFpbTmX4rheLs5uNSftiiRY"
	authorizedOrigin := "http://localhost:3000"

	rule := SignRule{
		Name:             "Test Rule",
		Contract:         contract,
		RuleType:         RuleTypeAutoSign,
		Enabled:          true,
		AuthorizedOrigin: &authorizedOrigin,
	}

	ruleId, err := cfg.AddSignRule(accountName, rule)
	assert.NoError(t, err)

	hasEnabled := cfg.HasEnabledRule(accountName)
	assert.True(t, hasEnabled)

	rulePtr := cfg.GetEnabledRuleForContract(accountName, &contract, &authorizedOrigin)
	assert.NotNil(t, rulePtr)

	// Disable the rule and check again
	rule.Enabled = false
	_, err = cfg.UpdateSignRule(accountName, ruleId, rule)
	assert.NoError(t, err)

	rulePtr = cfg.GetEnabledRuleForContract(accountName, &contract, &authorizedOrigin)
	assert.Nil(t, rulePtr)
}

func TestLegacyConfigHandling(t *testing.T) {
	tempdir := t.TempDir()
	if err := logger.InitializeGlobal(filepath.Join(tempdir, "unit-test.log")); err != nil {
		log.Fatalf("while initializing global logger: %s", err.Error())
	}

	t.Run("set default rule timeout when zero", func(t *testing.T) {
		legacyConfig := &Config{
			RuleTimeout: 0, // Legacy config with no timeout
			Accounts: map[string]AccountCfg{
				"test_account": {
					SignRules: []SignRule{},
				},
			},
		}

		err := legacyConfig.legacyConfigHandling()
		assert.NoError(t, err)
		assert.Equal(t, DefaultRuleTimeout, legacyConfig.RuleTimeout)
	})

	// Test case 2: Legacy AutoSign rule with no authorized origin should be deleted
	t.Run("delete legacy autosign rule with no authorized origin", func(t *testing.T) {
		legacyConfig := &Config{
			RuleTimeout: DefaultRuleTimeout,
			Accounts: map[string]AccountCfg{
				"test_account": {
					SignRules: []SignRule{
						{
							ID:       "legacy_rule_1",
							Name:     "Legacy AutoSign Rule",
							Contract: "AS12U4TZfNK7qoLyEERBBRDMu8nm5MKoRzPXDXans4v9wdATZedz9",
							RuleType: RuleTypeAutoSign,
							Enabled:  true,
							// No AuthorizedOrigin - this should be deleted
						},
						{
							ID:               "valid_rule_1",
							Name:             "Valid AutoSign Rule",
							Contract:         "AS133eqPPaPttJ6hJnk3sfoG5cjFFqBDi1VGxdo2wzWkq8AfZnan",
							RuleType:         RuleTypeAutoSign,
							Enabled:          true,
							AuthorizedOrigin: func() *string { s := "http://localhost:3000"; return &s }(),
						},
						{
							ID:       "disable_pwd_rule_1",
							Name:     "Disable Password Rule",
							Contract: "AS1hCJXjndR4c9vekLWsXGnrdigp4AaZ7uYG3UKFzzKnWVsrNLPJ",
							RuleType: RuleTypeDisablePasswordPrompt,
							Enabled:  true,
							// No AuthorizedOrigin - but this is not AutoSign, so it should remain
						},
					},
				},
			},
		}

		err := legacyConfig.legacyConfigHandling()
		assert.NoError(t, err)

		// Verify the legacy AutoSign rule was deleted
		account := legacyConfig.Accounts["test_account"]
		assert.Equal(t, 2, len(account.SignRules)) // Should have 2 rules left

		// Check that the valid AutoSign rule and DisablePasswordPrompt rule remain
		ruleNames := make([]string, len(account.SignRules))
		for i, rule := range account.SignRules {
			ruleNames[i] = rule.Name
		}

		assert.Contains(t, ruleNames, "Valid AutoSign Rule")
		assert.Contains(t, ruleNames, "Disable Password Rule")
		assert.NotContains(t, ruleNames, "Legacy AutoSign Rule")
	})

	// Test case 3: Legacy rules with no expiration date should get default expiration
	t.Run("add default expiration date to legacy rules", func(t *testing.T) {
		legacyConfig := &Config{
			RuleTimeout: DefaultRuleTimeout,
			Accounts: map[string]AccountCfg{
				"test_account": {
					SignRules: []SignRule{
						{
							ID:       "legacy_rule_2",
							Name:     "Legacy Rule No Expiration",
							Contract: "AS12U4TZfNK7qoLyEERBBRDMu8nm5MKoRzPXDXans4v9wdATZedz9",
							RuleType: RuleTypeDisablePasswordPrompt,
							Enabled:  true,
							// ExpireAfter is zero - should get default expiration
						},
						{
							ID:               "valid_rule_2",
							Name:             "Valid Rule With Expiration",
							Contract:         "AS133eqPPaPttJ6hJnk3sfoG5cjFFqBDi1VGxdo2wzWkq8AfZnan",
							RuleType:         RuleTypeAutoSign,
							Enabled:          true,
							AuthorizedOrigin: func() *string { s := "http://localhost:3000"; return &s }(),
							ExpireAfter:      time.Now().Add(24 * time.Hour), // Already has expiration
						},
					},
				},
			},
		}

		err := legacyConfig.legacyConfigHandling()
		assert.NoError(t, err)

		// Verify the legacy rule got an expiration date
		account := legacyConfig.Accounts["test_account"]
		assert.Equal(t, 2, len(account.SignRules))

		// Find the legacy rule
		var legacyRule *SignRule

		for i := range account.SignRules {
			if account.SignRules[i].Name == "Legacy Rule No Expiration" {
				legacyRule = &account.SignRules[i]
				break
			}
		}

		assert.NotNil(t, legacyRule)
		assert.False(t, legacyRule.ExpireAfter.IsZero())
		assert.True(t, legacyRule.ExpireAfter.After(time.Now()))
		assert.True(t, legacyRule.ExpireAfter.Before(time.Now().Add(time.Duration(DefaultRuleTimeout+1)*time.Second)))
	})

	// Test case 4: Combined scenario - multiple legacy issues
	t.Run("handle multiple legacy issues", func(t *testing.T) {
		legacyConfig := &Config{
			RuleTimeout: 0, // No timeout
			Accounts: map[string]AccountCfg{
				"test_account": {
					SignRules: []SignRule{
						{
							ID:       "legacy_autosign_no_origin",
							Name:     "Legacy AutoSign No Origin",
							Contract: "AS12U4TZfNK7qoLyEERBBRDMu8nm5MKoRzPXDXans4v9wdATZedz9",
							RuleType: RuleTypeAutoSign,
							Enabled:  true,
							// No AuthorizedOrigin - should be deleted
						},
						{
							ID:       "legacy_no_expiration",
							Name:     "Legacy No Expiration",
							Contract: "AS133eqPPaPttJ6hJnk3sfoG5cjFFqBDi1VGxdo2wzWkq8AfZnan",
							RuleType: RuleTypeDisablePasswordPrompt,
							Enabled:  true,
							// No expiration - should get default
						},
						{
							ID:               "valid_rule",
							Name:             "Valid Rule",
							Contract:         "AS1hCJXjndR4c9vekLWsXGnrdigp4AaZ7uYG3UKFzzKnWVsrNLPJ",
							RuleType:         RuleTypeAutoSign,
							Enabled:          true,
							AuthorizedOrigin: func() *string { s := "http://localhost:3000"; return &s }(),
							ExpireAfter:      time.Now().Add(24 * time.Hour),
						},
					},
				},
			},
		}

		err := legacyConfig.legacyConfigHandling()
		assert.NoError(t, err)

		// Verify RuleTimeout was set
		assert.Equal(t, DefaultRuleTimeout, legacyConfig.RuleTimeout)

		// Verify the legacy AutoSign rule was deleted
		account := legacyConfig.Accounts["test_account"]
		assert.Equal(t, 2, len(account.SignRules)) // Should have 2 rules left

		// Check that the valid rule and the rule with no expiration remain
		ruleNames := make([]string, len(account.SignRules))
		for i, rule := range account.SignRules {
			ruleNames[i] = rule.Name
		}

		assert.Contains(t, ruleNames, "Valid Rule")
		assert.Contains(t, ruleNames, "Legacy No Expiration")
		assert.NotContains(t, ruleNames, "Legacy AutoSign No Origin")

		// Verify the legacy rule got an expiration date
		var legacyRule *SignRule

		for i := range account.SignRules {
			if account.SignRules[i].Name == "Legacy No Expiration" {
				legacyRule = &account.SignRules[i]
				break
			}
		}

		assert.NotNil(t, legacyRule)
		assert.False(t, legacyRule.ExpireAfter.IsZero())
		assert.True(t, legacyRule.ExpireAfter.After(time.Now()))
	})
}
