package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/awnumar/memguard"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/cache"
	"github.com/massalabs/station-massa-wallet/pkg/config"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/stretchr/testify/assert"
)

const (
	callSCString = "AKT4CASAzuTNAqCNBgEAXBwUw39NBQYix8Ovph0TUiJuDDEnlFYUPgsbeMbrA4cLZm9yd2FyZEJ1cm7FAQDgfY7fLW7qpwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACoAAAAweGZDRERBRTI1MTAwNjIxYTViQzg4MTlkQzlEMzg0MjUzNEQ3QmY0NzYAAAAANQAAAEFTMTJUUm9TY01kd0xLOFlwdDZOQkFwcHl6Q0Z3N1FlRzVlM3hGdnhwQ0FuQW5ZTGZ1TVVUKgAAADB4NTM4NDRGOTU3N0MyMzM0ZTU0MUFlYzdEZjcxNzRFQ2U1ZEYxZkNmMKc2qgAAAAAA"
	contract     = "AS1hZpUH6TPiRxHtTKqAfXDmZ7Afa7UfS4rtYN7NxVwAaSAphCET"
)

func signTransaction(t *testing.T, api *operations.MassaWalletAPI, nickname string, body string, headers ...map[string]string) *httptest.ResponseRecorder {
	handler, exist := api.HandlerFor("post", "/api/accounts/{nickname}/sign")
	assert.True(t, exist)

	// Use provided headers if any, otherwise use empty map
	headerMap := make(map[string]string)
	if len(headers) > 0 {
		headerMap = headers[0]
	}

	resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/sign", nickname), body, headerMap)
	assert.NoError(t, err)

	return resp
}

func signMessage(t *testing.T, api *operations.MassaWalletAPI, nickname string, body string) *httptest.ResponseRecorder {
	handler, exist := api.HandlerFor("post", "/api/accounts/{nickname}/signMessage")
	assert.True(t, exist)

	resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/signMessage", nickname), body)
	assert.NoError(t, err)

	return resp
}

func verifySignResponse(t *testing.T, resp *httptest.ResponseRecorder) {
	var signResponse models.SignResponse
	err := json.Unmarshal(resp.Body.Bytes(), &signResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, signResponse.Signature)
	assert.NotEmpty(t, signResponse.Operation)
	assert.NotEmpty(t, signResponse.PublicKey)
}

func Test_walletSign_Handle(t *testing.T) {
	api, resChan, err := MockAPI()
	assert.NoError(t, err)

	transactionData := fmt.Sprintf(`{"chainId": `+strconv.FormatUint(ChainIDUnitTests, 10)+`, "operation":"%s"}`, callSCString)
	nickname := "walletToDelete"
	password := "zePassword"
	account := createAccount(password, nickname, t, prompterAppMock)

	t.Run("invalid nickname", func(t *testing.T) {
		resp := signTransaction(t, api, "Johnny", transactionData)
		verifyStatusCode(t, resp, http.StatusNotFound)
	})

	t.Run("sign transaction OK", func(t *testing.T) {
		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Password:    password,
				Fees:        "14400",
			}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp := signTransaction(t, api, nickname, transactionData)
		verifyStatusCode(t, resp, http.StatusOK)

		result := <-testResult

		checkResultChannel(t, result, true, "")
		verifySignResponse(t, resp)
	})

	t.Run("sign a plain text message", func(t *testing.T) {
		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Password:    password,
				Fees:        "1000",
			}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		message := "Test"
		requestBody := fmt.Sprintf(`{"message":"%s"}`, message)

		resp := signMessage(t, api, nickname, requestBody)
		verifyStatusCode(t, resp, http.StatusOK)

		result := <-testResult

		checkResultChannel(t, result, true, "")

		var signResponse models.SignResponse
		err := json.Unmarshal(resp.Body.Bytes(), &signResponse)
		assert.NoError(t, err)
		assert.NotEmpty(t, signResponse.Signature)
		assert.NotEmpty(t, signResponse.PublicKey)
	})

	// The handler will not return until a the good password is sent or the action is canceled
	t.Run("invalid password try, then valid password", func(t *testing.T) {
		testResult := make(chan walletapp.EventData)

		//nolint:staticcheck
		go func(res chan walletapp.EventData) {
			// Send wrong password to prompter app and wait for result
			prompterAppMock.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Password:    "this is not the password",
				Fees:        "1000",
			}
			// forward test result to test goroutine
			failRes := <-resChan

			checkResultChannel(t, failRes, false, utils.WrongPassword)

			// Send password to prompter app to unlock the handler
			prompterAppMock.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Password:    password,
				Fees:        "1000",
			}

			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp := signTransaction(t, api, nickname, transactionData)
		verifyStatusCode(t, resp, http.StatusOK)

		result := <-testResult

		checkResultChannel(t, result, true, "")
		verifySignResponse(t, resp)
	})

	t.Run("invalid password try, then action canceled by user", func(t *testing.T) {
		//nolint:staticcheck
		go func() {
			// Send wrong password to prompter app and wait for result
			prompterAppMock.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Password:    "this is not the password",
				Fees:        "1000",
			}
			// forward test result to test goroutine
			failRes := <-resChan

			checkResultChannel(t, failRes, false, utils.WrongPassword)

			// Send cancel to prompter app to unlock the handler
			prompterAppMock.App().CtrlChan <- walletapp.Cancel
		}()

		resp := signTransaction(t, api, nickname, transactionData)
		verifyStatusCode(t, resp, http.StatusUnauthorized)
	})

	t.Run("Auto Sign", func(t *testing.T) {
		cfg := config.Get()
		authorizedOrigin := "http://massa.network"
		
		// Add AutoSign rule
		ruleId, err := cfg.AddSignRule(nickname, config.SignRule{
			Name:             "test",
			Contract:         contract,
			RuleType:         config.RuleTypeAutoSign,
			Enabled:          true,
			AuthorizedOrigin: &authorizedOrigin,
		})
		assert.NoError(t, err)
		assert.True(t, cfg.HasEnabledRule(nickname))

		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Password:    password,
				Fees:        "14400",
			}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		// First sign with password prompt
		headers := map[string]string{
			originHeader: authorizedOrigin,
		}
		resp := signTransaction(t, api, nickname, transactionData, headers)
		verifyStatusCode(t, resp, http.StatusOK)

		result := <-testResult
		checkResultChannel(t, result, true, "")

		// check that privateKey is cached
		pkey, err := cache.PrivateKeyFromCache(account)
		assert.NoError(t, err)
		assert.NotNil(t, pkey)

		// sign again, should not prompt for password
		resp = signTransaction(t, api, nickname, transactionData, headers)
		verifyStatusCode(t, resp, http.StatusOK)
		verifySignResponse(t, resp)

		// Clean up
		err = cfg.DeleteSignRule(nickname, ruleId)
		assert.NoError(t, err)
		// Clear the private key from cache
		cacheInstance := cache.Init()
		address, err := account.Address.String()
		assert.NoError(t, err)
		cacheKey := "pkey" + address
		cacheInstance.Remove(cache.KeyHash([]byte(cacheKey)))
	})

	t.Run("Auto Sign failed if origin in header not the same", func(t *testing.T) {
		cfg := config.Get()
		authorizedOrigin := "http://massa.network"
		
		// Add AutoSign rule
		ruleId, err := cfg.AddSignRule(nickname, config.SignRule{
			Name:             "test",
			Contract:         contract,
			RuleType:         config.RuleTypeAutoSign,
			Enabled:          true,
			AuthorizedOrigin: &authorizedOrigin,
		})
		assert.NoError(t, err)
		assert.True(t, cfg.HasEnabledRule(nickname))

		// Cache the private key to test the AutoSign flow
		passwordBuffer := memguard.NewBufferFromBytes([]byte(password))
		defer passwordBuffer.Destroy()
		err = cache.CachePrivateKeyFromPassword(account, passwordBuffer)
		assert.NoError(t, err)

		// Verify private key is cached
		pkey, err := cache.PrivateKeyFromCache(account)
		assert.NoError(t, err)
		assert.NotNil(t, pkey)

		// Test with different origin
		headers := map[string]string{
			originHeader: "http://different-origin.com",
		}

		resp := signTransaction(t, api, nickname, transactionData, headers)
		verifyStatusCode(t, resp, http.StatusUnauthorized)

		var errorResponse models.Error
		err = json.Unmarshal(resp.Body.Bytes(), &errorResponse)
		assert.NoError(t, err)
		assert.Equal(t, "Unauthorized origin", errorResponse.Message)
		assert.Equal(t, "errorUnauthorizedOrigin", errorResponse.Code)

		// Clean up
		err = cfg.DeleteSignRule(nickname, ruleId)
		assert.NoError(t, err)
		// Clear the private key from cache
		cacheInstance := cache.Init()
		address, err := account.Address.String()
		assert.NoError(t, err)
		cacheKey := "pkey" + address
		cacheInstance.Remove(cache.KeyHash([]byte(cacheKey)))
	})

	t.Run("Auto Sign with enabled rule", func(t *testing.T) {
		// Clean cache
		testCache.Purge()

		authorizedOrigin := "http://massa.network"

		cfg := config.Get()
		ruleId, err := cfg.AddSignRule(nickname, config.SignRule{
			Name:             "test",
			Contract:         "AS1ZGF1upwp9kPRvDKLxFAKRebgg7b3RWDnhgV7VvdZkZsUL7Nuv",
			RuleType:         config.RuleTypeAutoSign,
			Enabled:          true,
			AuthorizedOrigin: &authorizedOrigin,
		})

		assert.NoError(t, err)

		assert.True(t, cfg.HasEnabledRule(nickname))

		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Password:    password,
				Fees:        "14400",
			}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp := signTransaction(t, api, nickname, transactionData)
		verifyStatusCode(t, resp, http.StatusOK)

		result := <-testResult

		checkResultChannel(t, result, true, "")

		// check that privateKey is cached
		pkey, err := cache.PrivateKeyFromCache(account)
		assert.NoError(t, err)
		assert.NotNil(t, pkey)

		// sign again, should prompt for password
		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Password:    password,
				Fees:        "14400",
			}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp = signTransaction(t, api, nickname, transactionData)
		verifyStatusCode(t, resp, http.StatusOK)

		result = <-testResult

		checkResultChannel(t, result, true, "")

		verifySignResponse(t, resp)

		err = cfg.DeleteSignRule(nickname, ruleId)
		assert.NoError(t, err)
	})

	t.Run("Auto Sign with enabled rule and authorized origin fail if origin is not the same", func(t *testing.T) {
		cfg := config.Get()
		authorizedOrigin := "http://massa.network"
		ruleId, err := cfg.AddSignRule(nickname, config.SignRule{
			Name:             "test",
			Contract:         contract,
			RuleType:         config.RuleTypeAutoSign,
			Enabled:          true,
			AuthorizedOrigin: &authorizedOrigin,
		})

		assert.NoError(t, err)

		assert.True(t, cfg.HasEnabledRule(nickname))

		headers := map[string]string{
			originHeader: "http://other-origin:3000",
		}

		resp := signTransaction(t, api, nickname, transactionData, headers)
		verifyStatusCode(t, resp, http.StatusUnauthorized)

		err = cfg.DeleteSignRule(nickname, ruleId)
		assert.NoError(t, err)
	})

	t.Run("Disable Password Prompt for contract", func(t *testing.T) {
		// Clean cache
		testCache.Purge()

		cfg := config.Get()
		ruleId, err := cfg.AddSignRule(nickname, config.SignRule{
			Name:     "test",
			Contract: contract,
			RuleType: config.RuleTypeDisablePasswordPrompt,
			Enabled:  true,
		})

		assert.NoError(t, err)

		assert.True(t, cfg.HasEnabledRule(nickname))

		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Password:    password,
				Fees:        "14400",
			}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp := signTransaction(t, api, nickname, transactionData)
		verifyStatusCode(t, resp, http.StatusOK)

		result := <-testResult

		checkResultChannel(t, result, true, "")

		// check that privateKey is cached
		pkey, err := cache.PrivateKeyFromCache(account)
		assert.NoError(t, err)
		assert.NotNil(t, pkey)

		// sign again, should prompt for user confirmation, no password needed
		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Fees:        "12345",
			}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp = signTransaction(t, api, nickname, transactionData)
		verifyStatusCode(t, resp, http.StatusOK)

		result = <-testResult

		checkResultChannel(t, result, true, "")

		verifySignResponse(t, resp)

		err = cfg.DeleteSignRule(nickname, ruleId)
		assert.NoError(t, err)
	})

	t.Run("Disable Password Prompt enabled for another contract", func(t *testing.T) {
		// Clean cache
		testCache.Purge()

		cfg := config.Get()
		ruleId, err := cfg.AddSignRule(nickname, config.SignRule{
			Name:     "test",
			Contract: "AS12LKs9txoSSy8JgFJgV96m8k5z9pgzjYMYSshwN67mFVuj3bdUV",
			RuleType: config.RuleTypeDisablePasswordPrompt,
			Enabled:  true,
		})

		assert.NoError(t, err)

		assert.True(t, cfg.HasEnabledRule(nickname))

		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Password:    password,
				Fees:        "14400",
			}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp := signTransaction(t, api, nickname, transactionData)
		verifyStatusCode(t, resp, http.StatusOK)

		result := <-testResult

		checkResultChannel(t, result, true, "")

		// check that privateKey is cached
		pkey, err := cache.PrivateKeyFromCache(account)
		assert.NoError(t, err)
		assert.NotNil(t, pkey)

		// sign again, should prompt for user confirmation, password is needed
		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Password:    password,
				Fees:        "14400",
			}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp = signTransaction(t, api, nickname, transactionData)
		verifyStatusCode(t, resp, http.StatusOK)

		result = <-testResult

		checkResultChannel(t, result, true, "")

		verifySignResponse(t, resp)

		err = cfg.DeleteSignRule(nickname, ruleId)
		assert.NoError(t, err)
	})

	t.Run("Disable Password Prompt enabled for all contract.", func(t *testing.T) {
		// Clean cache
		testCache.Purge()

		cfg := config.Get()
		ruleId, err := cfg.AddSignRule(nickname, config.SignRule{
			Name:     "test",
			Contract: "*",
			RuleType: config.RuleTypeDisablePasswordPrompt,
			Enabled:  true,
		})

		assert.NoError(t, err)

		assert.True(t, cfg.HasEnabledRule(nickname))

		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Password:    password,
				Fees:        "14400",
			}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp := signTransaction(t, api, nickname, transactionData)
		verifyStatusCode(t, resp, http.StatusOK)

		result := <-testResult

		checkResultChannel(t, result, true, "")

		// check that privateKey is cached
		pkey, err := cache.PrivateKeyFromCache(account)
		assert.NoError(t, err)
		assert.NotNil(t, pkey)

		// sign again, should prompt for user confirmation, password is not needed
		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Fees:        "14400",
			}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp = signTransaction(t, api, nickname, transactionData)
		verifyStatusCode(t, resp, http.StatusOK)

		result = <-testResult

		checkResultChannel(t, result, true, "")

		verifySignResponse(t, resp)

		err = cfg.DeleteSignRule(nickname, ruleId)
		assert.NoError(t, err)
	})

	t.Run("Auto Sign rule type origin validation", func(t *testing.T) {
		cfg := config.Get()
		authorizedOrigin := "http://massa.network"
		
		// Add AutoSign rule
		ruleId, err := cfg.AddSignRule(nickname, config.SignRule{
			Name:             "test",
			Contract:         contract,
			RuleType:         config.RuleTypeAutoSign,
			Enabled:          true,
			AuthorizedOrigin: &authorizedOrigin,
		})
		assert.NoError(t, err)
		assert.True(t, cfg.HasEnabledRule(nickname))

		// Cache the private key to test the AutoSign flow
		passwordBuffer := memguard.NewBufferFromBytes([]byte(password))
		defer passwordBuffer.Destroy()
		err = cache.CachePrivateKeyFromPassword(account, passwordBuffer)
		assert.NoError(t, err)

		// Verify private key is cached
		pkey, err := cache.PrivateKeyFromCache(account)
		assert.NoError(t, err)
		assert.NotNil(t, pkey)

		testCases := []struct {
			name           string
			origin         string
			expectedStatus int
			expectedError  string
			description    string
		}{
			{
				name:           "Exact match should succeed",
				origin:         "http://massa.network",
				expectedStatus: http.StatusOK,
				expectedError:  "",
				description:    "Origin exactly matches authorized origin",
			},
			{
				name:           "Different origin should fail",
				origin:         "http://other-origin:3000",
				expectedStatus: http.StatusUnauthorized,
				expectedError:  "Unauthorized origin",
				description:    "Origin doesn't match authorized origin",
			},
			{
				name:           "Missing origin should fail",
				origin:         "",
				expectedStatus: http.StatusUnauthorized,
				expectedError:  "Unauthorized origin",
				description:    "No origin provided",
			},
			{
				name:           "Different protocol should fail",
				origin:         "https://massa.network",
				expectedStatus: http.StatusUnauthorized,
				expectedError:  "Unauthorized origin",
				description:    "Different protocol (https vs http)",
			},
			{
				name:           "Different subdomain should fail",
				origin:         "http://sub.massa.network",
				expectedStatus: http.StatusUnauthorized,
				expectedError:  "Unauthorized origin",
				description:    "Different subdomain",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				headers := map[string]string{}
				if tc.origin != "" {
					headers[originHeader] = tc.origin
				}

				resp := signTransaction(t, api, nickname, transactionData, headers)
				verifyStatusCode(t, resp, tc.expectedStatus)

				if tc.expectedStatus == http.StatusUnauthorized {
					var errorResponse models.Error
					err := json.Unmarshal(resp.Body.Bytes(), &errorResponse)
					assert.NoError(t, err)
					assert.Equal(t, tc.expectedError, errorResponse.Message)
					assert.Equal(t, "errorUnauthorizedOrigin", errorResponse.Code)
				} else {
					// For successful cases, verify the transaction was signed
					var signResponse models.SignResponse
					err := json.Unmarshal(resp.Body.Bytes(), &signResponse)
					assert.NoError(t, err)
					assert.NotEmpty(t, signResponse.Signature)
					assert.NotEmpty(t, signResponse.Operation)
					assert.NotEmpty(t, signResponse.PublicKey)
				}
			})
		}

		// Clean up
		err = cfg.DeleteSignRule(nickname, ruleId)
		assert.NoError(t, err)
		// Clear the private key from cache
		cacheInstance := cache.Init()
		address, err := account.Address.String()
		assert.NoError(t, err)
		cacheKey := "pkey" + address
		cacheInstance.Remove(cache.KeyHash([]byte(cacheKey)))
	})
}
