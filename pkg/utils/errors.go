package utils

import (
	"errors"

	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/massalabs/station-massa-wallet/pkg/walletmanager"
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
	ErrPromptInputType   = "InvalidPromptInput-0001"
)

// Message codes
const (
	WrongPassword  = "WrongPassword-0001"
	ActionCanceled = "ActionCanceled-0001"
)

// Messages
const (
	MsgTransferRequest = "transfer-request"
)

// Sentinel errors
var (
	ErrCorrelationIDNotFound = errors.New("Correlation ID not found")
	ErrCache                 = errors.New("Error loading cache")
	ErrWrongPassword         = errors.New("wrong password")
	ErrActionCanceled        = errors.New("Action canceled by user")
	ErrTimeout               = errors.New("Password prompt reached timeout")
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

	if errors.Is(err, ErrActionCanceled) {
		return ActionCanceled
	}

	return ErrUnknown
}
