package assets

import (
	"os"
	"testing"

	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/stretchr/testify/assert"
)

const dummyAssetsJSON = `
{
	"accounts": {
		"dummyAccount": {
			"assets": [
				{
					"contractAddress": "0x1234567890abcdef",
					"name":            "DummyToken",
					"symbol":          "DT",
					"decimals":        18
				},
				{
					"contractAddress": "0xabcdef1234567890",
					"name":            "AnotherDummyToken",
					"symbol":          "ADT",
					"decimals":        6
				}
			]
		}
	}
}
`

func TestLoadAccountsStore(t *testing.T) {
	// Create a temporary testing JSON file with dummy data
	tmpFile, err := os.CreateTemp("", "dummy_assets.json")
	assert.NoError(t, err)
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name()) // Clean up the temporary file after the test
	}()

	_, err = tmpFile.Write([]byte(dummyAssetsJSON))
	assert.NoError(t, err)

	// Create a new instance of AssetsStore and load data from the testing file
	store, err := NewAssetsStore(tmpFile.Name())
	assert.NoError(t, err)

	// Validate the loaded data
	expectedAccountName := "dummyAccount"
	expectedContractAddress := "0x1234567890abcdef"
	asset, ok := store.Assets[expectedAccountName].ContractAssets[expectedContractAddress]
	assert.True(t, ok)

	assert.Equal(t, asset.Name, "DummyToken")
}

func TestAssetExists(t *testing.T) {
	// Create a temporary testing JSON file with dummy data
	tmpFile, err := os.CreateTemp("", "dummy_assets.json")
	assert.NoError(t, err)

	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name()) // Clean up the temporary file after the test
	}()

	_, err = tmpFile.Write([]byte(dummyAssetsJSON))
	assert.NoError(t, err)

	// Create a new instance of AssetsStore and load data from the testing file
	store, err := NewAssetsStore(tmpFile.Name())
	assert.NoError(t, err)

	// Test case 1: Check for an existing asset
	existingNickname := "dummyAccount"
	existingContractAddress := "0x1234567890abcdef"
	exists := store.AssetExists(existingNickname, existingContractAddress)

	assert.True(t, exists, "Expected asset to exist ")

	// Test case 2: Check for a non-existing asset
	nonExistingNickname := "nonExistingAccount"
	nonExistingContractAddress := "0xabcdefabcdefabcdef"
	notExists := store.AssetExists(nonExistingNickname, nonExistingContractAddress)

	assert.False(t, notExists, "Expected asset to not exist ")
}

func TestAddAndDeleteAsset(t *testing.T) {
	// Create a temporary testing JSON file with dummy data
	tmpFile, err := os.CreateTemp("", "dummy_assets.json")
	assert.NoError(t, err)

	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name()) // Clean up the temporary file after the test
	}()

	// Initial dummy JSON data with at least one entry
	initialDummyJSON := `
{
	"accounts": {
		"dummyAccount": {
			"assets": [
				{
					"contractAddress": "0xffffffffffffff",
					"name":            "DummyToken",
					"symbol":          "DT",
					"decimals":        18
				}
			]
		}
	}
}
`
	_, err = tmpFile.Write([]byte(initialDummyJSON))
	assert.NoError(t, err)

	// Create a new instance of AssetsStore and load data from the testing file
	store, err := NewAssetsStore(tmpFile.Name())
	assert.NoError(t, err)

	// Test case 1: Add an asset and check if it's saved to JSON
	nickname := "dummyAccount"
	assetAddress := "0x1234567890abcdef"
	assetInfo := models.AssetInfo{
		Address:  assetAddress,
		Name:     "TestToken",
		Symbol:   "TT",
		Decimals: new(int64),
	}
	*assetInfo.Decimals = 18

	// Add the asset
	err = store.AddAsset(nickname, assetAddress, assetInfo)
	assert.NoError(t, err)

	// Check if the added asset exists
	assert.True(t, store.AssetExists(nickname, assetAddress), "Added asset not found in the store")

	// Test case 2: Delete the added asset and check if it's removed from JSON
	err = store.DeleteAsset(nickname, assetAddress)
	assert.NoError(t, err)

	// Check if the deleted asset no longer exists
	assert.False(t, store.AssetExists(nickname, assetAddress), "Deleted asset still found in the store")
}
