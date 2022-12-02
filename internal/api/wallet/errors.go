package wallet

const (
	errorCodeWalletAlreadyExists              = "Wallet-0001"
	errorCodeWalletWrongPassword              = "Wallet-0002"
	errorCodeGetWallet                        = "Wallet-0003"
	ErrorCodeWalletCanceledAction             = "Wallet-0005"
	ErrorCodeWalletPasswordEmptyExecuteFct    = "Wallet-0007"
	errorCodeWalletCreateNoNickname           = "Wallet-1001"
	errorCodeWalletCreateNoPassword           = "Wallet-1002"
	errorCodeWalletCreateNew                  = "Wallet-1003"
	errorCodeWalletDeleteNoNickname           = "Wallet-2001"
	errorCodeWalletDeleteFile                 = "Wallet-2002"
	errorCodeWalletImportNew                  = "Wallet-3001"
	errorCodeWalletGetWallets                 = "Wallet-4001"
	errorCodeWalletSignOperationEmptyNickname = "Wallet-5001"
	errorCodeWalletSignOperationRead          = "Wallet-5001"
)
