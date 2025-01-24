package walletapp

// Enum Bindings

type EventType string

const (
	PromptResultEvent  string = "PROMPT_RESULT"
	PromptDataEvent    string = "PROMPT_DATA"
	PromptRequestEvent string = "PROMPT_REQUEST"
)

var EventTypes = []struct {
	Value  EventType
	TSName string
}{
	{EventType(PromptResultEvent), "promptResult"},
	{EventType(PromptDataEvent), "promptData"},
	{EventType(PromptRequestEvent), "promptRequest"},
}

type PromptRequestAction string

const (
	Delete         PromptRequestAction = "DELETE_ACCOUNT"
	NewPassword    PromptRequestAction = "CREATE_PASSWORD"
	Sign           PromptRequestAction = "SIGN"
	Import         PromptRequestAction = "IMPORT_ACCOUNT"
	Backup         PromptRequestAction = "BACKUP_ACCOUNT"
	TradeRolls     PromptRequestAction = "TRADE_ROLLS"
	Unprotect      PromptRequestAction = "UNPROTECT"
	AddSignRule    PromptRequestAction = "ADD_SIGN_RULE"
	DeleteSignRule PromptRequestAction = "DELETE_SIGN_RULE"
	UpdateSignRule PromptRequestAction = "UPDATE_SIGN_RULE"
)

var PromptRequest = []struct {
	Value  PromptRequestAction
	TSName string
}{
	{Delete, "delete"},
	{NewPassword, "newPassword"},
	{Sign, "sign"},
	{Import, "import"},
	{Backup, "backup"},
	{TradeRolls, "tradeRolls"},
	{Unprotect, "unprotect"},
	{AddSignRule, "addSignRule"},
	{DeleteSignRule, "deleteSignRule"},
	{UpdateSignRule, "updateSignRule"},
}
