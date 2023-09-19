// Package object provides functionalities for handling different types of Massa objects.
// It supports encoding and decoding operations with appropriate prefixes and versioning.
package object

import (
	"testing"
)

func TestNewKind(t *testing.T) {
	tests := []struct {
		name string
		args string
		want Kind
	}{
		{"Fake EncryptedPrivateKey", "S1ABDCDEF", EncryptedPrivateKey},
		{"Fake PublicKey", "P1ABDCDEF", PublicKey},
		{"Fake UserAddress", "AU1ABDCDEF", UserAddress},
		{"Fake SmartContractAddress", "AS1ABDCDEF", SmartContractAddress},
		{"Only one char", "S", Unknown},
		{"Two chars", "AS", SmartContractAddress},
		{"Unknown", "Z1ABDCDEF", Unknown},
		{"Empty string", "", Unknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewKind(tt.args); got != tt.want {
				t.Errorf("NewKind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKind_Prefix(t *testing.T) {
	tests := []struct {
		name string
		k    Kind
		want string
	}{
		{"EncryptedPrivateKey", EncryptedPrivateKey, "S"},
		{"PublicKey", PublicKey, "P"},
		{"UserAddress", UserAddress, "AU"},
		{"SmartContractAddress", SmartContractAddress, "AS"},
		{"Unknown", Unknown, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.k.Prefix(); got != tt.want {
				t.Errorf("Kind.Prefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
