package walletapp

import "github.com/awnumar/memguard"

const (
	PromptResultEvent  string = "promptResult"
	PromptDataEvent    string = "promptData"
	PromptRequestEvent string = "promptRequest"
)

type EventData struct {
	Success     bool
	CodeMessage string
	Data        interface{}
}

type SignPromptInput struct {
	Password string
	Fees     string
}

type SignPromptOutput struct {
	Password *memguard.LockedBuffer
	Fees     uint64
}

type PromptCtrl int

const (
	Cancel PromptCtrl = iota
)

type PromptRequestAction int

const (
	Delete PromptRequestAction = iota
	NewPassword
	Sign
	Import
	Backup
	TradeRolls
	Unprotect
)
