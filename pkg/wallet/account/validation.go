package account

import (
	"errors"
	"regexp"
)

const (
	MaxNicknameLength = 32
)

var ErrInvalidNickname = errors.New("nickname invalid")

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
