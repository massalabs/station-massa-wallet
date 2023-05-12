package utils

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MasToNano(t *testing.T) {
	tests := []struct {
		input string
		want  uint64
	}{
		{"0", uint64(0)},
		{"18446744073.709551615", math.MaxUint64},
		{"0.123456789", uint64(123456789)},
		{"123456789", uint64(123456789000000000)},
		{"123456789.123456789", uint64(123456789123456789)},
	}

	for _, test := range tests {
		nano, err := MasToNano(test.input)
		assert.NoError(t, err)
		assert.Equal(t, test.want, nano)

		mas, err := NanoToMas(nano)
		assert.NoError(t, err)
		assert.Equal(t, test.input, mas)
	}
}
