package account

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/awnumar/memguard"
	"github.com/btcsuite/btcutil/base58"
	"github.com/massalabs/station-massa-wallet/pkg/crypto"
	"github.com/massalabs/station-massa-wallet/pkg/types"
	"github.com/massalabs/station-massa-wallet/pkg/types/object"
	"gopkg.in/yaml.v2"
)

const (
	LastVersion    = 1
	UnknownVersion = 0
)

var (
	ErrInvalidPrivateKey = errors.New("invalid private key")
	ErrRuntime           = errors.New("runtime execution")
	ErrInvalidParameter  = errors.New("invalid parameter")
)

type Account struct {
	Version      uint8                      `yaml:"Version"`
	Nickname     string                     `yaml:"Nickname,omitempty"`
	Address      *types.Address             `yaml:"Address,omitempty"`
	Salt         [16]byte                   `yaml:"Salt,flow"`
	Nonce        [12]byte                   `yaml:"Nonce,flow"`
	CipheredData *types.EncryptedPrivateKey `yaml:"CipheredData,flow"`
	PublicKey    *types.PublicKey           `yaml:"PublicKey,flow"`
}

func New(
	version uint8,
	nickname string,
	address *types.Address,
	salt [16]byte,
	nonce [12]byte,
	encryptedPrivateKey *types.EncryptedPrivateKey,
	publicKey *types.PublicKey,
) (*Account, error) {
	if version > LastVersion {
		return nil, fmt.Errorf("%w: the provided version is invalid: %d", ErrInvalidParameter, version)
	}

	if !NicknameIsValid(nickname) {
		return nil, fmt.Errorf("%w: the provided nickname is invalid: %s", ErrInvalidNickname, nickname)
	}

	if err := address.Validate(address.Version, object.UserAddress, object.SmartContractAddress); err != nil {
		return nil, fmt.Errorf("%w: the provided address is invalid: %w", ErrInvalidParameter, err)
	}

	if err := encryptedPrivateKey.Validate(encryptedPrivateKey.Version, object.EncryptedPrivateKey); err != nil {
		return nil, fmt.Errorf("%w: the provided encrypted private key is invalid: %w", ErrInvalidPrivateKey, err)
	}

	if err := publicKey.Validate(publicKey.Version, object.PublicKey); err != nil {
		return nil, fmt.Errorf("%w: the provided public key is invalid: %w", ErrInvalidParameter, err)
	}

	return &Account{
		Version:      version,
		Nickname:     nickname,
		Address:      address,
		Salt:         salt,
		Nonce:        nonce,
		CipheredData: encryptedPrivateKey,
		PublicKey:    publicKey,
	}, nil
}

// Generate generates a new account with a random private key. It destroys the password.
func Generate(password *memguard.LockedBuffer, nickname string) (*Account, error) {
	version := uint8(LastVersion)

	publicKeyBytes, privateKeyBytes, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRuntime, err)
	}

	var salt [16]byte

	_, err = rand.Read(salt[:])
	if err != nil {
		return nil, fmt.Errorf("generating random salt: %w", err)
	}

	var nonce [12]byte

	_, err = rand.Read(nonce[:])
	if err != nil {
		return nil, fmt.Errorf("generating random nonce: %w", err)
	}

	privateKeyBytes = append([]byte{types.EncryptedPrivateKeyLastVersion}, []byte(privateKeyBytes)...)
	privateKey := memguard.NewBufferFromBytes(privateKeyBytes)

	encryptedSecret, err := seal(privateKey, password, salt[:], nonce[:])
	if err != nil {
		return nil, fmt.Errorf("sealing secret: %w", err)
	}

	password.Destroy()

	publicKey := types.PublicKey{
		Object: &object.Object{
			Kind:    object.PublicKey,
			Version: types.PublicKeyLastVersion,
			Data:    publicKeyBytes,
		},
	}

	address := types.NewAddressFromPublicKey(&publicKey)

	cipheredData := types.EncryptedPrivateKey{
		Object: &object.Object{
			Kind:    object.EncryptedPrivateKey,
			Version: types.EncryptedPrivateKeyLastVersion,
			Data:    encryptedSecret,
		},
	}

	return New(version, nickname, address, salt, nonce, &cipheredData, &publicKey)
}

// NewFromPrivateKey creates a new account from a private key. It destroys the password.
func NewFromPrivateKey(password *memguard.LockedBuffer, nickname string, privateKeyText *memguard.LockedBuffer) (*Account, error) {
	version := uint8(LastVersion)

	var salt [16]byte

	_, err := rand.Read(salt[:])
	if err != nil {
		return nil, fmt.Errorf("generating random salt: %w", err)
	}

	var nonce [12]byte

	_, err = rand.Read(nonce[:])
	if err != nil {
		return nil, fmt.Errorf("generating random nonce: %w", err)
	}

	seed, privateKeyVersion, err := base58.CheckDecode(string(privateKeyText.Bytes()[1:])) // omit the first byte because it's 'S' for secret key
	if err != nil {
		return nil, fmt.Errorf("%w: decoding base58 private key: %w", ErrInvalidPrivateKey, err)
	}

	seedBuffer := memguard.NewBufferFromBytes(seed)

	privateKeyBytes := ed25519.NewKeyFromSeed(seedBuffer.Bytes())
	seedBuffer.Destroy()

	privateKeyBytes = append([]byte{privateKeyVersion}, privateKeyBytes...)
	privateKey := memguard.NewBufferFromBytes(privateKeyBytes)

	encryptedSecret, err := seal(privateKey, password, salt[:], nonce[:])
	if err != nil {
		return nil, fmt.Errorf("sealing secret: %w", err)
	}

	encryptedPrivateKey := types.EncryptedPrivateKey{
		Object: &object.Object{
			Kind:    object.EncryptedPrivateKey,
			Version: privateKeyVersion,
			Data:    encryptedSecret,
		},
	}

	err = encryptedPrivateKey.Validate(privateKeyVersion, object.EncryptedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("validating encrypted private key: %w", err)
	}

	publicKey, err := encryptedPrivateKey.PublicKey(password, salt[:], nonce[:])
	password.Destroy()

	if err != nil {
		return nil, fmt.Errorf("getting public key from encrypted private key: %w", err)
	}

	address := types.NewAddressFromPublicKey(publicKey)

	return New(version, nickname, address, salt, nonce, &encryptedPrivateKey, publicKey)
}

// seal encrypts the private key with the password.
func seal(privateKey, password *memguard.LockedBuffer, salt, nonce []byte) ([]byte, error) {
	aeadCipher, secretKey, err := crypto.NewSecretCipher(password.Bytes(), salt[:])
	if err != nil {
		return nil, fmt.Errorf("creating secret cipher: %w", err)
	}

	encryptedSecret := crypto.SealSecret(aeadCipher, nonce[:], privateKey)

	secretKey.Destroy()

	return encryptedSecret, nil
}

// PrivateKeyTextInClear returns the private key in clear and destroys the password.
func (a *Account) PrivateKeyTextInClear(password *memguard.LockedBuffer) (*memguard.LockedBuffer, error) {
	return a.CipheredData.PrivateKeyTextInClear(password, a.Salt[:], a.Nonce[:])
}

func (a *Account) PrivateKeyBytesInClear(password *memguard.LockedBuffer) (*memguard.LockedBuffer, error) {
	return a.CipheredData.PrivateKeyBytesInClear(password, a.Salt[:], a.Nonce[:])
}

// Sign signs the data with the private key and destroys the password.
func (a *Account) Sign(password *memguard.LockedBuffer, data []byte) ([]byte, error) {
	return a.CipheredData.Sign(password, a.Salt[:], a.Nonce[:], data)
}

// SignWithPrivateKey signs the data with the private key and destroys the private key.
func (a *Account) SignWithPrivateKey(privateKey *memguard.LockedBuffer, data []byte) ([]byte, error) {
	return a.CipheredData.SignWithPrivateKey(privateKey, data)
}

func (a *Account) Marshal() ([]byte, error) {
	return yaml.Marshal(a)
}

// HasAccess returns true if the password is valid for the account. It destroys the password.
func (a *Account) HasAccess(password *memguard.LockedBuffer) bool {
	return a.CipheredData.HasAccess(password, a.Salt[:], a.Nonce[:])
}

func (a *Account) Unmarshal(data []byte) error {
	err := yaml.Unmarshal(data, &a)
	if err != nil {
		return err
	}

	if len(a.Salt) == 0 || a.Salt == [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		return fmt.Errorf("missing salt")
	}

	if len(a.Nonce) == 0 || a.Nonce == [12]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		return fmt.Errorf("missing nonce")
	}

	if a.CipheredData == nil || len(a.CipheredData.Data) == 0 {
		return fmt.Errorf("missing ciphered data")
	}

	if a.PublicKey == nil || len(a.PublicKey.Data) == 0 {
		return fmt.Errorf("missing public key")
	}

	if a.Version == UnknownVersion {
		return fmt.Errorf("invalid or missing version")
	}

	return nil
}
