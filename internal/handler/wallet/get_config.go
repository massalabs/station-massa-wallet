package wallet

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/config"
)

func NewGetConfig() operations.GetConfigHandler {
	return &getConfig{}
}

type getConfig struct{}

func (w *getConfig) Handle(_ operations.GetConfigParams) middleware.Responder {
	cfg := config.Get()
	modelConfig, err := newConfigModel(cfg)
	//nolint:wsl
	if err != nil {
		return operations.NewGetConfigInternalServerError().WithPayload(
			&models.Error{
				Code:    internalError,
				Message: "Unable to create config model",
			})
	}

	return operations.NewGetConfigOK().WithPayload(modelConfig)
}

func newConfigModel(cfg *config.Config) (*models.Config, error) {
	modelAccounts := make(map[string]models.AccountConfig)

	for nickname, accountConfig := range cfg.Accounts {
		modelSignRules := make([]*models.SignRule, len(accountConfig.SignRules))

		for i := range accountConfig.SignRules {
			rule := accountConfig.SignRules[i]

			modelSignRules[i] = &models.SignRule{
				ID:       &rule.ID,
				Name:     &rule.Name,
				Contract: &rule.Contract,
				Enabled:  &rule.Enabled,
				RuleType: (models.RuleType)(rule.RuleType),
			}
		}

		modelAccounts[nickname] = models.AccountConfig{
			SignRules: modelSignRules,
		}
	}

	return &models.Config{
		Accounts: modelAccounts,
	}, nil
}
