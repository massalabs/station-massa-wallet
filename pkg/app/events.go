package walletapp

import "github.com/awnumar/memguard"

const (
	PromptResultEvent  string = "promptResult"
	PromptDataEvent    string = "promptData"
	PromptRequestEvent string = "promptRequest"
)

// Events

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
	TradeRolls
	Unprotect
)

// User input interfaces for the channel

// CorrelationIdentifier interface that all message types will implement
type CorrelationIdentifier interface {
	GetCorrelationID() string
}

// BaseMessage contains the common CorrelationID attribute
type BaseMessage struct {
	CorrelationID string
}

// GetCorrelationID retrieves the CorrelationID from the BaseMessage struct
func (m *BaseMessage) GetCorrelationID() string {
	return m.CorrelationID
}

type StringPromptInput struct {
	BaseMessage
	Message string
}

type SignPromptInput struct {
	BaseMessage
	Password string
	Fees     string
}

type ImportPKeyPromptInput struct {
	BaseMessage
	PrivateKey *memguard.LockedBuffer
	Password   *memguard.LockedBuffer
	Nickname   string
}

// Output interface for the caller of the prompt

type SignPromptOutput struct {
	Password *memguard.LockedBuffer
	Fees     uint64
}
