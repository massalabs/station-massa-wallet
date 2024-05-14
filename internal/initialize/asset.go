package initialize

import "github.com/massalabs/station-massa-wallet/pkg/assets"

func Asset() error {
	err := assets.InitDefaultAsset()
	if err != nil {
		return err
	}

	return nil
}
