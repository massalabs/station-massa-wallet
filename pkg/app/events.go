package walletapp

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
	Transfer
	TradeRolls
	Unprotect
)
