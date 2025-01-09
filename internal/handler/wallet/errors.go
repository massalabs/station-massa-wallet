package wallet

const (
	_ = "Wallet-" + string('0'+(iota)/1000%10) + string('0'+(iota)/100%10) + string('0'+(iota)/10%10) + string('0'+(iota)/1%10)
	errorWrongPassword
	errorGetAccount
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
	errorSign
	errorSignLoadCache
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
	errorTradeRoll
	internalError
	errorAddSignRule
	errorDeleteSignRule
	errorUpdateSignRule
	errorCachePrivateKey
	errorGetAccountConfig
)
