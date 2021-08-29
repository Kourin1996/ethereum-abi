package abi

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

func CalculateSelector(f *Function) ([4]byte, error) {
	arguments := make([]string, len(f.Inputs))
	for idx, input := range f.Inputs {
		arguments[idx] = input.Type
	}
	key := fmt.Sprintf("%s(%s)", f.Name, strings.Join(arguments, ","))
	hash := crypto.Keccak256([]byte(key))

	var res [4]byte
	copy(res[:], hash[:4])
	return res, nil
}

var (
	RegexInt   = regexp.MustCompile("(u?int)([1-9][0-9]*)?")
	RegexArray = regexp.MustCompile("(\\w+)\\[([1-9][0-9]*)?\\]")
)

func IsDynamic(typename string) (bool, error) {
	if f, ok := isDynamicType[typename]; ok {
		return f, nil
	}

	// array
	_, num, ok := ParseArrayType(typename)
	if ok {
		return num == "", nil
	}

	return false, fmt.Errorf("wrong type %s", typename)
}

func GetKind(typename string) (ValueKind, error) {
	if k, ok := typeToKind[typename]; ok {
		return k, nil
	}
	_, _, ok := ParseArrayType(typename)
	if ok {
		return KindArray, nil
	}
	return KindWrong, fmt.Errorf("wrong type %s", typename)
}

func ParseIntType(typename string) (uint, bool, bool) {
	res := RegexInt.FindStringSubmatch(typename)
	if len(res) == 3 {
		signed := res[1] == "int"
		if len(res[2]) != 0 {
			i, err := strconv.ParseUint(res[2], 10, 32)
			if err != nil {
				panic(err)
			}
			return uint(i), signed, true
		} else {
			return 256, signed, true
		}
	}
	return 0, false, false
}

func ParseArrayType(typename string) (string, string, bool) {
	res := RegexArray.FindStringSubmatch(typename)
	if len(res) == 3 {
		return res[1], res[2], true
	}
	return "", "", false
}

func IsArray(typename string) bool {
	_, _, ok := ParseArrayType(typename)
	return ok
}

func IsFixedArray(typename string) bool {
	_, num, ok := ParseArrayType(typename)
	return ok && num != ""
}

func IsDynamicArray(typename string) bool {
	_, num, ok := ParseArrayType(typename)
	return ok && num == ""
}
