package wallet

import (
	"fmt"
	"net/http"
	"net/url"

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

	if newRule.RuleType == config.RuleTypeAutoSign {
		origin := params.Body.AuthorizedOrigin

		if origin == nil || len(*origin) == 0 {
			detectedOrigin, err := getOrigin(params.HTTPRequest)
			if err != nil {
				return newErrorResponse(err.Error(), errorAddSignRule, http.StatusBadRequest)
			}
			origin = detectedOrigin
			params.Body.AuthorizedOrigin = origin
		}

		newRule.AuthorizedOrigin = origin
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

	ruleType := config.RuleType(params.Body.RuleType)
	signRule := config.SignRule{
		Name:     params.Body.Name,
		Contract: *params.Body.Contract,
		RuleType: ruleType,
		Enabled:  *params.Body.Enabled,
	}

	if ruleType == config.RuleTypeAutoSign {
		signRule.AuthorizedOrigin = params.Body.AuthorizedOrigin
	}

	promptData := SignRulePromptData{
		WalletAddress: address,
		Nickname:      acc.Nickname,
		Description:   params.Body.Description,
		SignRule:      signRule,
	}

	promptRequest := prompt.PromptRequest{
		Action: walletapp.AddSignRule,
		Data:   promptData,
	}

	return &promptRequest, nil
}

const (
	originHeader     = "Origin"
	refererHeader    = "Referer"
	errMissingOrigin = "missing Origin and Referer headers"
	errInvalidURL    = "invalid URL format"
	errInvalidOrigin = "origin must only contain scheme, hostname, and optional port"
)

func getOrigin(r *http.Request) (*string, error) {
	var urlStr string
	if origin := r.Header.Get(originHeader); origin != "" {
		urlStr = origin
	} else if referer := r.Header.Get(refererHeader); referer != "" {
		urlStr = referer
	} else {
		return nil, errors.New(errMissingOrigin)
	}

	// Parse and validate the URL format
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, errInvalidURL)
	}

	// Validate that the URL only contains scheme, hostname, and optional port
	if parsedURL.Path != "" || parsedURL.RawQuery != "" || parsedURL.Fragment != "" {
		return nil, errors.New(errInvalidOrigin)
	}

	return &urlStr, nil
}
