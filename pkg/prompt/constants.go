package prompt

import (
	"errors"
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
	InputTypeErr          = "Invalid prompt input type"
	AlreadyListeningErr   = "prompter is already listening"
)

const (
	PASSWORD_MIN_LENGTH = 5
	TIMEOUT             = 5 * time.Minute
)

var passwordLengthErrMsg = fmt.Sprintf(PasswordLengthErr, PASSWORD_MIN_LENGTH)

var (
	ErrWrongPassword  = errors.New("wrong password")
	ErrActionCanceled = errors.New(ActionCanceledErr)
	ErrTimeout        = errors.New(TimeoutErr)
)
