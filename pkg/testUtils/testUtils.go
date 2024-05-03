package testutils

import (
	"path/filepath"

	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/massalabs/station/pkg/logger"
)

const (
	LogFileNameTest = "station-massa-wallet-test.log"
)

func LoggerTest() error {
	logPath, err := wallet.Path()
	if err != nil {
		return err
	}

	err = logger.InitializeGlobal(filepath.Join(logPath, LogFileNameTest))
	if err != nil {
		return err
	}

	return nil
}
