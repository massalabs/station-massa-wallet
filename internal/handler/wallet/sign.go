package wallet

import (
	"encoding/binary"
	"fmt"
	"net/http"
	"strconv"

	"github.com/awnumar/memguard"
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/assets"
	"github.com/massalabs/station-massa-wallet/pkg/cache"
	"github.com/massalabs/station-massa-wallet/pkg/config"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/node/sendoperation/buyrolls"
	"github.com/massalabs/station/pkg/node/sendoperation/callsc"
	"github.com/massalabs/station/pkg/node/sendoperation/executesc"
	"github.com/massalabs/station/pkg/node/sendoperation/sellrolls"
	"github.com/massalabs/station/pkg/node/sendoperation/transaction"
	onchain "github.com/massalabs/station/pkg/onchain"
	"github.com/pkg/errors"
)

const (
	RollPrice = 100
)

func NewSign(prompterApp prompt.WalletPrompterInterface) operations.SignHandler {
	return &walletSign{prompterApp: prompterApp}
}

type walletSign struct {
	prompterApp prompt.WalletPrompterInterface
}

func (w *walletSign) Handle(params operations.SignParams) middleware.Responder {
	acc, errResp := loadAccount(w.prompterApp.App().Wallet, params.Nickname)
	if errResp != nil {
		return errResp
	}

	promptRequest, fees, err := w.getPromptRequest(params, acc)
	if err != nil {
		return newErrorResponse(fmt.Sprintf("Error: %v", err.Error()), errorSignDecodeMessage, http.StatusBadRequest)
	}

	cfg := config.Get()

	var contract *string

	promptData, ok := promptRequest.Data.(prompt.PromptRequestSignData)
	if ok {
		contract = &promptData.Address
	}

	enabledRule := cfg.GetEnabledRuleForContract(acc.Nickname, contract)

	var privateKey *memguard.LockedBuffer

	skipPrompt := false

	if enabledRule != nil {
		promptData.EnabledSignRule = enabledRule

		// at this point, we have a rule enabled for the contract, if private key is cached, we don't need to prompt for password
		privateKey, err = cache.PrivateKeyFromCache(acc)
		if err != nil {
			logger.Warn("error retriving private key from cache: ", err)
		} else if privateKey != nil {
			// If privatekey is cached, we don't need to prompt for password
			promptRequest.DisablePassword = true

			// If the rule is AutoSign, we don't need to open wails prompt
			if *enabledRule == config.RuleTypeAutoSign {
				skipPrompt = true
			}
		}
	}

	if !skipPrompt {
		promptRequest.Data = promptData

		output, err := PromptForOperation(w.prompterApp, acc, promptRequest)
		if err != nil {
			msg := fmt.Sprintf("Unable to unprotect wallet: %s", err.Error())
			if errors.Is(err, utils.ErrWrongPassword) || errors.Is(err, utils.ErrActionCanceled) {
				return newErrorResponse(msg, errorGetWallets, http.StatusUnauthorized)
			}

			return newErrorResponse(msg, errorGetWallets, http.StatusInternalServerError)
		}

		fees = output.Fees

		if privateKey == nil {
			if output.Password != nil {
				privateKey, err = acc.PrivateKeyBytesInClear(output.Password)
				if err != nil {
					return newErrorResponse(err.Error(), errorWrongPassword, http.StatusInternalServerError)
				}
			}
		}

		if cfg.HasEnabledRule(acc.Nickname) {
			err = cache.CachePrivateKey(acc, privateKey)
			if err != nil {
				return newErrorResponse(err.Error(), errorCachePrivateKey, http.StatusInternalServerError)
			}
		}
	}

	operation, msgToSign, err := prepareOperation(acc, fees, params.Body.Operation.String(), *params.Body.ChainID)
	if err != nil {
		return newErrorResponse(err.Error(), errorSignDecodeOperation, http.StatusBadRequest)
	}

	signature := acc.SignWithPrivateKey(privateKey, msgToSign)

	if !acc.VerifySignature(msgToSign, signature) {
		return newErrorResponse("Error: signature verification failed", "errorSignVerifySignature", http.StatusInternalServerError)
	}

	return w.Success(acc, signature, operation)
}

func (w *walletSign) Success(acc *account.Account, signature []byte, operation []byte) middleware.Responder {
	publicKeyBytes, err := acc.PublicKey.MarshalText()
	if err != nil {
		return newErrorResponse(err.Error(), errorGetAccount, http.StatusInternalServerError)
	}

	return operations.NewSignOK().WithPayload(
		&models.SignResponse{
			PublicKey: string(publicKeyBytes),
			Signature: signature,
			Operation: operation,
		})
}

// prepareOperation prepares the operation to be signed.
// Returns the modified operation (fees change) and the operation to be signed (with public key).
// Returns an error if the operation cannot be decoded.
func prepareOperation(acc *account.Account, fees uint64, operationB64 string, chainID int64) ([]byte, []byte, error) {
	decodedMsg, _, expiry, err := sendoperation.DecodeMessage64(operationB64)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode operation for preparing before signing: %w", err)
	}

	operation := make([]byte, 0)
	operation = binary.AppendUvarint(operation, fees)
	operation = binary.AppendUvarint(operation, expiry)
	operation = append(operation, decodedMsg...)

	// operation to be signed
	msgToSign := make([]byte, len(operation))
	copy(msgToSign, operation)

	publicKey, err := acc.PublicKey.MarshalBinary()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to marshal public key: %w", err)
	}

	msgToSign = utils.PrepareSignData(uint64(chainID), append(publicKey, msgToSign...))

	return operation, msgToSign, nil
}

func (w *walletSign) getPromptRequest(params operations.SignParams, acc *account.Account) (*prompt.PromptRequest, uint64, error) {
	msgToSign := params.Body.Operation.String()

	decodedMsg, fees, _, err := sendoperation.DecodeMessage64(msgToSign)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode transaction message: %w", err)
	}

	opType, err := sendoperation.DecodeOperationType(decodedMsg)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode operation ID: %w", err)
	}

	var data prompt.PromptRequestSignData

	switch opType {
	case transaction.OpType:
		data, err = w.getTransactionPromptData(decodedMsg)

	case buyrolls.OpType, sellrolls.OpType:
		data, err = getRollPromptData(decodedMsg)

	case executesc.OpType:
		data, err = getExecuteSCPromptData(decodedMsg)

	case callsc.OpType:
		data, err = getCallSCPromptData(decodedMsg)

	default:
		return nil, 0, fmt.Errorf("unhandled operation type: %d", opType)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode message of operation type: %d: %w", opType, err)
	}

	address, err := acc.Address.String()
	if err != nil {
		return nil, 0, fmt.Errorf("err: %w", err)
	}

	_, minimalFees, err := network.GetNodeInfo()
	if err != nil {
		minimalFees = "0"
	}

	data.Description = params.Body.Description
	data.Fees = strconv.FormatUint(fees, 10)
	data.MinFees = minimalFees
	data.WalletAddress = address
	data.Nickname = acc.Nickname
	data.OperationType = int(opType)
	data.AllowFeeEdition = *params.AllowFeeEdition
	data.ChainID = *params.Body.ChainID
	data.Assets = convertAssetsToModel(assets.Store.All(acc.Nickname, int(data.ChainID)))

	promptRequest := prompt.PromptRequest{
		Action:          walletapp.Sign,
		Data:            data,
		DisablePassword: false,
	}

	return &promptRequest, fees, nil
}

func getCallSCPromptData(
	decodedMsg []byte,
) (prompt.PromptRequestSignData, error) {
	msg, err := callsc.DecodeMessage(decodedMsg)
	if err != nil {
		return prompt.PromptRequestSignData{}, err
	}

	return prompt.PromptRequestSignData{
		Coins:      strconv.FormatUint(msg.Coins, 10),
		Address:    msg.Address,
		Function:   msg.Function,
		Parameters: msg.Parameters,
	}, nil
}

func getExecuteSCPromptData(
	decodedMsg []byte,
) (prompt.PromptRequestSignData, error) {
	msg, err := executesc.DecodeMessage(decodedMsg)
	if err != nil {
		return prompt.PromptRequestSignData{}, err
	}

	promptReq := prompt.PromptRequestSignData{
		MaxCoins: strconv.FormatUint(msg.MaxCoins, 10),
	}

	// Check the datastore to know whether the ExecuteSC is a DeploySC or not
	if msg.DataStore == nil { // the executeSC is not a deploySC
		return promptReq, nil
	}

	dataStore, err := onchain.DeSerializeDatastore(msg.DataStore)
	if err != nil {
		return prompt.PromptRequestSignData{}, err
	}

	deployedContract, isDeployDatastore := onchain.DatastoreToDeployedContract(dataStore)
	if isDeployDatastore { // the executeSC is a deploySC
		promptReq.DeployedByteCodeSize = uint(len(deployedContract.ByteCode))
		promptReq.DeployedCoins = deployedContract.Coins
	}

	return promptReq, nil
}

func getRollPromptData(
	decodedMsg []byte,
) (prompt.PromptRequestSignData, error) {
	msg, err := sendoperation.RollDecodeMessage(decodedMsg)
	if err != nil {
		return prompt.PromptRequestSignData{}, err
	}

	return prompt.PromptRequestSignData{
		RollCount: msg.RollCount,
		Coins:     strconv.FormatUint(msg.RollCount*RollPrice, 10),
	}, nil
}

func (w *walletSign) getTransactionPromptData(
	decodedMsg []byte,
) (prompt.PromptRequestSignData, error) {
	msg, err := transaction.DecodeMessage(decodedMsg)
	if err != nil {
		return prompt.PromptRequestSignData{}, err
	}

	var recipientNickname string

	recipientAcc, err := w.prompterApp.App().Wallet.GetAccountFromAddress(msg.RecipientAddress)
	if err != nil {
		recipientNickname = ""
	} else {
		recipientNickname = recipientAcc.Nickname
	}

	return prompt.PromptRequestSignData{
		RecipientAddress:  msg.RecipientAddress,
		RecipientNickname: recipientNickname,
		Amount:            strconv.FormatUint(msg.Amount, 10),
	}, nil
}

func convertAssetsToModel(assetsWithBalance []*assets.AssetInfoWithBalances) []models.AssetInfo {
	result := make([]models.AssetInfo, 0)

	for _, asset := range assetsWithBalance {
		assetInfo := models.AssetInfo{
			Address:  asset.AssetInfo.Address,
			Decimals: asset.AssetInfo.Decimals,
			Name:     asset.AssetInfo.Name,
			Symbol:   asset.AssetInfo.Symbol,
			ChainID:  asset.AssetInfo.ChainID,
		}
		result = append(result, assetInfo)
	}

	return result
}
