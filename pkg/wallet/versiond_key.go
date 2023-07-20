package wallet

type VersionedKey []byte

// AddVersion adds the version byte to the key.
func (v VersionedKey) AddVersion(version byte) []byte {
	return append([]byte{version}, v...)
}

// RemoveVersion removes the version byte from the key.
func (v VersionedKey) RemoveVersion() []byte {
	return v[1:]
}
