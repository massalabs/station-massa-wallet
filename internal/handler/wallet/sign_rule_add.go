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

type SignRulePromptData struct {
	Nickname      string
	WalletAddress string
	Description   string
	SignRule      config.SignRule
}

func NewAddSignRuleHandler(prompterApp prompt.WalletPrompterInterface) operations.AddSignRuleHandler {
	return &addSignRuleHandler{prompterApp: prompterApp}
}

type addSignRuleHandler struct {
	prompterApp prompt.WalletPrompterInterface
}

func (w *addSignRuleHandler) Handle(params operations.AddSignRuleParams) middleware.Responder {
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

	if exists := cfg.IsExistingRule(acc.Nickname, newRule); exists {
		return newErrorResponse("Rule already exists", errorAddSignRule, http.StatusBadRequest)
	}

	if err := config.ValidateRule(newRule); err != nil {
		return operations.NewAddSignRuleBadRequest().WithPayload(&models.Error{Code: errorInvalidAssetAddress, Message: err.Error()})
	}

	promptRequest, err := w.getPromptRequest(params, acc)
	if err != nil {
		return newErrorResponse(fmt.Sprintf("Error: %v", err.Error()), errorAddSignRule, http.StatusBadRequest)
	}

	password, err := PromptPassword(w.prompterApp, acc, promptRequest)
	if err != nil {
		msg := fmt.Sprintf("Unable to unprotect wallet: %s", err.Error())
		if errors.Is(err, utils.ErrWrongPassword) || errors.Is(err, utils.ErrActionCanceled) {
			return newErrorResponse(msg, errorGetWallets, http.StatusUnauthorized)
		}

		return newErrorResponse(msg, errorGetWallets, http.StatusInternalServerError)
	}

	ruleID, err := cfg.AddSignRule(acc.Nickname, newRule)
	if err != nil {
		return newErrorResponse(fmt.Sprintf("Error: %v", err.Error()), errorAddSignRule, http.StatusInternalServerError)
	}

	if cfg.HasEnabledRule(acc.Nickname) {
		err = cache.CachePrivateKeyFromPassword(acc, password)
		if err != nil {
			return newErrorResponse(err.Error(), errorCachePrivateKey, http.StatusInternalServerError)
		}
	}

	return w.Success(ruleID)
}

func (w *addSignRuleHandler) Success(ruleID string) middleware.Responder {
	return operations.NewAddSignRuleOK().WithPayload(
		&models.AddSignRuleResponse{
			ID: ruleID,
		})
}

func (w *addSignRuleHandler) getPromptRequest(params operations.AddSignRuleParams, acc *account.Account) (*prompt.PromptRequest, error) {
	address, err := acc.Address.String()
	if err != nil {
		return nil, fmt.Errorf("failed to stringify address: %w", err)
	}

	promptData := SignRulePromptData{
		WalletAddress: address,
		Nickname:      acc.Nickname,
		Description:   params.Body.Description,
		SignRule: config.SignRule{
			Name:     params.Body.Name,
			Contract: *params.Body.Contract,
			RuleType: config.RuleType(params.Body.RuleType),
			Enabled:  *params.Body.Enabled,
		},
	}

	promptRequest := prompt.PromptRequest{
		Action: walletapp.AddSignRule,
		Data:   promptData,
	}

	return &promptRequest, nil
}
