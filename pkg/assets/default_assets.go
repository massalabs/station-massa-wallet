package assets

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
)

func GetDefaultAssets() ([]models.AssetInfo, error) {
	defaultAssetsJSONPath, err := getDefaultAssetsJSONPath()
	if err != nil {
		return nil, err
	}

	defaultAssets, err := loadDefaultAssets(defaultAssetsJSONPath)
	if err != nil {
		return nil, err
	}

	return defaultAssets, nil
}

// loadDefaultAssets loads the default assets from the JSON file.
func loadDefaultAssets(path string) ([]models.AssetInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var defaultAssets []models.AssetInfo
	if err := json.NewDecoder(file).Decode(&defaultAssets); err != nil {
		return nil, err
	}

	return defaultAssets, nil
}

func InitDefaultAsset() error {
	// Get the path to the default assets JSON file
	defaultAssetsJSONPath, err := getDefaultAssetsJSONPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(defaultAssetsJSONPath); os.IsNotExist(err) {
		if err := createFileDefaultAssets(defaultAssetsJSONPath); err != nil {
			return err
		}
	}

	return nil
}

// getDefaultAssetsJSONPath returns the path to the default assets JSON file.
func getDefaultAssetsJSONPath() (string, error) {
	walletPath, err := wallet.Path()
	if err != nil {
		return "", err
	}

	return filepath.Join(walletPath, "assets_default.json"), nil
}

// createFileDefaultAssets creates the default assets JSON file with the default assets.
func createFileDefaultAssets(path string) error {
	if err := os.WriteFile(path, []byte(`[
	{
		"address": "AS12k8viVmqPtRuXzCm6rKXjLgpQWqbuMjc37YHhB452KSUUb9FgL",
		"name": "Sepolia USDC",
		"symbol": "USDC.s",
		"decimals": 6
	},
	{
		"address": "AS12LpYyAjYRJfYhyu7fkrS224gMdvFHVEeVWoeHZzMdhis7UZ3Eb",
		"name": "Sepolia tDAI",
		"symbol": "tDAI.s",
		"decimals": 18
	},
	{
		"address": "AS1gt69gqYD92dqPyE6DBRJ7KjpnQHqFzFs2YCkBcSnuxX5bGhBC",
		"name": "sepolia WETH",
		"symbol": "WETH.s",
		"decimals": 18
	},
	{
		"address": "AS12U4TZfNK7qoLyEERBBRDMu8nm5MKoRzPXDXans4v9wdATZedz9",
		"name": "Wrapped Massa",
		"symbol": "WMAS",
		"decimals": 9
	},
	{
		"address": "AS1hCJXjndR4c9vekLWsXGnrdigp4AaZ7uYG3UKFzzKnWVsrNLPJ",
		"name": "USD Coin",
		"symbol": "USDC.e",
		"decimals": 6
	},
	{
		"address": "AS1ZGF1upwp9kPRvDKLxFAKRebgg7b3RWDnhgV7VvdZkZsUL7Nuv",
		"name": "Dai Stablecoin",
		"symbol": "DAI.e",
		"decimals": 18
	},
	{
		"address": "AS124vf3YfAJCSCQVYKczzuWWpXrximFpbTmX4rheLs5uNSftiiRY",
		"name": "Wrapped Ether",
		"symbol": "WETH.e",
		"decimals": 18
	}
]`), permissionUrwGrOr); err != nil {
		return err
	}

	return nil
}
