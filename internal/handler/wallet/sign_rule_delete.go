package wallet

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/config"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/pkg/errors"
)

func NewDeleteSignRuleHandler(prompterApp prompt.WalletPrompterInterface) operations.DeleteSignRuleHandler {
	return &deleteSignRuleHandler{prompterApp: prompterApp}
}

type deleteSignRuleHandler struct {
	prompterApp prompt.WalletPrompterInterface
}

func (w *deleteSignRuleHandler) Handle(params operations.DeleteSignRuleParams) middleware.Responder {
	acc, errResp := loadAccount(w.prompterApp.App().Wallet, params.Nickname)
	if errResp != nil {
		return errResp
	}

	cfg := config.Get()

	signRule := cfg.GetSignRule(acc.Nickname, params.RuleID)
	if signRule == nil {
		return newErrorResponse(fmt.Sprintf("Rule ID %s not found", params.RuleID), errorDeleteSignRule, http.StatusInternalServerError)
	}

	promptRequest, err := w.getPromptRequest(acc, signRule)
	if err != nil {
		return newErrorResponse(fmt.Sprintf("Error: %v", err.Error()), errorDeleteSignRule, http.StatusBadRequest)
	}

	_, err = PromptPassword(w.prompterApp, acc, promptRequest)
	if err != nil {
		msg := fmt.Sprintf("Unable to unprotect wallet: %s", err.Error())
		if errors.Is(err, utils.ErrWrongPassword) || errors.Is(err, utils.ErrActionCanceled) {
			return newErrorResponse(msg, errorGetWallets, http.StatusUnauthorized)
		}

		return newErrorResponse(msg, errorGetWallets, http.StatusInternalServerError)
	}

	err = cfg.DeleteSignRule(acc.Nickname, params.RuleID)
	if err != nil {
		return newErrorResponse(fmt.Sprintf("Error: %v", err.Error()), errorDeleteSignRule, http.StatusInternalServerError)
	}

	return w.Success()
}

func (w *deleteSignRuleHandler) Success() middleware.Responder {
	return operations.NewDeleteSignRuleOK()
}

func (w *deleteSignRuleHandler) getPromptRequest(acc *account.Account, signRule *config.SignRule) (*prompt.PromptRequest, error) {
	address, err := acc.Address.String()
	if err != nil {
		return nil, fmt.Errorf("failed to stringify address: %w", err)
	}

	promptData := SignRulePromptData{
		WalletAddress: address,
		Nickname:      acc.Nickname,
		SignRule:      *signRule,
	}

	promptRequest := prompt.PromptRequest{
		Action: walletapp.DeleteSignRule,
		Data:   promptData,
	}

	return &promptRequest, nil
}
