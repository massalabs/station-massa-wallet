package wallet

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/cache"
	"github.com/massalabs/station-massa-wallet/pkg/config"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/pkg/errors"
)

func NewUpdateSignRuleHandler(prompterApp prompt.WalletPrompterInterface) operations.UpdateSignRuleHandler {
	return &updateSignRuleHandler{prompterApp: prompterApp}
}

type updateSignRuleHandler struct {
	prompterApp prompt.WalletPrompterInterface
}

func (w *updateSignRuleHandler) Handle(params operations.UpdateSignRuleParams) middleware.Responder {
	acc, errResp := loadAccount(w.prompterApp.App().Wallet, params.Nickname)
	if errResp != nil {
		return errResp
	}

	newRule := config.SignRule{
		Name:     params.Body.Name,
		Contract: *params.Body.Contract,
		RuleType: config.RuleType(params.Body.RuleType),
		Enabled:  *params.Body.Enabled,
	}

	cfg := config.Get()

	if signRule := cfg.GetSignRule(acc.Nickname, params.RuleID); signRule == nil {
		return newErrorResponse(fmt.Sprintf("Rule ID %s not found", params.RuleID), errorUpdateSignRule, http.StatusInternalServerError)
	}

	if err := config.ValidateRule(newRule); err != nil {
		return operations.NewUpdateSignRuleBadRequest().WithPayload(&models.Error{Code: errorInvalidAssetAddress, Message: err.Error()})
	}

	promptRequest, err := w.getPromptRequest(acc, &newRule)
	if err != nil {
		return newErrorResponse(fmt.Sprintf("Error: %v", err.Error()), errorUpdateSignRule, http.StatusBadRequest)
	}

	password, err := PromptPassword(w.prompterApp, acc, promptRequest)
	if err != nil {
		msg := fmt.Sprintf("Unable to unprotect wallet: %s", err.Error())
		if errors.Is(err, utils.ErrWrongPassword) || errors.Is(err, utils.ErrActionCanceled) {
			return newErrorResponse(msg, errorGetWallets, http.StatusUnauthorized)
		}

		return newErrorResponse(msg, errorGetWallets, http.StatusInternalServerError)
	}

	ruleID, err := cfg.UpdateSignRule(acc.Nickname, params.RuleID, newRule)
	if err != nil {
		return newErrorResponse(fmt.Sprintf("Error: %v", err.Error()), errorUpdateSignRule, http.StatusInternalServerError)
	}

	if cfg.HasEnabledRule(acc.Nickname) {
		err = cache.CachePrivateKeyFromPassword(acc, password)
		if err != nil {
			return newErrorResponse(err.Error(), errorCachePrivateKey, http.StatusInternalServerError)
		}
	}

	return w.Success(ruleID)
}

func (w *updateSignRuleHandler) Success(ruleID string) middleware.Responder {
	return operations.NewUpdateSignRuleOK().WithPayload(
		&models.UpdateSignRuleResponse{
			ID: ruleID,
		})
}

func (w *updateSignRuleHandler) getPromptRequest(acc *account.Account, signRule *config.SignRule) (*prompt.PromptRequest, error) {
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
		Action: walletapp.UpdateSignRule,
		Data:   promptData,
	}

	return &promptRequest, nil
}
