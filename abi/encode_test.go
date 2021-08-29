package abi

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	ZeroBytes32 = [32]byte{}
)

func Test_encodeInt(t *testing.T) {
	tests := []struct {
		typename string
		value    interface{}
		succeed  bool
		expected []byte
	}{
		{
			typename: "uint32",
			value:    255,
			succeed:  true,
			expected: padLeft([]byte{0xff}, 32, 0x0),
		},
		{
			typename: "uint64",
			value:    big.NewInt(65535),
			succeed:  true,
			expected: padLeft([]byte{0xff, 0xff}, 32, 0x0),
		},
		{
			typename: "uint256",
			value:    -255,
			succeed:  false,
			expected: nil,
		},
		{
			typename: "int32",
			value:    255,
			succeed:  true,
			expected: padLeft([]byte{0xff}, 32, 0x0),
		},
		{
			typename: "int64",
			value:    big.NewInt(65535),
			succeed:  true,
			expected: padLeft([]byte{0xff, 0xff}, 32, 0x0),
		},
		{
			typename: "int256",
			value:    -128,
			succeed:  true,
			expected: padLeft([]byte{0x80}, 32, 0xff),
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("encode %+v to %s", tt.value, tt.typename)
		t.Run(testname, func(t *testing.T) {
			res, err := encodeInt(tt.typename, tt.value)
			if tt.succeed {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, res)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
