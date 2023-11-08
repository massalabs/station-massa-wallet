package account

import (
	"testing"

	"github.com/awnumar/memguard"
	"github.com/stretchr/testify/assert"
)

const (
	nickname       = "bonjour"
	password       = "bonjour"
	privateKeyText = "S12eCL2rGvRT4wZKaH7KdLd7fuhCF1Vt34SrNnRDEtduMZrjMxHz"
)

func NewAccount(t *testing.T) *Account {
	// Create test values for the password and nickname
	samplePassword := memguard.NewBufferFromBytes([]byte(password))
	privateKey := memguard.NewBufferFromBytes([]byte(privateKeyText))

	// Call the NewFromPrivateKey function with the test values
	account, err := NewFromPrivateKey(samplePassword, nickname, privateKey)
	assert.NoError(t, err)

	return account
}
