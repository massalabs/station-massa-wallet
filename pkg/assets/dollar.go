package assets

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"strconv"
)

const url = "https://api.mexc.com/api/v3/avgPrice?symbol="

type AvgPrice struct {
	Price string `json:"price"`
}

func DollarValue(balance string, symbol string, decimals int64) (string, error) {
	ticker, err := SymbolToTicker(symbol)
	if err != nil {
		return "", fmt.Errorf("Error converting symbol to ticker: %s\n", err)
	}

	price := 0.0

	if ticker == "USD" {
		price = 1.0
	} else {
		price, err = DollarPrice(ticker)
		if err != nil {
			return "", fmt.Errorf("Error getting dollar price: %s\n", err)
		}
	}

	balanceInt, success := new(big.Int).SetString(balance, 10)
	if !success {
		return "", fmt.Errorf("error converting balance to big.Int")
	}

	// Calculate 10^decimals as a big.Float because big.Int doesn't support floating point operations
	powTenDecimals := new(big.Float).SetFloat64(math.Pow(10, float64(decimals)))

	// Convert balanceInt to big.Float for division
	balanceFloat := new(big.Float).SetInt(balanceInt)

	// Divide balanceFloat by 10^decimals to adjust for decimals
	actualBalanceBigFloat := new(big.Float).Quo(balanceFloat, powTenDecimals)

	// Convert actualBalanceBigFloat to a float64 to multiply with price (assuming price is not a very large number)
	actualBalance, _ := actualBalanceBigFloat.Float64()

	// Calculate dollar value
	dollarValue := actualBalance * price

	return fmt.Sprintf("%.2f", dollarValue), nil
}

// DollarPrice returns the value of 1 MAS in USD
func DollarPrice(ticker string) (float64, error) {
	resp, err := http.Get(url + ticker)
	if err != nil {
		return 0.0, fmt.Errorf("Error making HTTP request: %s\n", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return 0.0, fmt.Errorf("Error reading response body: %s\n", err)
	}

	var data AvgPrice

	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0.0, fmt.Errorf("Error parsing JSON: %s\n", err)
	}

	result, err := strconv.ParseFloat(data.Price, 64)
	if err != nil {
		return 0.0, fmt.Errorf("Error converting string to float: %s\n", err)
	}

	return result, nil
}

func SymbolToTicker(symbol string) (string, error) {
	switch symbol {
	case "tDAI.s", "DAI.e":
		return "USD", nil // or DAIUSDT
	case "WETH.e", "WETH.s":
		return "ETHUSDT", nil
	case "MAS", "WMAS":
		return "MASUSDT", nil
	case "USDC.s", "USDC.e":
		return "USD", nil // or USDCUSDT
	}

	return "", fmt.Errorf("Invalid symbol: %s", symbol)
}
