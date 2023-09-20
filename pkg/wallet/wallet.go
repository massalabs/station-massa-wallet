package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/labstack/gommon/log"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
	"lukechampine.com/blake3"
)

// TODO: rename to 'account'

const (
	SecretKeyLength           = 32
	PBKDF2NbRound             = 600_000
	FileModeUserReadWriteOnly = 0o600
	Base58Version             = 0x00
	AddressVersion            = 0x00
	PubKeyVersion             = 0x00
	SignatureVersion          = 0x00
	PrivKeyVersion            = 0x00
	UserAddressPrefix         = "AU"
	PublicKeyPrefix           = "P"
	PrivateKeyPrefix          = "S"
	MaxNicknameLength         = 32
	AccountVersion            = 1
	StatusOK                  = "ok"
	StatusCorrupted           = "corrupted"
)

func ErrorAccountNotFound(nickname string) error {
	return fmt.Errorf("account '%s' not found", nickname)
}

type WalletError struct {
	Err     error
	CodeErr string // Sentinel error code from utils package, can be used as a translation key.
}

// KeyPair structure contains all the information necessary to save a key pair securely.
type KeyPair struct {
	PrivateKey VersionedKey
	PublicKey  VersionedKey
	Salt       [16]byte
	Nonce      [12]byte
}

// Wallet structure allows to link a nickname, an address and a version to one or more key pairs.
type Wallet struct {
	Version  uint8
	Nickname string
	Address  string
	KeyPair  KeyPair
	Status   string
}

type AccountSerialized struct {
	Version      *uint8       `yaml:"Version"`
	Nickname     string       `yaml:"Nickname"`
	Address      string       `yaml:"Address"`
	Salt         [16]byte     `yaml:"Salt,flow"`
	Nonce        [12]byte     `yaml:"Nonce,flow"`
	CipheredData []byte       `yaml:"CipheredData,flow"`
	PublicKey    VersionedKey `yaml:"PublicKey,flow"`
}

// toAccount returns a Wallet from an AccountSerialized.
func (accountSerialized *AccountSerialized) toAccount() (Wallet, error) {
	publicKey, err := accountSerialized.PublicKey.CheckVersion([]byte{PubKeyVersion})
	if err != nil {
		return Wallet{}, fmt.Errorf("while checking public key version: %w", err)
	}

	wallet := Wallet{
		Version:  *accountSerialized.Version,
		Nickname: accountSerialized.Nickname,
		Address:  accountSerialized.Address,
		KeyPair: KeyPair{
			PrivateKey: accountSerialized.CipheredData,
			PublicKey:  publicKey,
			Salt:       accountSerialized.Salt,
			Nonce:      accountSerialized.Nonce,
		},
	}

	return wallet, nil
}

// toAccountSerialized returns an AccountSerialized from a Wallet.
func (account *Wallet) toAccountSerialized() AccountSerialized {
	accountSerialized := AccountSerialized{
		Version:      &account.Version,
		Nickname:     account.Nickname,
		Address:      account.Address,
		Salt:         account.KeyPair.Salt,
		Nonce:        account.KeyPair.Nonce,
		CipheredData: account.KeyPair.PrivateKey, // account is protected so PrivateKey is encrypted
		PublicKey:    account.KeyPair.PublicKey,
	}

	return accountSerialized
}

// aead returns an authenticated encryption with associated data.
func aead(password []byte, salt []byte) (cipher.AEAD, error) {
	secretKey := pbkdf2.Key([]byte(password), salt, PBKDF2NbRound, SecretKeyLength, sha256.New)

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, fmt.Errorf("initializing block ciphering: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("initializing the AES block cipher wrapped in a Galois Counter Mode ciphering: %w", err)
	}

	return aesGCM, nil
}

// Protect encrypts the private key.
// The encryption algorithm used to protect the private key is AES-GCM and
// the secret key is derived from the given password using the PBKDF2 algorithm.
func (w *Wallet) Protect(password string) error {
	aead, err := aead([]byte(password), w.KeyPair.Salt[:])
	if err != nil {
		return fmt.Errorf("while protecting wallet: %w", err)
	}

	w.KeyPair.PrivateKey = aead.Seal(
		nil,
		w.KeyPair.Nonce[:],
		w.KeyPair.PrivateKey,
		nil)

	return nil
}

// Unprotect decrypts the private key using the given GUI Modal.
// The encryption algorithm used to unprotect the private key is AES-GCM and
// the secret key is derived from the given password using the PBKDF2 algorithm.
func (w *Wallet) Unprotect(password string) *WalletError {
	aead, err := aead([]byte(password), w.KeyPair.Salt[:])
	if err != nil {
		return &WalletError{fmt.Errorf("while unprotecting wallet: %w", err), utils.ErrUnknown}
	}

	pk, err := aead.Open(nil, w.KeyPair.Nonce[:], w.KeyPair.PrivateKey, nil)
	if err != nil {
		return &WalletError{fmt.Errorf("opening the private key seal: %w", err), utils.WrongPassword}
	}

	privateKey, err := VersionedKey(pk).CheckVersion([]byte{PrivKeyVersion})
	if err != nil {
		return &WalletError{fmt.Errorf("while checking private key version: %w", err), utils.ErrInvalidPrivateKey}
	}

	w.KeyPair.PrivateKey = privateKey

	return nil
}

func (w *Wallet) UnprotectFromCorrelationId(fromCache []byte, correlationId models.CorrelationID) error {
	pk, err := Xor(fromCache, correlationId)
	if err != nil {
		return fmt.Errorf("decrypt the private key: %w", err)
	}
	w.KeyPair.PrivateKey = pk

	return nil
}

func Xor(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("length of two arrays must be same")
	}
	result := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		result[i] = a[i] ^ b[i]
	}
	return result, nil
}

// Persist stores the wallet on the file system.
// Note: the wallet is stored in YAML format and in Massa Station wallet directory.
func (w *Wallet) Persist() error {
	accountSerialized := w.toAccountSerialized()

	yamlMarshaled, err := yaml.Marshal(accountSerialized)
	if err != nil {
		return fmt.Errorf("marshalling wallet: %w", err)
	}

	filePath, err := FilePath(w.Nickname)
	if err != nil {
		return fmt.Errorf("getting file path for '%s': %w", w.Nickname, err)
	}

	err = os.WriteFile(filePath, yamlMarshaled, FileModeUserReadWriteOnly)
	if err != nil {
		return fmt.Errorf("writing wallet to '%s: %w", filePath, err)
	}

	return nil
}

// MigrateWallet moves the wallet from the old location (GetWorkDir) to the new one (AccountPath).
func MigrateWallet() error {
	oldPath, err := GetWorkDir()
	if err != nil {
		return fmt.Errorf("reading config directory '%s': %w", oldPath, err)
	}

	files, err := os.ReadDir(oldPath)
	if err != nil {
		return fmt.Errorf("reading working directory '%s': %w", oldPath, err)
	}

	newPath, err := AccountPath()
	if err != nil {
		return fmt.Errorf("getting account directory '%s': %w", newPath, err)
	}

	for _, f := range files {
		fileName := f.Name()
		oldFilePath := path.Join(oldPath, fileName)

		if strings.HasPrefix(fileName, "wallet_") && strings.HasSuffix(fileName, ".yaml") {
			newFilePath := path.Join(newPath, fileName)

			// Skip if new file path exists
			if _, err := os.Stat(newFilePath); err == nil {
				continue
			}

			fmt.Println("Migrating wallet from", oldFilePath, "to", newFilePath) // Log
			err = os.Rename(oldFilePath, newFilePath)
			if err != nil {
				return fmt.Errorf("moving account file from '%s' to '%s': %w", oldFilePath, newFilePath, err)
			}
		}
	}

	return nil
}

func GetWorkDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("getting executable path: %w", err)
	}

	if runtime.GOOS == "darwin" {
		// On macOS, the executable is in a subdirectory of the working directory.
		// We need to go up 4 levels to get the working directory.
		// wallet-plugin.app/Contents/MacOS/wallet-plugin
		return filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(ex)))), nil
	}

	dir := filepath.Dir(ex)

	// Helpful when developing:
	// when running `go run`, the executable is in a temporary directory.
	if strings.Contains(dir, "go-build") {
		return ".", nil
	}

	return filepath.Dir(ex), nil
}

// AccountPath returns the path where the account yaml file are stored.
// Note: the wallet directory is the folder where the wallet plugin binary resides.
func AccountPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("getting user config directory: %w", err)
	}

	accountPath := filepath.Join(configDir, "massa-station-wallet")

	// create the directory if it doesn't exist
	if _, err := os.Stat(accountPath); os.IsNotExist(err) {
		err = os.MkdirAll(accountPath, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("creating account directory '%s': %w", accountPath, err)
		}
	}

	return accountPath, nil
}

// LoadAll loads all the wallets in the working directory.
// Note: a wallet must have: `wallet_` prefix and a `.yaml` extension.
func LoadAll() ([]Wallet, error) {
	walletDir, err := AccountPath()
	if err != nil {
		return nil, fmt.Errorf("reading config directory '%s': %w", walletDir, err)
	}

	files, err := os.ReadDir(walletDir)
	if err != nil {
		return nil, fmt.Errorf("reading working directory '%s': %w", walletDir, err)
	}

	wallets := []Wallet{}
	for _, f := range files {
		fileName := f.Name()
		filePath := path.Join(walletDir, fileName)

		if strings.HasPrefix(fileName, "wallet_") && strings.HasSuffix(fileName, ".yaml") {
			wallet, loadErr := LoadFile(filePath)
			wallets = append(wallets, wallet)
			if loadErr != nil {
				log.Errorf("while loading wallet '%s': %s", filePath, loadErr.Err)
				continue
			}

		}
	}

	return wallets, nil
}

// Load loads the wallet that match the given name in the working directory
// Note: `wallet_` prefix and a `.yaml` extension are automatically added.
func Load(nickname string) (*Wallet, error) {
	if len(nickname) == 0 {
		return nil, fmt.Errorf("nickname is required")
	}

	filePath, err := FilePath(nickname)
	if err != nil {
		return nil, fmt.Errorf("getting file path for '%s': %w", nickname, err)
	}

	if _, err := os.Stat(filePath); err != nil {
		return nil, ErrorAccountNotFound(nickname)
	}

	wallet, loadErr := LoadFile(filePath)
	if loadErr != nil {
		return nil, loadErr.Err
	}

	return &wallet, nil
}

func LoadFile(filePath string) (Wallet, *WalletError) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return Wallet{}, &WalletError{fmt.Errorf("reading file '%s': %w", filePath, err), utils.ErrAccountFile}
	}

	accountSerialized := AccountSerialized{}

	err = yaml.Unmarshal(content, &accountSerialized)
	if err != nil {
		return Wallet{}, &WalletError{fmt.Errorf("unmarshalling file '%s': %w", filePath, err), utils.ErrAccountFile}
	}
	walletHasNickname := len(accountSerialized.Nickname) > 0
	nicknameFromFileName := NicknameFromFilePath(filePath)
	baseAccount := Wallet{
		Nickname: nicknameFromFileName,
	}
	if walletHasNickname {
		baseAccount.Nickname = accountSerialized.Nickname
	}
	errMissingFields := checkMandatoryFields(accountSerialized)
	if errMissingFields != nil {
		return baseAccount, errMissingFields
	}

	account, err := accountSerialized.toAccount()
	if err != nil {
		return baseAccount, &WalletError{fmt.Errorf("deserializing account '%s': %w", filePath, err), utils.ErrAccountFile}
	}
	account.Status = StatusOK

	return account, nil
}

func checkMandatoryFields(accountSerialized AccountSerialized) *WalletError {
	if len(accountSerialized.Salt) == 0 {
		return &WalletError{fmt.Errorf("missing salt"), utils.ErrInvalidFileFormat}
	}

	if len(accountSerialized.Nonce) == 0 {
		return &WalletError{fmt.Errorf("missing nonce"), utils.ErrInvalidFileFormat}
	}

	if len(accountSerialized.CipheredData) == 0 {
		return &WalletError{fmt.Errorf("missing ciphered data"), utils.ErrInvalidFileFormat}
	}

	if len(accountSerialized.PublicKey) == 0 {
		return &WalletError{fmt.Errorf("missing public key"), utils.ErrInvalidFileFormat}
	}

	if accountSerialized.Version == nil {
		return &WalletError{fmt.Errorf("missing version"), utils.ErrInvalidFileFormat}
	}

	return nil
}

// Generate instantiates a new wallet, protects its private key and persists it.
// Everything is dynamically generated except from the nickname.
func Generate(nickname string, password string) (*Wallet, *WalletError) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, &WalletError{fmt.Errorf("generating ed25519 keypair: %w", err), utils.ErrUnknown}
	}

	wallet, createErr := createAccountFromKeys(nickname, VersionedKey(privateKey), VersionedKey(publicKey), password)
	if createErr != nil {
		return nil, createErr
	}

	err = wallet.Persist()
	if err != nil {
		return nil, &WalletError{fmt.Errorf("persisting the new wallet: %w", err), utils.ErrAccountFile}
	}

	return wallet, nil
}

// Delete removes wallet from file system
func (w *Wallet) DeleteFile() (err error) {
	return DeleteAccount(w.Nickname)
}

func DeleteAccount(nickname string) error {
	filePath, err := FilePath(nickname)
	if err != nil {
		return fmt.Errorf("getting file path for '%s': %w", nickname, err)
	}

	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("deleting wallet file '%s': %w", filePath, err)
	}

	return nil
}

func NicknameFromFilePath(filePath string) string {
	_, nicknameFromFileName := filepath.Split(filePath)
	nicknameFromFileName = strings.TrimPrefix(nicknameFromFileName, "wallet_")
	return strings.TrimSuffix(nicknameFromFileName, ".yaml")
}

// filename returns the wallet filename based on the given nickname.
func Filename(nickname string) string {
	return fmt.Sprintf("wallet_%s.yaml", nickname)
}

// FilePath returns the wallet file path base on the given nickname.
// Files are stored in
func FilePath(nickname string) (string, error) {
	walletDir, err := AccountPath()
	if err != nil {
		return "", fmt.Errorf("getting wallet directory: %w", err)
	}

	return path.Join(walletDir, Filename(nickname)), nil
}

// Filename returns the wallet filename.
func (w *Wallet) Filename() string {
	return Filename(w.Nickname)
}

// FilePath returns the wallet file path base.
func (w *Wallet) FilePath() (string, error) {
	return FilePath(w.Nickname)
}

// Import instantiates a new wallet from a private key, protects with the given password, and persists it.
func Import(nickname string, privateKeyB58V string, password string) (*Wallet, *WalletError) {
	if len(privateKeyB58V) < 2 {
		return nil, &WalletError{fmt.Errorf("invalid private key"), utils.ErrInvalidPrivateKey}
	}

	seed, version, err := base58.CheckDecode(privateKeyB58V[1:]) // omit the first byte because it's 'S' for secret key
	if err != nil {
		return nil, &WalletError{fmt.Errorf("decoding private key: %w", err), utils.ErrInvalidPrivateKey}
	}
	if !VersionIsKnown(version, []byte{PrivKeyVersion}) {
		return nil, &WalletError{fmt.Errorf("unknown private key version: %d", version), utils.ErrInvalidPrivateKey}
	}

	// The ed25519 seed is in fact what we call a private key in cryptography.
	privateKey := ed25519.NewKeyFromSeed(seed)

	pubKeyBytes := reflect.ValueOf(privateKey.Public()).Bytes() // force conversion to byte array

	wallet, createErr := createAccountFromKeys(nickname, VersionedKey(privateKey), pubKeyBytes, password)
	if createErr != nil {
		return nil, &WalletError{fmt.Errorf("creating account: %w", createErr.Err), createErr.CodeErr}
	}

	err = wallet.Persist()
	if err != nil {
		return nil, &WalletError{fmt.Errorf("persisting the new account: %w", err), utils.ErrAccountFile}
	}

	return wallet, nil
}

// createAccountFromKeys creates a new account from a private key and a public key.
// It add the versions to the keys.
// It protects the private key with the given password.
func createAccountFromKeys(nickname string, privateKey, publicKey VersionedKey, password string) (*Wallet, *WalletError) {
	var salt [16]byte
	_, err := rand.Read(salt[:])
	if err != nil {
		return nil, &WalletError{fmt.Errorf("generating random salt: %w", err), utils.ErrUnknown}
	}

	var nonce [12]byte
	_, err = rand.Read(nonce[:])
	if err != nil {
		return nil, &WalletError{fmt.Errorf("generating random nonce: %w", err), utils.ErrUnknown}
	}

	// Validate nickname
	if !NicknameIsValid(nickname) {
		return nil, &WalletError{fmt.Errorf("invalid nickname"), utils.ErrInvalidNickname}
	}

	address := addressFromPublicKey(publicKey)

	// Validate unique private key
	err = AddressIsUnique(address)
	if err != nil {
		return nil, &WalletError{err, utils.ErrDuplicateKey}
	}

	// Validate nickname uniqueness
	err = NicknameIsUnique(nickname)
	if err != nil {
		return nil, &WalletError{err, utils.ErrDuplicateNickname}
	}

	wallet := Wallet{
		Version:  AccountVersion,
		Nickname: nickname,
		Address:  address,
		KeyPair: KeyPair{
			PrivateKey: privateKey.AddVersion(PrivKeyVersion),
			PublicKey:  publicKey.AddVersion(PubKeyVersion),
			Salt:       salt,
			Nonce:      nonce,
		},
	}

	err = wallet.Protect(password)
	if err != nil {
		return nil, &WalletError{fmt.Errorf("protecting the new wallet: %w", err), utils.ErrUnknown}
	}

	return &wallet, nil
}

// Validation functions

func NicknameIsUnique(nickname string) error {
	// Load all accounts
	wallets, err := LoadAll()
	if err != nil {
		return fmt.Errorf("loading wallets: %w", err)
	}

	// Check if nickname is unique
	for _, wallet := range wallets {
		if strings.EqualFold(wallet.Nickname, nickname) {
			return fmt.Errorf("This account name already exists")
		}
	}

	return nil
}

// NicknameIsValid validates the nickname using the following rules:
// - must have at least 1 character
// - must contain only alphanumeric characters, underscores and dashes
// - must not exceed MaxNicknameLength characters
func NicknameIsValid(nickname string) bool {
	return CheckAlphanumeric(nickname) && len(nickname) <= MaxNicknameLength
}

func CheckAlphanumeric(str string) bool {
	regex := regexp.MustCompile("^[a-zA-Z0-9-_]+$")
	return regex.MatchString(str)
}

func AddressIsUnique(address string) error {
	wallets, err := LoadAll()
	if err != nil {
		return fmt.Errorf("loading accounts: %w", err)
	}

	if slices.IndexFunc(
		wallets,
		func(w Wallet) bool { return w.Address == address },
	) != -1 {
		return fmt.Errorf("importing new account: duplicate account with different name (but same keys).")
	}

	return nil
}

// Helpers

// GetPupKey returns the public key of the wallet.
func (wallet *Wallet) GetPupKey() string {
	return PublicKeyPrefix + base58.CheckEncode(wallet.KeyPair.PublicKey.RemoveVersion(), PubKeyVersion)
}

// GetPrivKey returns the versioned string representation of private key of the wallet.
// This function requires that the private key is not protected.
func (wallet *Wallet) GetPrivKey() string {
	seed := ed25519.PrivateKey(wallet.KeyPair.PrivateKey.RemoveVersion()).Seed()
	return PrivateKeyPrefix + base58.CheckEncode(seed, PrivKeyVersion)
}

// GetSalt returns the versioned string representation of the salt of the wallet.
func (wallet *Wallet) GetSalt() string {
	return base58.CheckEncode(wallet.KeyPair.Salt[:], Base58Version)
}

// GetNonce returns the versioned string representation of the nonce of the wallet.
func (wallet *Wallet) GetNonce() string {
	return base58.CheckEncode(wallet.KeyPair.Nonce[:], Base58Version)
}

// addressFromPublicKey returns the versioned string representation of the address of the wallet.
func addressFromPublicKey(pubKeyBytes VersionedKey) string {
	addr := blake3.Sum256(pubKeyBytes.AddVersion(PubKeyVersion))
	return UserAddressPrefix + base58.CheckEncode(addr[:], AddressVersion)
}

// Sign signs the given data with the wallet.
// To sign an operation, set operation to true. The operation is a base64 encoded string.
// To sign a message, set operation to false. The message is a byte array.
// This function requires that the private key is not protected.
func (wallet *Wallet) Sign(operation bool, data []byte) []byte {
	privKey := wallet.KeyPair.PrivateKey

	var digest [32]byte
	if operation {
		digest = blake3.Sum256(append(wallet.KeyPair.PublicKey, data...))
	} else {
		digest = blake3.Sum256(data)
	}

	signature := append([]byte{SignatureVersion}, ed25519.Sign(privKey.RemoveVersion(), digest[:])...)

	return signature
}
