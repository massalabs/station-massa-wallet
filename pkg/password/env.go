package password

// In this file is implemented an environment version of the password.PasswordAsker interface.

import (
	"errors"
	"os"
)

// EnvPrompter is a struct that implements the password.PasswordAsker interface.
type EnvPrompter struct{}

// Ask checks if the password is set in the environment variable WALLET_${nickname}_PASSWORD.
//
// It returns the entered password and any error that may have occurred.
func (e *EnvPrompter) Ask(name string) (string, error) {
	password := os.Getenv("WALLET_" + name + "_PASSWORD")
	if password == "" {
		return "", errors.New("password not found in environment variable")
	}
	return password, nil
}

// NewEnvPrompter creates a new password prompter.
func NewEnvPrompter() *EnvPrompter {
	return &EnvPrompter{}
}

// Verifies at compilation time that EnvPrompter implements Asker interface.
var _ Asker = &EnvPrompter{}
