package wallet

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
)

func NewWalletExportFile(wallet *wallet.Wallet) operations.ExportAccountFileHandler {
	return &walletExportFile{wallet: wallet}
}

type walletExportFile struct {
	wallet *wallet.Wallet
}

// Handle handles an export file request
// It will serve the yaml file so that the client can download it.
func (w *walletExportFile) Handle(params operations.ExportAccountFileParams) middleware.Responder {
	// params.Nickname length is already checked by go swagger
	acc, errResp := loadAccount(w.wallet, params.Nickname)
	if errResp != nil {
		return errResp
	}

	pathToAccount, err := w.wallet.AccountPath(acc.Nickname)
	if err != nil {
		return operations.NewExportAccountFileInternalServerError().WithPayload(
			&models.Error{
				Code:    errorExportWallet,
				Message: err.Error(),
			})
	}

	responder := middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		file, err := os.Open(pathToAccount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filepath.Base(file.Name())))

		if _, err := io.Copy(w, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
	})

	return responder
}
