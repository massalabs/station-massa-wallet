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

func Test_NanoToMas(t *testing.T) {
	tests := []struct {
		input uint64
		want  string
	}{
		{uint64(0), "0"},
		{math.MaxUint64, "18446744073.709551615"},
		{uint64(123456789), "0.123456789"},
		{uint64(123456789000000000), "123456789"},
		{uint64(123456789123456789), "123456789.123456789"},
		{uint64(1), "0.000000001"},
		{uint64(999999999999999999), "999999999.999999999"},
	}

	for _, test := range tests {
		mas, err := NanoToMas(test.input)
		if assert.NoError(t, err, "Expected no error for NanoToMas conversion of %d", test.input) {
			assert.Equal(t, test.want, mas, "NanoToMas conversion for input: %d", test.input)
		}
	}
}
