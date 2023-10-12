package utils

import (
	"errors"

	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/massalabs/station-massa-wallet/pkg/walletmanager"
)

// Sentinel errors
var (
	ErrInvalidPrivateKeySentinel = errors.New(ErrInvalidPrivateKey)
	ErrCorrelationIDNotFound     = errors.New("Correlation ID not found")
	ErrCache                     = errors.New("Error loading cache")
)

func WailsErrorCode(err error) string {
	if err == nil {
		return ""
	}

	if errors.Is(err, walletmanager.ErrNicknameNotUnique) {
		return ErrDuplicateNickname
	}

	if errors.Is(err, walletmanager.ErrAddressNotUnique) {
		return ErrDuplicateKey
	}

	if errors.Is(err, account.ErrInvalidNickname) {
		return ErrInvalidNickname
	}

	if errors.Is(err, account.ErrInvalidPrivateKey) {
		return ErrInvalidPrivateKey
	}

	if errors.Is(err, walletmanager.ErrUnmarshalAccount) {
		return ErrAccountFile
	}

	return ErrUnknown
}
