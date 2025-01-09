package utils

import (
	"errors"

	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
)

// Error codes
const (
	ErrInvalidNickname   = "Nickname-0001"
	ErrInvalidPrivateKey = "PrivateKey-0001"
	ErrAccountFile       = "AccountFile-0001"
	ErrNoFile            = "NoFile-0001"
	ErrDuplicateKey      = "DuplicateKey-0001"
	ErrUnknown           = "Unknown-0001"
	ErrDuplicateNickname = "DuplicateNickname-001"
	ErrTimeoutMsg        = "Timeout-0001"
	ErrNetwork           = "Network-0001"
)

// Message codes
const (
	WrongPassword  = "WrongPassword-0001"
	ActionCanceled = "ActionCanceled-0001"
	InvalidFees    = "InvalidFees-0001"
)

// Messages
const (
	MsgTransferRequest = "transfer-request"
)

// Sentinel errors
var (
	ErrPrivateKeyCache  = errors.New("private key not found in cache")
	ErrCache            = errors.New("error loading cache")
	ErrWrongPassword    = errors.New("wrong password")
	ErrActionCanceled   = errors.New("action canceled by user")
	ErrInvalidInputType = errors.New("invalid prompt input type")
	ErrTimeout          = errors.New("password prompt reached timeout")
)

func WailsErrorCode(err error) string {
	if err == nil {
		return ""
	}

	if errors.Is(err, wallet.ErrNicknameNotUnique) {
		return ErrDuplicateNickname
	}

	if errors.Is(err, wallet.ErrAddressNotUnique) {
		return ErrDuplicateKey
	}

	if errors.Is(err, account.ErrInvalidNickname) {
		return ErrInvalidNickname
	}

	if errors.Is(err, account.ErrInvalidPrivateKey) {
		return ErrInvalidPrivateKey
	}

	if errors.Is(err, wallet.ErrUnmarshalAccount) {
		return ErrAccountFile
	}

	if errors.Is(err, ErrActionCanceled) {
		return ActionCanceled
	}

	return ErrUnknown
}
