package abi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDynamic(t *testing.T) {
	tests := []struct {
		typename string
		expected bool
	}{
		{
			typename: "address",
			expected: false,
		},
		{
			typename: "bool",
			expected: false,
		},
		{
			typename: "uint",
			expected: false,
		},
		{
			typename: "uint8",
			expected: false,
		},
		{
			typename: "uint32",
			expected: false,
		},
		{
			typename: "uint64",
			expected: false,
		},
		{
			typename: "int",
			expected: false,
		},
		{
			typename: "int8",
			expected: false,
		},
		{
			typename: "int64",
			expected: false,
		},
		{
			typename: "int128",
			expected: false,
		},
		{
			typename: "int256",
			expected: false,
		},
		{
			typename: "fixed",
			expected: false,
		},
		{
			typename: "fixed128x20",
			expected: false,
		},
		{
			typename: "ufixed",
			expected: false,
		},
		{
			typename: "fixed256x80",
			expected: false,
		},
		{
			typename: "address[5]",
			expected: false,
		},
		{
			typename: "uint[10]",
			expected: false,
		},
		{
			typename: "int[20]",
			expected: false,
		},
		{
			typename: "bytes",
			expected: true,
		},
		{
			typename: "string",
			expected: true,
		},
		{
			typename: "address[]",
			expected: true,
		},
		{
			typename: "string[]",
			expected: true,
		},
		{
			typename: "uint8[]",
			expected: true,
		},
		{
			typename: "fixed8x10[]",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.typename, func(t *testing.T) {
			res, err := IsDynamic(tt.typename)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, res, "IsDynamic returns unexpected value")
		})
	}
}
