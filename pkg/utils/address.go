package utils

import "github.com/massalabs/station/pkg/node/base58"

// IsValidAddress checks if the address is valid based on the prefix rule, non-empty rule, and successful decoding.
func IsValidAddress(addr string) bool {
	if addr == "" {
		return false
	}

	addressPrefix := addr[:2]
	addressWithoutPrefix := addr[2:]

	if addressPrefix == "AS" && len(addressWithoutPrefix) > 0 {
		_, _, err := base58.VersionedCheckDecode(addressWithoutPrefix)
		if err == nil {
			return true
		}
	}

	return false
}
