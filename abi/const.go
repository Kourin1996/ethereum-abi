package abi

import "fmt"

var isDynamicType = map[string]bool{
	// static
	"address": false,
	"bool":    false,

	// dynamic
	"bytes":  true,
	"string": true,
}

type ValueKind string

const (
	KindWrong   ValueKind = "wrong_kind"
	KindUint    ValueKind = "uint_kind"
	KindInt     ValueKind = "int_kind"
	KindFixed   ValueKind = "fixed_kind"
	KindUFixed  ValueKind = "ufixed_kind"
	KindBool    ValueKind = "bool_kind"
	KindAddress ValueKind = "address_kind"
	KindBytes   ValueKind = "bytes_kind"
	KindString  ValueKind = "string_kind"
	KindArray   ValueKind = "array_kind"
)

var typeToKind = map[string]ValueKind{
	"bool":    KindBool,
	"address": KindAddress,
	"bytes":   KindBytes,
	"string":  KindString,
}

func init() {
	initIsDynamic()
	initTypeToKind()
}

func initIsDynamic() {
	// integer & fixed
	isDynamicType["uint"] = false
	isDynamicType["int"] = false
	isDynamicType["fixed"] = false
	isDynamicType["ufixed"] = false
	for i := 1; i <= 256/8; i++ {
		b := i * 8 // bits length
		isDynamicType[fmt.Sprintf("uint%d", b)] = false
		isDynamicType[fmt.Sprintf("int%d", b)] = false
		for k := 1; k <= 80; k++ {
			isDynamicType[fmt.Sprintf("fixed%dx%d", b, k)] = false
			isDynamicType[fmt.Sprintf("ufixed%dx%d", b, k)] = false
		}
	}
	// bytes
	for i := 1; i <= 32; i++ {
		isDynamicType[fmt.Sprintf("bytes%d", i)] = false
	}
}

func initTypeToKind() {
	// integer & fixed
	typeToKind["uint"] = KindUint
	typeToKind["int"] = KindInt
	typeToKind["fixed"] = KindFixed
	typeToKind["ufixed"] = KindUFixed
	for i := 1; i <= 256/8; i++ {
		b := i * 8 // bits length
		typeToKind[fmt.Sprintf("uint%d", b)] = KindUint
		typeToKind[fmt.Sprintf("int%d", b)] = KindInt
		for k := 1; k <= 80; k++ {
			typeToKind[fmt.Sprintf("fixed%dx%d", b, k)] = KindFixed
			typeToKind[fmt.Sprintf("ufixed%dx%d", b, k)] = KindUFixed
		}
	}
	// bytes
	for i := 1; i <= 32; i++ {
		typeToKind[fmt.Sprintf("bytes%d", i)] = KindBytes
	}
}
