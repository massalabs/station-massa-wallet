package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func strptr(s string) *string { return &s }

// Sample valid contract addresses taken from existing tests/fixtures in the repo
const (
	validContract1 = "AS12U4TZfNK7qoLyEERBBRDMu8nm5MKoRzPXDXans4v9wdATZedz9"
	validContract2 = "AS133eqPPaPttJ6hJnk3sfoG5cjFFqBDi1VGxdo2wzWkq8AfZnan"
)

func TestValidateRule(t *testing.T) {
	empty := ""

	cases := []struct {
		name        string
		rule        SignRule
		wantErr     bool
		errContains string
	}{
		{name: "Disable Password allows wildcard", rule: SignRule{Contract: "*", RuleType: RuleTypeDisablePasswordPrompt}},
		{name: "Disable Password allows specific contract", rule: SignRule{Contract: validContract1, RuleType: RuleTypeDisablePasswordPrompt}},
		{name: "Disable Password rejects invalid contract", rule: SignRule{Contract: "AA_invalid", RuleType: RuleTypeDisablePasswordPrompt}, wantErr: true, errContains: "invalid contract address"},

		{name: "Disable Password valid origin host", rule: SignRule{Contract: validContract1, RuleType: RuleTypeDisablePasswordPrompt, AuthorizedOrigin: strptr("https://example.com")}},
		{name: "Disable Password valid origin host+port", rule: SignRule{Contract: validContract1, RuleType: RuleTypeDisablePasswordPrompt, AuthorizedOrigin: strptr("https://example.com:8080")}},
		{name: "Disable Password invalid origin with path", rule: SignRule{Contract: validContract1, RuleType: RuleTypeDisablePasswordPrompt, AuthorizedOrigin: strptr("https://example.com/path")}, wantErr: true, errContains: "invalid AuthorizedOrigin URL"},
		{name: "Disable Password invalid blank origin", rule: SignRule{Contract: validContract1, RuleType: RuleTypeDisablePasswordPrompt, AuthorizedOrigin: strptr(" ")}, wantErr: true, errContains: "invalid AuthorizedOrigin URL"},

		{name: "AutoSign rejects wildcard contract", rule: SignRule{Contract: "*", RuleType: RuleTypeAutoSign, AuthorizedOrigin: strptr("https://example.com")}, wantErr: true, errContains: "cannot have a wildcard contract"},
		{name: "AutoSign requires non-empty origin (nil)", rule: SignRule{Contract: validContract2, RuleType: RuleTypeAutoSign}, wantErr: true, errContains: "must have a non-empty AuthorizedOrigin"},
		{name: "AutoSign requires non-empty origin (empty)", rule: SignRule{Contract: validContract2, RuleType: RuleTypeAutoSign, AuthorizedOrigin: &empty}, wantErr: true, errContains: "must have a non-empty AuthorizedOrigin"},
		{name: "AutoSign accepts valid origin", rule: SignRule{Contract: validContract2, RuleType: RuleTypeAutoSign, AuthorizedOrigin: strptr("https://example.org")}},
		{name: "AutoSign rejects malformed origin", rule: SignRule{Contract: validContract2, RuleType: RuleTypeAutoSign, AuthorizedOrigin: strptr("http://example.org/path?x=1")}, wantErr: true, errContains: "invalid AuthorizedOrigin URL"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ValidateRule(c.rule)
			if c.wantErr {
				require.Error(t, err)

				if c.errContains != "" {
					require.Contains(t, err.Error(), c.errContains)
				}

				return
			}

			require.NoError(t, err)
		})
	}
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

func TestGetEnabledRuleForContract(t *testing.T) {
	origin1 := "https://example.com"
	origin2 := "https://other.com"
	c1 := validContract1
	c2 := validContract2

	cases := []struct {
		name        string
		rules       []SignRule
		contract    *string
		origin      *string
		expectNil   bool
		expectType  RuleType
		expectContr string
	}{
		{
			name:      "no rules returns nil",
			rules:     []SignRule{},
			contract:  &c1,
			origin:    &origin1,
			expectNil: true,
		},
		{
			name:        "DisabledPassword wildcard with matching origin",
			rules:       []SignRule{{Contract: "*", RuleType: RuleTypeDisablePasswordPrompt, Enabled: true, AuthorizedOrigin: &origin1}},
			contract:    &c1,
			origin:      &origin1,
			expectNil:   false,
			expectType:  RuleTypeDisablePasswordPrompt,
			expectContr: "*",
		},
		{
			name:        "DisabledPassword specific contract match, all origins",
			rules:       []SignRule{{Contract: c1, RuleType: RuleTypeDisablePasswordPrompt, Enabled: true, AuthorizedOrigin: nil}},
			contract:    &c1,
			origin:      &origin1,
			expectNil:   false,
			expectType:  RuleTypeDisablePasswordPrompt,
			expectContr: c1,
		},
		{
			name:        "DisabledPassword all contracts and all origins",
			rules:       []SignRule{{Contract: "*", RuleType: RuleTypeDisablePasswordPrompt, Enabled: true, AuthorizedOrigin: nil}},
			contract:    nil,
			origin:      &origin1,
			expectNil:   false,
			expectType:  RuleTypeDisablePasswordPrompt,
			expectContr: "*",
		},
		{
			name:        "DisabledPassword specific contracts and specific origins",
			rules:       []SignRule{{Contract: c1, RuleType: RuleTypeDisablePasswordPrompt, Enabled: true, AuthorizedOrigin: &origin1}},
			contract:    &c1,
			origin:      &origin1,
			expectNil:   false,
			expectType:  RuleTypeDisablePasswordPrompt,
			expectContr: c1,
		},
		{
			name:      "DisabledPassword Origin mismatch returns nil",
			rules:     []SignRule{{Contract: "*", RuleType: RuleTypeDisablePasswordPrompt, Enabled: true, AuthorizedOrigin: &origin1}},
			contract:  &c1,
			origin:    &origin2,
			expectNil: true,
		},
		{
			name:      "DisabledPassword rule is ignored if disabled",
			rules:     []SignRule{{Contract: "*", RuleType: RuleTypeDisablePasswordPrompt, Enabled: false, AuthorizedOrigin: &origin1}},
			contract:  &c1,
			origin:    &origin1,
			expectNil: true,
		},
		{
			name:        "AutoSign exact contract and origin",
			rules:       []SignRule{{Contract: c1, RuleType: RuleTypeAutoSign, Enabled: true, AuthorizedOrigin: &origin1}},
			contract:    &c1,
			origin:      &origin1,
			expectNil:   false,
			expectType:  RuleTypeAutoSign,
			expectContr: c1,
		},
		{
			name: "Both DisabledPassword and AutoSign present, AutoSign has priority",
			rules: []SignRule{
				{Contract: "*", RuleType: RuleTypeDisablePasswordPrompt, Enabled: true, AuthorizedOrigin: &origin1},
				{Contract: c1, RuleType: RuleTypeAutoSign, Enabled: true, AuthorizedOrigin: &origin1},
			},
			contract:    &c1,
			origin:      &origin1,
			expectNil:   false,
			expectType:  RuleTypeAutoSign,
			expectContr: c1,
		},
		{
			name: "DisabledPassword matches when AutoSign origin mismatches",
			rules: []SignRule{
				{Contract: c1, RuleType: RuleTypeAutoSign, Enabled: true, AuthorizedOrigin: &origin2},
				{Contract: "*", RuleType: RuleTypeDisablePasswordPrompt, Enabled: true, AuthorizedOrigin: &origin1},
			},
			contract:    &c1,
			origin:      &origin1,
			expectNil:   false,
			expectType:  RuleTypeDisablePasswordPrompt,
			expectContr: "*",
		},
		{
			name: "DisabledPassword matches when AutoSign contract mismatches",
			rules: []SignRule{
				{Contract: c1, RuleType: RuleTypeAutoSign, Enabled: true, AuthorizedOrigin: &origin1},
				{Contract: c2, RuleType: RuleTypeDisablePasswordPrompt, Enabled: true, AuthorizedOrigin: &origin1},
			},
			contract:    &c2,
			origin:      &origin1,
			expectNil:   false,
			expectType:  RuleTypeDisablePasswordPrompt,
			expectContr: c2,
		},
		{
			name: "DisabledPassword matches when AutoSign is disabled",
			rules: []SignRule{
				{Contract: c1, RuleType: RuleTypeAutoSign, Enabled: false, AuthorizedOrigin: &origin1},
				{Contract: c1, RuleType: RuleTypeDisablePasswordPrompt, Enabled: true, AuthorizedOrigin: &origin1},
			},
			contract:    &c1,
			origin:      &origin1,
			expectNil:   false,
			expectType:  RuleTypeDisablePasswordPrompt,
			expectContr: c1,
		},
		{
			name:      "No match when contract different for AutoSign",
			rules:     []SignRule{{Contract: c2, RuleType: RuleTypeAutoSign, Enabled: true, AuthorizedOrigin: &origin1}},
			contract:  &c1,
			origin:    &origin1,
			expectNil: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// reset account rules for isolation
			cfg.Accounts[accountName] = AccountCfg{SignRules: append([]SignRule(nil), c.rules...)}

			got := cfg.GetEnabledRuleForContract(accountName, c.contract, c.origin)
			if c.expectNil {
				require.Nil(t, got)
				return
			}

			require.NotNil(t, got)
			require.Equal(t, c.expectType, got.RuleType)
			require.Equal(t, c.expectContr, got.Contract)
		})
	}
}
