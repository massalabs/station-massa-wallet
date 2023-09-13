package wallet

const (
	_ = "Wallet-" + string('0'+(iota)/1000%10) + string('0'+(iota)/100%10) + string('0'+(iota)/10%10) + string('0'+(iota)/1%10)
	errorWrongPassword
	errorGetWallet
	errorCanceledAction
	errorCreateNoNickname
	errorCreateNoPassword
	errorCreateNew
	errorDeleteNoNickname
	errorDeleteFile
	errorGetWallets
	errorSignEmptyNickname
	errorSignDecodeOperation
	errorSignDecodeMessage
	errorSignRead
	errorSignLoadCache
	errorSignCorrelationIdNotFound
	errorSignGenerateCorrelationId
	errorImportWalletCanceled
	errorImportNickNameAlreadyTaken
	errorImportWallet
	errorExportWallet
	errorTransferCoin
	errorSaveAccount
	errorInvalidAssetAddress
	errorAssetExists
	errorFetchAssetSC
	errorAddAssetJSON
	errorFetchAssetBalance
	errorAssetNotExists
	errorDeleteAssetJSON
	errorAssetNotFound
)
