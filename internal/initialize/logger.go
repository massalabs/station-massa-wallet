package initialize

import (
	"path/filepath"

	"github.com/massalabs/station-massa-wallet/pkg/walletmanager"
	"github.com/massalabs/station/pkg/logger"
)

const (
	LogFileName = "station-massa-wallet.log"
)

func Logger() error {
	logPath, err := walletmanager.Path()
	if err != nil {
		return err
	}

	err = logger.InitializeGlobal(filepath.Join(logPath, LogFileName))
	if err != nil {
		return err
	}

	return nil
}
