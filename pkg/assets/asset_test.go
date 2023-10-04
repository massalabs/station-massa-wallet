package assets

import (
	"os"
	"testing"

	"github.com/massalabs/station-massa-wallet/api/server/models"
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
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name()) // Clean up the temporary file after the test
	}()

	_, err = tmpFile.Write([]byte(dummyAssetsJSON))
	if err != nil {
		t.Fatalf("Error writing to temporary file: %v", err)
	}

	// Create a new instance of AssetsStore and load data from the testing file
	store, err := NewAssetsStore(tmpFile.Name())
	if err != nil {
		t.Fatalf("Error creating AssetsStore: %v", err)
	}

	// Validate the loaded data
	expectedAccountName := "dummyAccount"
	expectedContractAddress := "0x1234567890abcdef"
	asset, ok := store.Assets[expectedAccountName].ContractAssets[expectedContractAddress]
	if !ok {
		t.Fatalf("Expected asset not found in loaded data")
	}

	// Example: Validate the name of the loaded asset
	if asset.Name != "DummyToken" {
		t.Errorf("Unexpected asset name. Got %s, expected %s", asset.Name, "DummyToken")
	}
}

func TestAssetExists(t *testing.T) {
	// Create a temporary testing JSON file with dummy data
	tmpFile, err := os.CreateTemp("", "dummy_assets.json")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name()) // Clean up the temporary file after the test
	}()

	_, err = tmpFile.Write([]byte(dummyAssetsJSON))
	if err != nil {
		t.Fatalf("Error writing to temporary file: %v", err)
	}

	// Create a new instance of AssetsStore and load data from the testing file
	store, err := NewAssetsStore(tmpFile.Name())
	if err != nil {
		t.Fatalf("Error creating AssetsStore: %v", err)
	}

	// Test case 1: Check for an existing asset
	existingNickname := "dummyAccount"
	existingContractAddress := "0x1234567890abcdef"
	exists := store.AssetExists(existingNickname, existingContractAddress)
	if !exists {
		t.Errorf("Expected asset to exist, but AssetExists returned false")
	}

	// Test case 2: Check for a non-existing asset
	nonExistingNickname := "nonExistingAccount"
	nonExistingContractAddress := "0xabcdefabcdefabcdef"
	notExists := store.AssetExists(nonExistingNickname, nonExistingContractAddress)
	if notExists {
		t.Errorf("Expected asset not to exist, but AssetExists returned true")
	}
}

func TestAddAndDeleteAsset(t *testing.T) {
	// Create a temporary testing JSON file with dummy data
	tmpFile, err := os.CreateTemp("", "dummy_assets.json")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
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
	if err != nil {
		t.Fatalf("Error writing to temporary file: %v", err)
	}

	// Create a new instance of AssetsStore and load data from the testing file
	store, err := NewAssetsStore(tmpFile.Name())
	if err != nil {
		t.Fatalf("Error creating AssetsStore: %v", err)
	}

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
	if err != nil {
		t.Fatalf("Error adding asset: %v", err)
	}

	// Check if the added asset exists
	if !store.AssetExists(nickname, assetAddress) {
		t.Errorf("Added asset not found in reloaded store")
	}

	// Test case 2: Delete the added asset and check if it's removed from JSON
	err = store.DeleteAsset(nickname, assetAddress)
	if err != nil {
		t.Fatalf("Error deleting asset: %v", err)
	}

	// Reload the AssetsStore from the JSON file after deletion
	deletedStore, err := NewAssetsStore(tmpFile.Name())
	if err != nil {
		t.Fatalf("Error reloading AssetsStore after deletion: %v", err)
	}

	// Check if the deleted asset no longer exists
	if deletedStore.AssetExists(nickname, assetAddress) {
		t.Errorf("Deleted asset still found in reloaded store")
	}
}
