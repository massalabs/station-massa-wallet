package wallet

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
)

// HandleExportFile handles a export file request
// It will serve the yaml file so that the client can download it.
func HandleExportFile(params operations.ExportAccountFileParams) middleware.Responder {
	// params.Nickname length is already checked by go swagger
	wlt, resp := loadWallet(params.Nickname)
	if resp != nil {
		return resp
	}

	pathToWallet, err := wlt.FilePath()
	if err != nil {
		pathToWallet = "wallet_unknown.yaml"
	}

	tempFile, err := os.CreateTemp("", "wallet-*.dat")
	if err != nil {
		return operations.NewExportAccountFileInternalServerError().WithPayload(
			&models.Error{
				Code:    errorExportWallet,
				Message: err.Error(),
			})
	}
	defer tempFile.Close()

	err = wlt.SaveFile(tempFile.Name())
	if err != nil {
		return operations.NewExportAccountFileInternalServerError().WithPayload(
			&models.Error{
				Code:    errorExportWallet,
				Message: err.Error(),
			})
	}

	responder := middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		file, err := os.Open(pathToWallet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filepath.Base(file.Name())))
		if _, err := io.Copy(w, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	return responder
}
