package assets

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const defaultAssetsFilename = "assets_default.json"

type DefaultAssetInfo struct {
	Address    string `json:"address"`
	Name       string `json:"name"`
	Symbol     string `json:"symbol"`
	Decimals   int64  `json:"decimals"`
	MEXCSymbol string `json:"MEXCSymbol"`
}

func (s *AssetsStore) Default() ([]DefaultAssetInfo, error) {
	defaultAssetsJSONPath, err := getDefaultJSONPath(s.assetsJSONDir)
	if err != nil {
		return nil, err
	}

	defaultAssets, err := s.loadDefaultAssets(defaultAssetsJSONPath)
	if err != nil {
		return nil, err
	}

	return defaultAssets, nil
}

// loadDefaultAssets loads the default assets from the JSON file.
func (s *AssetsStore) loadDefaultAssets(path string) ([]DefaultAssetInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var defaultAssets []DefaultAssetInfo
	if err := json.NewDecoder(file).Decode(&defaultAssets); err != nil {
		return nil, err
	}

	return defaultAssets, nil
}

func (s *AssetsStore) InitDefault() error {
	// Get the path to the default assets JSON file
	defaultAssetsJSONPath, err := getDefaultJSONPath(s.assetsJSONDir)
	if err != nil {
		return err
	}

	// if the file does not exist, create the default assets JSON file
	_, err = os.Stat(defaultAssetsJSONPath)
	if os.IsNotExist(err) {
		// Create the default assets JSON file
		return s.createFileDefault(defaultAssetsJSONPath)
	}

	// if the file exists, read the content and compare it with the default assets
	if err == nil {
		// read the content of the default assets JSON file
		content, err := os.ReadFile(defaultAssetsJSONPath)
		if err != nil {
			return err
		}

		// if the content is different, overwrite the default assets JSON file
		if string(content) != assetsJSON {
			if err := s.createFileDefault(defaultAssetsJSONPath); err != nil {
				return err
			}
		}
	}

	return err
}

// getDefaultJSONPath returns the path to the default assets JSON file.
func getDefaultJSONPath(assetsJSONDir string) (string, error) {
	return filepath.Join(assetsJSONDir, defaultAssetsFilename), nil
}

// createFileDefault creates the default assets JSON file with the default assets.
func (s *AssetsStore) createFileDefault(path string) error {
	if err := os.WriteFile(path, []byte(assetsJSON), permissionUrwGrOr); err != nil {
		return err
	}

	return nil
}

const assetsJSON = `[
	{
		"address": "AS12k8viVmqPtRuXzCm6rKXjLgpQWqbuMjc37YHhB452KSUUb9FgL",
		"name": "Sepolia USDC",
		"symbol": "USDC.s",
		"decimals": 6,
		"MEXCSymbol": "USD"
	},
	{
		"address": "AS12LpYyAjYRJfYhyu7fkrS224gMdvFHVEeVWoeHZzMdhis7UZ3Eb",
		"name": "Sepolia tDAI",
		"symbol": "tDAI.s",
		"decimals": 18,
		"MEXCSymbol": "USD"
	},
	{
		"address": "AS1gt69gqYD92dqPyE6DBRJ7KjpnQHqFzFs2YCkBcSnuxX5bGhBC",
		"name": "sepolia WETH",
		"symbol": "WETH.s",
		"decimals": 18,
		"MEXCSymbol": "ETHUSDT"
	},
	{
		"address": "AS12U4TZfNK7qoLyEERBBRDMu8nm5MKoRzPXDXans4v9wdATZedz9",
		"name": "Wrapped Massa",
		"symbol": "WMAS",
		"decimals": 9,
		"MEXCSymbol": "MASUSDT"
	},
	{
		"address": "AS1hCJXjndR4c9vekLWsXGnrdigp4AaZ7uYG3UKFzzKnWVsrNLPJ",
		"name": "USD Coin",
		"symbol": "USDC.e",
		"decimals": 6,
		"MEXCSymbol": "USD"
	},
	{
		"address": "AS1ZGF1upwp9kPRvDKLxFAKRebgg7b3RWDnhgV7VvdZkZsUL7Nuv",
		"name": "Dai Stablecoin",
		"symbol": "DAI.e",
		"decimals": 18,
		"MEXCSymbol": "USD"
	},
	{
		"address": "AS124vf3YfAJCSCQVYKczzuWWpXrximFpbTmX4rheLs5uNSftiiRY",
		"name": "Wrapped Ether",
		"symbol": "WETH.e",
		"decimals": 18,
		"MEXCSymbol": "ETHUSDT"
	},
	{
		"address": "",
		"name": "Wrapped Ether",
		"symbol": "WETH.b",
		"decimals": 18,
		"MEXCSymbol": "ETHUSDT"
	},
	{
		"address": "",
		"name": "Wrapped Binance USD",
		"symbol": "USDT.b",
		"decimals": 18,
		"MEXCSymbol": "USD"
	},
	{
		"address": "AS12RmCXTA9NZaTBUBnRJuH66AGNmtEfEoqXKxLdmrTybS6GFJPFs",
		"name": "Wrapped Ether",
		"symbol": "WETH.bt",
		"decimals": 18,
		"MEXCSymbol": "ETHUSDT"
	},
	{
		"address": "AS12ix1Qfpue7BB8q6mWVtjNdNE9UV3x4MaUo7WhdUubov8sJ3CuP",
		"name": "Wrapped Binance USD",
		"symbol": "USDT.bt",
		"decimals": 18,
		"MEXCSymbol": "USD"
	}
]`
