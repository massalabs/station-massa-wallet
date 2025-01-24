package utils

import "github.com/massalabs/station/pkg/node/base58"

func IsValidAddress(addr string) bool {
	if len(addr) <= 2 {
		return false
	}

	addressWithoutPrefix := addr[2:]

	_, _, err := base58.VersionedCheckDecode(addressWithoutPrefix)

	return err == nil
}

func IsValidContract(addr string) bool {
	if !IsValidAddress(addr) {
		return false
	}
	addressPrefix := addr[:2]

	return addressPrefix == "AS"
}
