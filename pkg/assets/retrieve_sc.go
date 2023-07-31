package assets

import (
	"fmt"
	"sync"

	"github.com/go-openapi/swag"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	network "github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station/pkg/convert"
	"github.com/massalabs/station/pkg/node"
)

const (
	NAME_KEY           = "NAME"
	SYMBOL_KEY         = "SYMBOL"
	DECIMALS_KEY       = "DECIMALS"
	BALANCE_KEY_PREFIX = "BALANCE"
)

// Function to convert an address to a storage key using the balance key prefix
func balanceKey(address string) []byte {
	return convert.ToBytes(BALANCE_KEY_PREFIX + address)
}

// Balance retrieves the balance of a user for a given asset contract address by making a smart contract call.
func (s *AssetsStore) Balance(assetContractAddress, userAddress string) (string, error) {
	client, err := network.NewMassaClient()
	if err != nil {
		return "", fmt.Errorf("failed to create Massa client: %w", err)
	}

	balanceData, err := node.DatastoreEntry(client, assetContractAddress, balanceKey(userAddress))
	if err != nil {
		return "", fmt.Errorf("failed to fetch user balance: %w", err)
	}

	balanceValue, err := convert.BytesToU256(balanceData.CandidateValue)
	if err != nil {
		return "", fmt.Errorf("failed to parse user balance: %w", err)
	}

	return fmt.Sprint(balanceValue.String()), nil
}

// AssetInfo retrieves the asset information for a given contract address by making a smart contract call.
func (s *AssetsStore) AssetInfo(contractAddress string) (*models.AssetInfo, error) {
	client, err := network.NewMassaClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Massa client: %w", err)
	}

	// Create WaitGroup to wait for all Goroutines to finish
	var wg sync.WaitGroup
	wg.Add(3)

	// Channels to receive the results of each Goroutine
	nameCh := make(chan string, 1)
	symbolCh := make(chan string, 1)
	decimalsCh := make(chan uint8, 1)
	errCh := make(chan error, 3) // Channel to collect errors

	// Concurrently fetch asset name
	go func() {
		defer wg.Done()
		nameData, err := node.DatastoreEntry(client, contractAddress, convert.ToBytes(NAME_KEY))
		if err != nil {
			errCh <- fmt.Errorf("failed to fetch asset name: %w", err)
			return
		}
		nameCh <- string(nameData.CandidateValue)
	}()

	// Concurrently fetch asset symbol
	go func() {
		defer wg.Done()
		symbolData, err := node.DatastoreEntry(client, contractAddress, convert.ToBytes(SYMBOL_KEY))
		if err != nil {
			errCh <- fmt.Errorf("failed to fetch asset symbol: %w", err)
			return
		}
		symbolCh <- string(symbolData.CandidateValue)
	}()

	// Concurrently fetch asset decimals
	go func() {
		defer wg.Done()
		decimals, err := node.DatastoreEntry(client, contractAddress, convert.ToBytes(DECIMALS_KEY))
		if err != nil {
			errCh <- fmt.Errorf("failed to fetch asset decimals: %w", err)
			return
		}
		decimalsCh <- uint8((decimals.CandidateValue)[0])
	}()

	// Wait for all Goroutines to finish
	wg.Wait()
	close(nameCh)
	close(symbolCh)
	close(decimalsCh)
	close(errCh) // Close the errCh channel after all Goroutines have finished

	// Collect any errors from the error channel
	for err := range errCh {
		return nil, err
	}

	// Read results from channels
	name := <-nameCh
	symbol := <-symbolCh
	decimals := <-decimalsCh

	// Create the AssetInfo struct with the retrieved information
	assetInfoFromSC := &models.AssetInfo{
		Name:     name,
		Symbol:   symbol,
		Decimals: swag.Int64(int64(decimals)),
	}

	return assetInfoFromSC, nil
}
