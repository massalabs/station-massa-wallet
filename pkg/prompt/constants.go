package prompt

import (
	"time"
)

const (
	InputTypeErr        = "invalid prompt input type"
	AlreadyListeningErr = "prompter is already listening"
)

const (
	TIMEOUT = 5 * time.Minute
)
