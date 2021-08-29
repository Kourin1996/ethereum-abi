package main

import (
	"fmt"
	"math/big"

	"github.com/Kourin1996/ethereum-abi/abi"
)

func main() {
	fmt.Printf("Hello, World\n")

	abiData, err := abi.Parse(ABITest1)
	if err != nil {
		fmt.Printf("err=%+v\n", err)
	}
	res, err := abiData.EncodeCallFunctionData("setValue", big.NewInt(12345))
	if err != nil {
		fmt.Printf("err=%+v\n", err)
		return
	}
	fmt.Printf("res=%x\n", res)
}
