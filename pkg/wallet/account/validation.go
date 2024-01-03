package account

import (
	"errors"
	"regexp"
)

const (
	MaxNicknameLength = 32
)

var ErrInvalidNickname = errors.New("nickname invalid")

// NicknameIsValidForLoading validates the nickname using the following rules:
// - must have at least 1 character
// - must contain only alphanumeric characters, underscores and dashes
// - must not exceed 53 characters
// The length is 53 because for importing or loading an account file, we want to accept the address as a nickname,
// because massa-client uses the address as the nickname when saving the account file.
func NicknameIsValidForLoading(nickname string) bool {
	return CheckAlphanumeric(nickname) && len(nickname) <= 53
}

// NicknameIsValid validates the nickname using the following rules:
// - must have at least 1 character
// - must contain only alphanumeric characters, underscores and dashes
// - must not exceed MaxNicknameLength characters
func NicknameIsValid(nickname string) bool {
	return CheckAlphanumeric(nickname) && len(nickname) <= MaxNicknameLength
}

func CheckAlphanumeric(str string) bool {
	regex := regexp.MustCompile("^[a-zA-Z0-9-_]+$")
	return regex.MatchString(str)
}
