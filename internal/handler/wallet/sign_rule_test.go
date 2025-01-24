package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/massalabs/station-massa-wallet/api/server/models"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/cache"
	"github.com/massalabs/station-massa-wallet/pkg/config"
	"github.com/stretchr/testify/assert"
)

func Test_signrule_Handlers(t *testing.T) {
	assert := assert.New(t)

	api, resChan, err := MockAPI()
	assert.NoError(err)

	addhandler, exist := api.HandlerFor("post", "/api/accounts/{nickname}/signrules")
	assert.True(exist)

	updatehandler, exist := api.HandlerFor("put", "/api/accounts/{nickname}/signrules/{ruleId}")
	assert.True(exist)

	deletehandler, exist := api.HandlerFor("delete", "/api/accounts/{nickname}/signrules/{ruleId}")
	assert.True(exist)

	nickname := "zeWahLetName"
	password := "zePassword"
	account := createAccount(password, nickname, t, prompterAppMock)

	cfg := config.Get()

	t.Run("Add sign rule with invalid contract", func(t *testing.T) {
		contract := "Invalid Contract"
		ruleType := config.RuleTypeAutoSign
		ruleName := "Test Rule"
		enabled := true
		body := fmt.Sprintf(`{
			"name": "%s",
			"contract": "%s",
			"ruleType": "%s",
			"enabled": %t,
			"description": "Test Description"
		}`, ruleName, contract, ruleType, enabled)

		resp, err := handleHTTPRequest(addhandler, "POST", fmt.Sprintf("/api/accounts/%s/signrules", nickname), body)
		assert.NoError(err)

		verifyStatusCode(t, resp, http.StatusBadRequest)

		// check rule is not added
		rulePtr := cfg.GetEnabledRuleForContract(account.Nickname, &contract)
		assert.Nil(rulePtr)

		hasRule := cfg.HasEnabledRule(account.Nickname)
		assert.False(hasRule)
	})

	t.Run("Add sign rule", func(t *testing.T) {
		contract := "AS12UMSUxgpRBB6ArZDJ19arHoxNkkpdfofQGekAiAJqsuE6PEFJy"
		ruleType := config.RuleTypeAutoSign
		ruleName := "Test Rule"
		enabled := true
		body := fmt.Sprintf(`{
			"name": "%s",
			"contract": "%s",
			"ruleType": "%s",
			"enabled": %t,
			"description": "Test Description"
		}`, ruleName, contract, ruleType, enabled)

		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.StringPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Message:     password,
			}

			res <- (<-resChan)
		}(testResult)

		resp, err := handleHTTPRequest(addhandler, "POST", fmt.Sprintf("/api/accounts/%s/signrules", nickname), body)
		assert.NoError(err)

		verifyStatusCode(t, resp, http.StatusOK)

		result := <-testResult

		checkResultChannel(t, result, true, "")

		var addRuleResponse models.AddSignRuleResponse
		err = json.Unmarshal(resp.Body.Bytes(), &addRuleResponse)
		assert.NoError(err)
		assert.NotEmpty(addRuleResponse.ID)

		// check rule is added
		rulePtr := cfg.GetEnabledRuleForContract(account.Nickname, &contract)
		assert.NotNil(rulePtr)
		assert.Equal(*rulePtr, ruleType)

		hasRule := cfg.HasEnabledRule(account.Nickname)
		assert.True(hasRule)

		// check rule contains expected fields
		rule := cfg.GetSignRule(account.Nickname, addRuleResponse.ID)
		assert.NotNil(rule)
		assert.Equal(rule.Name, ruleName)
		assert.Equal(rule.Contract, contract)
		assert.Equal(rule.RuleType, ruleType)
		assert.Equal(rule.Enabled, enabled)

		// check that privateKey is cached
		pkey, err := cache.PrivateKeyFromCache(account)
		assert.NoError(err)
		assert.NotNil(pkey)
	})

	t.Run("Update sign rule", func(t *testing.T) {
		contract := "*"
		ruleType := config.RuleTypeDisablePasswordPrompt
		ruleName := "Test Rule"
		enabled := true
		body := fmt.Sprintf(`{
			"name": "%s",
			"contract": "%s",
			"ruleType": "%s",
			"enabled": %t,
			"description": "Test Description"
		}`, ruleName, contract, ruleType, enabled)

		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.StringPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Message:     password,
			}

			res <- (<-resChan)
		}(testResult)

		// Add a sign rule first
		resp, err := handleHTTPRequest(addhandler, "POST", fmt.Sprintf("/api/accounts/%s/signrules", nickname), body)
		assert.NoError(err)
		verifyStatusCode(t, resp, http.StatusOK)

		var addRuleResponse models.AddSignRuleResponse
		err = json.Unmarshal(resp.Body.Bytes(), &addRuleResponse)
		assert.NoError(err)
		assert.NotEmpty(addRuleResponse.ID)

		// Update the sign rule
		updatedRuleName := "Updated Test Rule"
		enabled = false
		updatedBody := fmt.Sprintf(`{
			"name": "%s",
			"contract": "%s",
			"ruleType": "%s",
			"enabled": %t,
			"description": "Updated Test Description"
		}`, updatedRuleName, contract, ruleType, enabled)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.StringPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Message:     password,
			}

			res <- (<-resChan)
		}(testResult)

		resp, err = handleHTTPRequest(updatehandler, "PUT", fmt.Sprintf("/api/accounts/%s/signrules/%s", nickname, addRuleResponse.ID), updatedBody)
		assert.NoError(err)
		verifyStatusCode(t, resp, http.StatusOK)

		var updateRuleResponse models.UpdateSignRuleResponse
		err = json.Unmarshal(resp.Body.Bytes(), &updateRuleResponse)
		assert.NoError(err)
		assert.NotEmpty(updateRuleResponse.ID)

		// Check rule is updated
		rule := cfg.GetSignRule(account.Nickname, updateRuleResponse.ID)
		assert.NotNil(rule)
		assert.Equal(rule.Name, updatedRuleName)
		assert.Equal(rule.Contract, contract)
		assert.Equal(rule.RuleType, ruleType)
		assert.Equal(rule.Enabled, enabled)

		// Check that privateKey is still cached
		pkey, err := cache.PrivateKeyFromCache(account)
		assert.NoError(err)
		assert.NotNil(pkey)
	})

	t.Run("Delete sign rule", func(t *testing.T) {
		contract := "AS1hCJXjndR4c9vekLWsXGnrdigp4AaZ7uYG3UKFzzKnWVsrNLPJ"
		ruleType := config.RuleTypeAutoSign
		ruleName := "Test Rule"
		enabled := true
		body := fmt.Sprintf(`{
			"name": "%s",
			"contract": "%s",
			"ruleType": "%s",
			"enabled": %t
		}`, ruleName, contract, ruleType, enabled)

		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.StringPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Message:     password,
			}

			res <- (<-resChan)
		}(testResult)

		// Add a sign rule first
		resp, err := handleHTTPRequest(addhandler, "POST", fmt.Sprintf("/api/accounts/%s/signrules", nickname), body)
		assert.NoError(err)

		verifyStatusCode(t, resp, http.StatusOK)

		var addRuleResponse models.AddSignRuleResponse
		err = json.Unmarshal(resp.Body.Bytes(), &addRuleResponse)
		assert.NoError(err)
		assert.NotEmpty(addRuleResponse.ID)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.StringPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Message:     password,
			}

			res <- (<-resChan)
		}(testResult)

		// Now delete the sign rule
		resp, err = handleHTTPRequest(deletehandler, "DELETE", fmt.Sprintf("/api/accounts/%s/signrules/%s", nickname, addRuleResponse.ID), "")
		assert.NoError(err)
		verifyStatusCode(t, resp, http.StatusOK)

		// Check rule is deleted
		rulePtr := cfg.GetEnabledRuleForContract(account.Nickname, &contract)
		assert.Nil(rulePtr)

		hasRule := cfg.HasEnabledRule(account.Nickname)
		assert.True(hasRule)

		// Check that privateKey is still cached
		pkey, err := cache.PrivateKeyFromCache(account)
		assert.NoError(err)
		assert.NotNil(pkey)
	})
}
