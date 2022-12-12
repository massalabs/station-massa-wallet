package wallet

const (
	_ = "Wallet-" + string('0'+(iota)/10%10) + string('0'+(iota)/1%10)
	errorAlreadyExists
	errorWrongPassword
	errorGetWallet
	errorCanceledAction
	errorPasswordEmptyExecuteFct
	errorCreateNoNickname
	errorCreateNoPassword
	errorCreateNew
	errorDeleteNoNickname
	errorDeleteFile
	errorImportNew
	errorGetWallets
	errorSignOperationEmptyNickname
	errorSignOperationRead
	errorAlreadyImported
)
