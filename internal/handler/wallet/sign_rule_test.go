package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/massalabs/station-massa-wallet/api/server/models"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	assertLib "github.com/stretchr/testify/assert"
)

func Test_signrule_Handlers(t *testing.T) {
	assert := assertLib.New(t)

	api, resChan, err := MockAPI()
	assert.NoError(err)

	addhandler, exist := api.HandlerFor("post", "/api/accounts/{nickname}/signrules")
	assert.True(exist)

	nickname := "zeWahLetName"
	password := "zePassword"
	account := createAccount(password, nickname, t, prompterAppMock)

	t.Run("Add sign rule", func(t *testing.T) {
		body := `{
			"name": "Test Rule",
			"contract": "Test Contract",
			"ruleType": "auto_sign",
			"enabled": true,
			"description": "Test Description"
		}`

		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.StringPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Message:     password,
			}

			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp, err := handleHTTPRequest(addhandler, "POST", fmt.Sprintf("/api/accounts/%s/signrules", nickname), body)
		assert.NoError(err)
		fmt.Println(resp)

		verifyStatusCode(t, resp, http.StatusOK)

		result := <-testResult

		checkResultChannel(t, result, true, "")

		var addRuleResponse models.AddSignRuleResponse
		err = json.Unmarshal(resp.Body.Bytes(), &addRuleResponse)
		assert.NoError(err)
		assert.NotEmpty(addRuleResponse.ID)

		// check that privateKey is cached
		pkey, err := privateKeyFromCache(testCache, account)
		assert.NoError(err)
		assert.NotNil(pkey)
	})
}

// func TestDeleteSignRule(t *testing.T) {
// 	api, _, _, resultChannel, err := MockAPI()
// 	assert.NoError(t, err)

// 	body := `{
// 		"nickname": "Test Nickname",
// 		"ruleID": "Test Rule ID"
// 	}`

// 	resp, err := processHTTPRequest(api, "DELETE", "/sign-rule", body)
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusOK, resp.Code)

// 	result := <-resultChannel
// 	checkResultChannel(t, result, true, "Sign rule deleted successfully")
// }

// func TestUpdateSignRule(t *testing.T) {
// 	api, _, _, resultChannel, err := MockAPI()
// 	assert.NoError(t, err)

// 	body := `{
// 		"nickname": "Test Nickname",
// 		"ruleID": "Test Rule ID",
// 		"name": "Updated Rule",
// 		"contract": "Updated Contract",
// 		"ruleType": "Updated Type",
// 		"enabled": false
// 	}`

// 	resp, err := processHTTPRequest(api, "PUT", "/sign-rule", body)
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusOK, resp.Code)

// 	result := <-resultChannel
// 	checkResultChannel(t, result, true, "Sign rule updated successfully")
// }
