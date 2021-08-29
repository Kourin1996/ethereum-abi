package main

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/Kourin1996/ethereum-abi/abi"
)

func main() {
	fmt.Printf("Hello, World\n")

	abiData := abi.ABI{}
	if err := json.Unmarshal([]byte(ABITest1), &abiData); err != nil {
		fmt.Println(err)
	}
	res, err := abiData.EncodeCallFunctionData("setValue", big.NewInt(-12345))
	if err != nil {
		fmt.Printf("err=%+v\n", err)
		return
	}
	fmt.Printf("res=%x\n", res)
}
