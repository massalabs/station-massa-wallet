package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	chainID := uint64(1)
	data := []byte("test")

	expected := []byte{0, 0, 0, 0, 0, 0, 0, 1, 116, 101, 115, 116}
	actual := PrepareSignData(chainID, data)

	assert.Equal(t, len(expected), len(actual))
	assert.Equal(t, expected, actual)
}
