package wallet

import "fmt"

type VersionedKey []byte

var ErrUnknownVersion = fmt.Errorf("invalid version byte")

// AddVersion adds the version byte to the key.
func (v VersionedKey) AddVersion(version byte) []byte {
	return append([]byte{version}, v...)
}

// RemoveVersion removes the version byte from the key.
func (v VersionedKey) RemoveVersion() []byte {
	return v[1:]
}

// Version returns the version byte of the key.
func (v VersionedKey) Version() byte {
	return v[0]
}

// CheckVersion checks if the version byte is in `knownVersions`.
// Return an error if the version byte is unknown.
// Return the key with the version byte otherwise.
func (v VersionedKey) CheckVersion(knownVersions []byte) ([]byte, error) {
	if VersionIsKnown(v.Version(), knownVersions) {
		return v, nil
	}

	return nil, ErrUnknownVersion
}

// VersionIsKnown checks if the `version` byte is in `knownVersions`.
func VersionIsKnown(version byte, knownVersions []byte) bool {
	for _, v := range knownVersions {
		if version == v {
			return true
		}
	}

	return false
}
