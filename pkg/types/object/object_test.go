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

func TestObject_Validate(t *testing.T) {
	// Define some sample test cases
	testCases := []struct {
		name          string
		object        Object
		lastVersion   byte
		expectedKinds []Kind
		expectedError error
	}{
		{
			name:          "ValidObject",
			object:        Object{Data: []byte{}, Kind: EncryptedPrivateKey, Version: 1},
			lastVersion:   2,
			expectedKinds: []Kind{EncryptedPrivateKey},
			expectedError: nil,
		},
		{
			name:          "UnsupportedVersion",
			object:        Object{Data: []byte{}, Kind: EncryptedPrivateKey, Version: 3},
			lastVersion:   2,
			expectedKinds: []Kind{EncryptedPrivateKey},
			expectedError: ErrUnsupportedVersion,
		},
		{
			name:          "InvalidKind",
			object:        Object{Data: []byte{}, Kind: Unknown, Version: 1},
			lastVersion:   2,
			expectedKinds: []Kind{EncryptedPrivateKey},
			expectedError: ErrInvalidType,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.object.Validate(testCase.lastVersion, testCase.expectedKinds...)
			if err != testCase.expectedError {
				t.Errorf("Expected error: %v, but got error: %v", testCase.expectedError, err)
			}
		})
	}
}
