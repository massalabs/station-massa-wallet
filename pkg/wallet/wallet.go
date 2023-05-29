package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/go-openapi/strfmt"
	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/pkg/utils"
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
	UserAddressPrefix         = "AU"
	PublicKeyPrefix           = "P"
	PrivateKeyPrefix          = "S"
	MaxNicknameLength         = 32
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
	PrivateKey []byte
	PublicKey  []byte
	Salt       [16]byte
	Nonce      [12]byte
}

// Wallet structure allows to link a nickname, an address and a version to one or more key pairs.
type Wallet struct {
	Version  uint8
	Nickname string
	Address  string
	KeyPair  KeyPair
}

type AccountSerialized struct {
	Version      *uint8   `yaml:"Version"`
	Nickname     string   `yaml:"Nickname"`
	Address      string   `yaml:"Address"`
	Salt         [16]byte `yaml:"Salt,flow"`
	Nonce        [12]byte `yaml:"Nonce,flow"`
	CipheredData []byte   `yaml:"CipheredData,flow"`
	PublicKey    []byte   `yaml:"PublicKey,flow"`
}

func (accountSerialized *AccountSerialized) ToAccount() Wallet {
	wallet := Wallet{
		Version:  *accountSerialized.Version,
		Nickname: accountSerialized.Nickname,
		Address:  accountSerialized.Address,
		KeyPair: KeyPair{
			PrivateKey: make([]byte, 0),
			PublicKey:  accountSerialized.PublicKey,
			Salt:       accountSerialized.Salt,
			Nonce:      accountSerialized.Nonce,
		},
	}

	return wallet
}

func (account *Wallet) ToAccountSerialized() AccountSerialized {
	accountSerialized := AccountSerialized{
		Version:      &account.Version,
		Nickname:     account.Nickname,
		Address:      account.Address,
		Salt:         account.KeyPair.Salt,
		Nonce:        account.KeyPair.Nonce,
		CipheredData: make([]byte, 0),
		PublicKey:    account.KeyPair.PublicKey,
	}

	return accountSerialized
}

// aead returns a authenticated encryption with associated data.
func aead(password []byte, salt []byte) (cipher.AEAD, error) {
	secretKey := pbkdf2.Key([]byte(password), salt, PBKDF2NbRound, SecretKeyLength, sha256.New)

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, fmt.Errorf("intializing block ciphering: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("intializing the AES block cipher wrapped in a Gallois Counter Mode ciphering: %w", err)
	}

	return aesGCM, nil
}

// Protect encrypts the private key using the given guiModal.
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

	w.KeyPair.PrivateKey = pk

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
// Note: the wallet is stored in YAML format and in Thyra config directory.
func (w *Wallet) Persist() error {
	accountSerialized := w.ToAccountSerialized()

	// account is protected so PrivateKey is encrypted
	accountSerialized.CipheredData = w.KeyPair.PrivateKey

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

func GetWorkDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("getting executable path: %w", err)
	}

	dir := filepath.Dir(ex)

	// Helpful when developing:
	// when running `go run`, the executable is in a temporary directory.
	if strings.Contains(dir, "go-build") {
		return ".", nil
	}
	return filepath.Dir(ex), nil
}

// GetWalletDir returns the path where the account yaml file are stored.
// Note: the wallet directory is the folder where the wallet plugin binary resides.
func GetWalletDir() (string, error) {
	return GetWorkDir()
}

// LoadAll loads all the wallets in the working directory.
// Note: a wallet must have: `wallet_` prefix and a `.yaml` extension.
func LoadAll() ([]Wallet, error) {
	walletDir, err := GetWalletDir()
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
			if loadErr != nil {
				return nil, loadErr.Err
			}

			wallets = append(wallets, wallet)
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

	errMissingFields := checkMandatoryFields(accountSerialized)
	if errMissingFields != nil {
		return Wallet{}, errMissingFields
	}

	account := accountSerialized.ToAccount()
	account.KeyPair.PrivateKey = accountSerialized.CipheredData

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

	wallet, createErr := createAccountFromKeys(nickname, privateKey, publicKey, password)
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
	filePath, err := FilePath(w.Nickname)
	if err != nil {
		return fmt.Errorf("getting file path for '%s': %w", w.Nickname, err)
	}

	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("deleting wallet file '%s': %w", filePath, err)
	}

	return nil
}

// filename returns the wallet filename based on the given nickname.
func Filename(nickname string) string {
	return fmt.Sprintf("wallet_%s.yaml", nickname)
}

// FilePath returns the wallet file path base on the given nickname.
// Files are stored in
func FilePath(nickname string) (string, error) {
	walletDir, err := GetWalletDir()
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

func Import(nickname string, privateKeyB58V string, password string) (*Wallet, *WalletError) {
	if len(privateKeyB58V) < 2 {
		return nil, &WalletError{fmt.Errorf("invalid private key"), utils.ErrInvalidPrivateKey}
	}

	privateKeyBytes, _, err := base58.CheckDecode(privateKeyB58V[1:])
	if err != nil {
		return nil, &WalletError{fmt.Errorf("decoding private key: %w", err), utils.ErrInvalidPrivateKey}
	}

	// The ed25519 seed is in fact what we call a private key in cryptography...
	privateKey := ed25519.NewKeyFromSeed(privateKeyBytes)

	pubKeyBytes := reflect.ValueOf(privateKey.Public()).Bytes() // force conversion to byte array

	wallet, createErr := createAccountFromKeys(nickname, privateKey, pubKeyBytes, password)
	if createErr != nil {
		return nil, &WalletError{fmt.Errorf("creating account: %w", createErr.Err), createErr.CodeErr}
	}

	err = wallet.Persist()
	if err != nil {
		return nil, &WalletError{fmt.Errorf("persisting the new account: %w", err), utils.ErrAccountFile}
	}

	return wallet, nil
}

func createAccountFromKeys(nickname string, privateKey []byte, publicKey []byte, password string) (*Wallet, *WalletError) {
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
		Version:  0,
		Nickname: nickname,
		Address:  address,
		KeyPair: KeyPair{
			PrivateKey: privateKey,
			PublicKey:  publicKey,
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

func NicknameIsUnique(nickname string) error {
	// Load all accounts
	wallets, err := LoadAll()
	if err != nil {
		return fmt.Errorf("loading wallets: %w", err)
	}

	// Check if nickname is unique
	for _, wallet := range wallets {
		if wallet.Nickname == nickname {
			return fmt.Errorf("This account name already exists")
		}
	}

	return nil
}

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

// GetPupKey returns the public key of the wallet.
func (wallet *Wallet) GetPupKey() string {
	return PublicKeyPrefix + base58.CheckEncode(wallet.KeyPair.PublicKey, Base58Version)
}

// GetPrivKey returns the private key of the wallet.
// This function requires that the private key is not protected.
func (wallet *Wallet) GetPrivKey() string {
	seed := ed25519.PrivateKey(wallet.KeyPair.PrivateKey).Seed()
	return PrivateKeyPrefix + base58.CheckEncode(seed, Base58Version)
}

func (wallet *Wallet) GetSalt() string {
	return base58.CheckEncode(wallet.KeyPair.Salt[:], Base58Version)
}

func (wallet *Wallet) GetNonce() string {
	return base58.CheckEncode(wallet.KeyPair.Nonce[:], Base58Version)
}

func addressFromPublicKey(pubKeyBytes []byte) string {
	addr := blake3.Sum256(pubKeyBytes)
	return UserAddressPrefix + base58.CheckEncode(addr[:], Base58Version)
}

// Sign signs the given operation with the wallet.
// The operation is a base64 encoded string.
func (wallet *Wallet) Sign(operation *strfmt.Base64) ([]byte, error) {
	pubKey := wallet.KeyPair.PublicKey
	privKey := wallet.KeyPair.PrivateKey

	digest, err := hash(operation, pubKey)
	if err != nil {
		return nil, err
	}

	signature := ed25519.Sign(privKey, digest[:])
	return signature, nil
}

// hash prepares the digest for signature.
func hash(operation *strfmt.Base64, publicKey []byte) ([32]byte, error) {
	op, err := base64.StdEncoding.DecodeString(operation.String())
	if err != nil {
		return [32]byte{}, fmt.Errorf("decoding operation: %w", err)
	}

	digest := blake3.Sum256(append(publicKey, op...))

	return digest, nil
}
