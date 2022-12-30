package wallet

import (
	"log"
	"os"
	"path/filepath"
)

// Clean up test data by listing all created wallets with tests and deleting them.
func cleanupTestData(walletNicknames []string) error {
	log.Printf("\n................... Start cleanupTestData ...................\n")
	// get the current working directory
	path, err := os.Getwd()
	if err != nil {
		return err
	}

	// loop through the list of file names
	for _, fileName := range walletNicknames {
		// construct the full path to the file
		fullPath := filepath.Join(path, "wallet_"+fileName+".json")

		// delete the file
		err := os.Remove(fullPath)
		log.Printf("\n................... %s cleaned ...................\n", "wallet_"+fileName)
		if err != nil {
			return err
		}
	}
	log.Printf("\n................... cleanupTestData complete ...................\n")
	return nil
}
