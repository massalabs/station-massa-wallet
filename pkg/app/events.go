package walletapp

const (
	PromptResultEvent  string = "promptResult"
	PromptRequestEvent string = "promptRequest"
)

type EventData struct {
	Success bool
	Data    interface{}
	Error   string
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
	Export
	Transfer
	TradeRolls
)

type PromptRequestData struct {
	Action PromptRequestAction
	Msg    string
	Data   interface{}
}
