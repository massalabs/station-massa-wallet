package wallet

import (
	"testing"
)

func Test_walletImport_Handle(t *testing.T) {
	api, _, err := MockAPI()
	if err != nil {
		panic(err)
	}

	testsImport := []struct {
		name       string
		body       string
		statusCode int
	}{
		{"passing", `{
		"address": "A12TFdPyw8Sg9qouzgTWwW5yo5PBDu5C3BWEGPjB9vRx9s3b42qv",
		"keyPair": 
			{
				"nonce": "86zrpLuzBXBtePQiC5b4d1",
				"privateKey": "HvAH6XuMNamRCuCuAaGsUKrCSjFibwyZ35aHZ4zBd5iNM5x2YM74vLUUH9KhAxDKGxWJ4V3YWNvGGiziGjC4yA1J72NKmcVMitHZM23eW44FAHay4iA",
				"publicKey": "ub4aTM9RSBydGJCbkxe8v7GqWpZNNXuh7uGQgthBpaWhocvA1",
				"salt": "4B28WQKc6jaYN6ymx6xoX8GzwHqF"
			},
		"nickname": "imported" }`, 204},
		{"fail_empty_fields", `{}`, 422}}
	for _, tt := range testsImport {
		t.Run(tt.name, func(t *testing.T) {

			resp, err := processHTTPRequest(api, "PUT", "/rest/wallet", tt.body)
			if err != nil {
				t.Fatalf("while serving HTTP request: %s", err)
			}

			if resp.Result().StatusCode != tt.statusCode {
				t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, tt.statusCode)
			}

			if tt.statusCode != 204 {
				return
			}

			err = cleanupTestData([]string{"imported"})
			if err != nil {
				t.Fatalf("while cleaning up TestData: %s", err)
			}
		})
	}
}
