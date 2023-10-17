package assets

import (
	"fmt"
	"sync"

	"github.com/go-openapi/swag"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	network "github.com/massalabs/station-massa-wallet/pkg/network"
)

// AssetInfo retrieves the asset information for a given contract address by making a smart contract call.
func AssetInfo(contractAddress string, massaClient network.NodeFetcherInterface) (*models.AssetInfo, error) {
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

		nameData, err := massaClient.DatastoreAssetName(contractAddress)
		if err != nil {
			errCh <- fmt.Errorf("failed to fetch asset name: %w", err)
			return
		}
		nameCh <- nameData
	}()

	// Concurrently fetch asset symbol
	go func() {
		defer wg.Done()

		symbolData, err := massaClient.DatastoreAssetSymbol(contractAddress)
		if err != nil {
			errCh <- fmt.Errorf("failed to fetch asset symbol: %w", err)
			return
		}
		symbolCh <- symbolData
	}()

	// Concurrently fetch asset decimals
	go func() {
		defer wg.Done()

		decimals, err := massaClient.DatastoreAssetDecimals(contractAddress)
		if err != nil {
			errCh <- fmt.Errorf("failed to fetch asset decimals: %w", err)
			return
		}
		decimalsCh <- decimals
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
		Address:  contractAddress,
		Name:     name,
		Symbol:   symbol,
		Decimals: swag.Int64(int64(decimals)),
	}

	return assetInfoFromSC, nil
}
