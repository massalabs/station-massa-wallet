package utils

import (
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"
)

const MasDecimals = 9

func MasToNano(masAmount string) (uint64, error) {
	dec, err := decimal.NewFromString(masAmount)
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(dec.Shift(MasDecimals).StringFixed(0), 10, 64)
}

func NanoToMas(nanoMasAmount uint64) (string, error) {
	dec, err := decimal.NewFromString(fmt.Sprint(nanoMasAmount))
	if err != nil {
		return "", err
	}

	return dec.Shift(-MasDecimals).String(), nil
}
