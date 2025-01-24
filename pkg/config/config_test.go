package config

import (
	"log"
	"os"
	"testing"

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
	rule := SignRule{
		Name:     "Test Rule",
		Contract: "AS12U4TZfNK7qoLyEERBBRDMu8nm5MKoRzPXDXans4v9wdATZedz9",
		RuleType: RuleTypeAutoSign,
		Enabled:  true,
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

	rule := SignRule{
		Name:     "Test Rule",
		Contract: contract,
		RuleType: RuleTypeAutoSign,
		Enabled:  true,
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
	rule := SignRule{
		Name:     "Test Rule",
		Contract: contract,
		RuleType: RuleTypeAutoSign,
		Enabled:  true,
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
	rule := SignRule{
		Name:     "Test Rule",
		Contract: contract,
		RuleType: RuleTypeAutoSign,
		Enabled:  true,
	}

	ruleId, err := cfg.AddSignRule(accountName, rule)
	assert.NoError(t, err)

	hasEnabled := cfg.HasEnabledRule(accountName)
	assert.True(t, hasEnabled)

	rulePtr := cfg.GetEnabledRuleForContract(accountName, &contract)
	assert.NotNil(t, rulePtr)

	// Disable the rule and check again
	rule.Enabled = false
	_, err = cfg.UpdateSignRule(accountName, ruleId, rule)
	assert.NoError(t, err)

	rulePtr = cfg.GetEnabledRuleForContract(accountName, &contract)
	assert.Nil(t, rulePtr)
}
