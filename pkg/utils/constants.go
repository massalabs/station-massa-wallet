package utils

// Error codes
const (
	ErrInvalidNickname      = "Nickname-0001"
	ErrInvalidPassword      = "Password-0001"
	ErrInvalidPrivateKey    = "PrivateKey-0001"
	ErrAccountFile          = "AccountFile-0001" // for errors related to folder, read/write file, unmarshal...
	ErrNoFile               = "NoFile-0001"
	ErrInvalidFileFormat    = "InvalidFileFormat-0001"    // for errors related to file content
	ErrInvalidFileExtension = "InvalidFileExtension-0001" // for errors related to file name extension
	ErrDuplicateKey         = "DuplicateKey-0001"
	ErrUnknown              = "Unknown-0001"
	ErrDuplicateNickname    = "DuplicateNickname-001"
	ErrTimeout              = "Timeout-0001"
	ErrNetwork              = "Network-0001"
	ErrPromptInputType      = "InvalidPromptInput-0001"
)

// Message codes
const (
	WrongPassword  = "WrongPassword-0001"
	ActionCanceled = "ActionCanceled-0001"
)

// Messages
const (
	MsgAccountCreated     = "New password created"
	MsgAccountDeleted     = "Delete Success"
	MsgAccountUnprotected = "Unprotect Success"
	MsgAccountImported    = "Import Success"
	MsgTransferSuccess    = "Transfer Success"
	MsgRollTradeSuccess   = "Trade rolls Success"
	MsgBackupSuccess      = "Backup Success"
)
