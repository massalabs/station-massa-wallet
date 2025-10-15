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

	cfg := config.Get()

	signRule := cfg.GetSignRule(acc.Nickname, params.RuleID)
	if signRule == nil {
		return newErrorResponse(fmt.Sprintf("Rule ID %s not found", params.RuleID), errorUpdateSignRule, http.StatusInternalServerError)
	}

	newRule := config.SignRule{
		Name:             params.Body.Name,
		Contract:         *params.Body.Contract,
		RuleType:         config.RuleType(params.Body.RuleType),
		Enabled:          *params.Body.Enabled,
		AuthorizedOrigin: signRule.AuthorizedOrigin,
	}

	if w.preventOverwritingRule(signRule, newRule, acc.Nickname, cfg) {
		return newErrorResponse("A similar rule already exists", errorUpdateSignRule, http.StatusBadRequest)
	}

	if err := config.ValidateRule(newRule); err != nil {
		return operations.NewUpdateSignRuleBadRequest().WithPayload(&models.Error{Code: errorInvalidAssetAddress, Message: err.Error()})
	}

	promptRequest, err := w.getPromptRequest(acc, &newRule, params.Body.Description)
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

// check that the new rule is not already existing in the config.
func (w *updateSignRuleHandler) preventOverwritingRule(oldRule *config.SignRule, newRule config.SignRule, nickname string, cfg *config.Config) bool {
	if oldRule.Contract != newRule.Contract ||
		oldRule.RuleType != newRule.RuleType ||
		!utils.PtrEqual(oldRule.AuthorizedOrigin, newRule.AuthorizedOrigin) {
		return cfg.IsExistingRule(nickname, newRule)
	}

	return false
}

func (w *updateSignRuleHandler) Success(ruleID string) middleware.Responder {
	return operations.NewUpdateSignRuleOK().WithPayload(
		&models.UpdateSignRuleResponse{
			ID: ruleID,
		})
}

func (w *updateSignRuleHandler) getPromptRequest(acc *account.Account, signRule *config.SignRule, description string) (*prompt.PromptRequest, error) {
	address, err := acc.Address.String()
	if err != nil {
		return nil, fmt.Errorf("failed to stringify address: %w", err)
	}

	promptData := SignRulePromptData{
		WalletAddress: address,
		Nickname:      acc.Nickname,
		Description:   description,
		SignRule:      *signRule,
	}

	promptRequest := prompt.PromptRequest{
		Action: walletapp.UpdateSignRule,
		Data:   promptData,
	}

	return &promptRequest, nil
}
