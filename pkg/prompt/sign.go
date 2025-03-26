package prompt

import (
	"fmt"
	"strconv"

	"github.com/awnumar/memguard"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/config"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
)

type PromptRequestSignData struct {
	Description          string
	Fees                 string
	MinFees              string
	OperationType        int
	Coins                string
	Address              string
	Function             string
	MaxCoins             string
	WalletAddress        string
	Nickname             string
	RollCount            uint64
	RecipientAddress     string
	RecipientNickname    string
	Amount               string
	PlainText            string
	AllowFeeEdition      bool
	ChainID              int64
	Assets               []models.AssetInfo
	Parameters           []byte
	DeployedByteCodeSize uint   // for executeSC of type deploySC
	DeployedCoins        uint64 // for executeSC of type DeploySC; the number of coins sent to the deployed contract
	EnabledSignRule      *config.RuleType
}

// handleSignPrompt returns the password as a LockedBuffer, or an error if the input is not a string.
func handleSignPrompt(prompterApp WalletPrompterInterface, req PromptRequest, input interface{}, acc *account.Account) (*walletapp.SignPromptOutput, bool, error) {
	inputObject, ok := input.(*walletapp.SignPromptInput)
	if !ok {
		return nil, true, InputTypeError(prompterApp)
	}

	fees, err := strconv.ParseUint(inputObject.Fees, 10, 64)
	if err != nil {
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.InvalidFees})

		return nil, true, fmt.Errorf("failed to parse fees: %w", err)
	}

	_, ok = req.Data.(PromptRequestSignData)

	if ok && req.DisablePassword {
		// if sign rule is enabled, we don't need to check password
		return &walletapp.SignPromptOutput{
			Fees: fees,
		}, false, nil
	}

	inputString := inputObject.Password

	password := memguard.NewBufferFromBytes([]byte(inputString))

	// password will be destroy in acc.HasAccess, so we need to create a new one.
	passwordReturned := memguard.NewBufferFromBytes([]byte(inputString))

	if acc != nil && !acc.HasAccess(password) {
		msg := fmt.Sprintf("Invalid password for account %s", acc.Nickname)

		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.WrongPassword})

		return nil, true, fmt.Errorf("%w: %s", utils.ErrWrongPassword, msg)
	}

	output := &walletapp.SignPromptOutput{
		PasswordPromptOutput: walletapp.PasswordPromptOutput{
			Password: passwordReturned,
		},
		Fees: fees,
	}

	return output, false, nil
}
