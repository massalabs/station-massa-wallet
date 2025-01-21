package cache

import (
	"crypto/ed25519"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/awnumar/memguard"
	"github.com/bluele/gcache"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"lukechampine.com/blake3"
)

var (
	cache gcache.Cache
	once  sync.Once
)

func Init() gcache.Cache {
	once.Do(func() {
		cache = gcache.New(cacheSize).
			LRU().
			Build()
	})

	return cache
}

const (
	defaultExpirationTime = time.Hour * 24 * 30
	cacheKeyPrefix        = "pkey"
	notFoundMsg           = "private key not found in cache"
	cacheSize             = 100
)

func CachePrivateKeyFromPassword(account *account.Account, password *memguard.LockedBuffer) error {
	privateKey, err := account.PrivateKeyBytesInClear(password)
	if err != nil {
		return fmt.Errorf("error caching private key: %w", err)
	}

	err = CachePrivateKey(account, privateKey)
	if err != nil {
		return fmt.Errorf("error caching private key: %w", err)
	}

	return nil
}

func CachePrivateKey(account *account.Account, privateKey *memguard.LockedBuffer) error {
	cacheKey, err := privateKeyCacheKey(account)
	if err != nil {
		return fmt.Errorf("%w: %w", utils.ErrPrivateKeyCache, err)
	}

	key := KeyHash([]byte(cacheKey))

	// Concatenate the key with itself to make it 64 bytes long
	var extendedKey [64]byte
	copy(extendedKey[:], append(key[:], key[:]...))

	cipheredPrivateKey, err := xor(privateKey, extendedKey)
	if err != nil {
		return fmt.Errorf("%w: %w", utils.ErrCache, err)
	}

	cacheValue := make([]byte, ed25519.PrivateKeySize)
	copy(cacheValue, cipheredPrivateKey.Bytes())
	cipheredPrivateKey.Destroy()

	err = cache.SetWithExpire(key, cacheValue, expirationDuration())
	if err != nil {
		return fmt.Errorf("error set private key in cache: %w", err)
	}

	return nil
}

// privateKeyFromCache return the private key from the cache or an error.
func PrivateKeyFromCache(
	acc *account.Account,
) (*memguard.LockedBuffer, error) {
	cacheKey, err := privateKeyCacheKey(acc)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", utils.ErrPrivateKeyCache, err)
	}

	keyHash := KeyHash([]byte(cacheKey))

	value, err := cache.Get(keyHash)
	if err != nil {
		if err.Error() == gcache.KeyNotFoundError.Error() {
			return nil, fmt.Errorf("%w: %w", utils.ErrPrivateKeyCache, err)
		}

		return nil, fmt.Errorf("%w: %w", utils.ErrCache, err)
	}

	if value == nil {
		return nil, fmt.Errorf("%w: %s", utils.ErrCache, notFoundMsg)
	}

	byteValue, ok := value.([]byte)
	if !ok {
		return nil, fmt.Errorf("%w: %s", utils.ErrCache, "value is not a byte array")
	}

	cacheValue := make([]byte, ed25519.PrivateKeySize)
	copy(cacheValue, byteValue)

	cipheredPrivateKey := memguard.NewBufferFromBytes(cacheValue)

	// Concatenate the key with itself to make it 64 bytes long
	var extendedKey [64]byte
	copy(extendedKey[:], append(keyHash[:], keyHash[:]...))

	privateKey, err := xor(cipheredPrivateKey, extendedKey)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", utils.ErrCache, err)
	}

	cipheredPrivateKey.Destroy()

	return privateKey, nil
}

func xor(pkey *memguard.LockedBuffer, cacheKeyHash [64]byte) (*memguard.LockedBuffer, error) {
	a := pkey.Bytes()

	if len(a) != len(cacheKeyHash) {
		return nil, fmt.Errorf("length of two arrays must be same, %d and %d", len(a), len(cacheKeyHash))
	}
	result := make([]byte, len(a))

	for i := 0; i < len(a); i++ {
		result[i] = a[i] ^ cacheKeyHash[i]
	}

	return memguard.NewBufferFromBytes(result), nil
}

func privateKeyCacheKey(account *account.Account) ([]byte, error) {
	address, err := account.Address.String()
	if err != nil {
		return []byte(""), fmt.Errorf("err: %w", err)
	}

	return []byte(cacheKeyPrefix + address), nil
}

func KeyHash(cacheKeyHash []byte) [32]byte {
	return blake3.Sum256(cacheKeyHash)
}

func expirationDuration() time.Duration {
	fromEnv := os.Getenv("PKEY_CACHE_EXPIRATION_TIME")

	if fromEnv == "" {
		return defaultExpirationTime
	}

	duration, err := time.ParseDuration(fromEnv)
	if err != nil {
		return defaultExpirationTime
	}

	return duration
}
