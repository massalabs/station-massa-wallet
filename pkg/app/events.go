package walletapp

const (
	PasswordResultEvent string = "passwordResult" // TODO: rename to ResultEvent
	PromptRequesEvent   string = "promptRequest"
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

type PromptRequest int

const (
	Password PromptRequest = iota
	NewPassword
	Sign
	Import
	Export
	Transfer
	TradeRolls
)

type PromptRequestData struct {
	Action PromptRequest
	Msg    string
	Data   interface{}
}
