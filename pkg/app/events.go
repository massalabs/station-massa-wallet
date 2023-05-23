package walletapp

const (
	PromptResultEvent  string = "promptResult"
	PromptRequestEvent string = "promptRequest"
)

type EventData struct {
	Success     bool
	CodeMessage string
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
)

type PromptRequestData struct {
	Action PromptRequestAction
	Msg    string
	Data   interface{}
}
