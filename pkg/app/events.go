package walletapp

import "github.com/awnumar/memguard"

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

// User input interfaces for the channel

// EventInterface interface that all message types will implement
// Unused for now
type EventInterface interface{}

// BaseMessage contains the common attribute
type BaseMessage struct{}

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

type PasswordPromptOutput struct {
	Password *memguard.LockedBuffer
}

type SignPromptOutput struct {
	PasswordPromptOutput
	Fees uint64
}
