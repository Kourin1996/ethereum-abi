package abi

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"reflect"
)

func (abi *ABI) EncodeCallConstructorData(args ...interface{}) ([]byte, error) {
	if abi.Constructor == nil {
		return nil, errors.New("constructor is not set in ABI")
	}
	// todo: check
	return encodeCallData(abi.Constructor, args)
}

func (abi *ABI) EncodeCallFunctionData(name string, args ...interface{}) ([]byte, error) {
	f, ok := abi.Functions[name]
	if !ok {
		return nil, fmt.Errorf("function %s is not set in ABI", name)
	}
	argData, err := encodeCallData(&f, args)
	if err != nil {
		return nil, err
	}
	//todo: check
	return append(f.Selector[:], argData...), nil
}

func encodeCallData(f *Function, args []interface{}) ([]byte, error) {
	fmt.Printf("encodeCallData name=%s function=%+v, arguments=%+v\n", f.Name, f, args)
	return encodeTuple(f.Inputs, args)
}

// todo: change type
func encodeTuple(types []FunctionValue, values []interface{}) ([]byte, error) {
	if len(types) != len(values) {
		return nil, fmt.Errorf("mismatch num values in tuple, required=%d, actual=%d", len(types), len(values))
	}
	var heads, tails bytes.Buffer
	for idx, t := range types {
		head, tail, err := encodeValue(t, values[idx])
		if err != nil {
			return nil, err
		}
		heads.Write(head)
		if len(tail) > 0 {
			tails.Write(tail)
		}
	}
	return append(heads.Bytes(), tails.Bytes()...), nil
}

func encodeValue(t FunctionValue, v interface{}) ([]byte, []byte, error) {
	kind, err := GetKind(t.Type)
	if err != nil {
		return nil, nil, err
	}
	switch kind {
	case KindUint, KindInt:
		{
			head, err := encodeInt(t.Type, v)
			return head, nil, err
		}
	case KindFixed:
	case KindUFixed:
	case KindBool:

	case KindAddress:
	case KindBytes:
	case KindString:
	case KindArray:
	}
	// shouldn't reach here
	return nil, nil, fmt.Errorf("wrong kind type=%s, kind=%s", t.Type, kind)
}

func encodeInt(typename string, v interface{}) ([]byte, error) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return encodeFromBigInt(typename, new(big.Int).SetUint64(rv.Uint()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return encodeFromBigInt(typename, big.NewInt(rv.Int()))
	case reflect.TypeOf(big.Int{}).Kind():
		v := rv.Interface().(big.Int)
		return encodeFromBigInt(typename, &v)
	case reflect.Ptr:
		if rv.IsNil() {
			return nil, fmt.Errorf("value mustn't be nil, expected type = %s", typename)
		}
		return encodeInt(typename, rv.Elem().Interface())
	}
	return nil, fmt.Errorf("wrong type, expected int value but actual = %T", v)
}

var (
	mask256 = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
)

func encodeFromBigInt(typename string, v *big.Int) ([]byte, error) {
	_, signed, ok := ParseIntType(typename)
	if !ok {
		return nil, fmt.Errorf("expected int typename, but actual=%s", typename)
	}
	if !signed && v.Sign() < 0 {
		return nil, fmt.Errorf("expected unsigned int value, but actual=%s", v.String())
	}

	if v.Sign() >= 0 {
		// positive or zero
		return padLeft(v.Bytes(), 32, 0x0), nil
	} else {
		// negative
		v = v.And(v, mask256)
		// calculate two's complement encoded value
		v = v.Add(v.Not(v), big.NewInt(1))
		return padLeft(v.Bytes(), 32, 0xff), nil
	}
}

func getPadding(size int, e byte) []byte {
	padding := make([]byte, size)
	for i := 0; i < size; i++ {
		padding[i] = e
	}
	return padding
}

func padLeft(data []byte, size int, e byte) []byte {
	l := len(data)
	if l == size {
		return data
	}
	if l > size {
		return data[(l - size):]
	}
	padding := getPadding(size-l, e)
	return append(padding, data...)
}

func padRight(data []byte, size int, e byte) []byte {
	l := len(data)
	if l == size {
		return data
	}
	if l > size {
		return data[(l - size):]
	}
	padding := getPadding(size-l, e)
	return append(data, padding...)
}
