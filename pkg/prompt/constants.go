package prompt

import (
	"fmt"
	"time"
)

const (
	InvalidAccountFileErr = "invalid account file"
	AccountLoadErr        = "unable to load account file"
	ImportPrivateKeyErr   = "unable to import private key"
	ActionCanceledErr     = "Action canceled by user"
	TimeoutErr            = "Password prompt reached timeout"
	PasswordLengthErr     = "password length must be %d characters minimum"
	UnprotectErr          = "error unprotecting wallet"
	InputTypeErr          = "Invalid prompt input type"
	AlreadyListeningErr   = "prompter is already listening"
	UserChoiceErr         = "Invalid user choice input"
)

var passwordLengthErr = fmt.Sprintf(PasswordLengthErr, PASSWORD_MIN_LENGTH)

const TIMEOUT = 5 * time.Minute
