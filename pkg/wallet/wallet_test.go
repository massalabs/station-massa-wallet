package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkProtectWallet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		iteration()
	}
}

func iteration() {
	password := "password"
	wlt, errGenerate := Generate("nickname", password)
	if errGenerate != nil {
		panic(errGenerate.Err)
	}
	err := wlt.Protect(password)
	if err != nil {
		panic(err)
	}
}

func TestExtractNickname(t *testing.T) {
	tests := []struct {
		filePath string
		expected string
	}{
		{"/home/moi/wallet_{{nickname}}.yaml", "{{nickname}}"},
		{"/home/user/wallet_example.yml", "example"},
		{"/path/to/wallet_abc123.yaml", "abc123"},
		{"/another/directory/wallet_unit_test.yaml", "unit_test"},
		{"wallet_nickname.yaml", "nickname"},
		{"/wallet_123", "123"},
		{"/no_prefix.yaml", ""},
	}

	for _, test := range tests {
		nickname := ExtractNickname(test.filePath)
		assert.Equal(t, test.expected, nickname)
	}
}
