package main

// Test1
/*
pragma solidity ^0.8.7;
pragma abicoder v2;

contract Test1 {
    uint32 public _value;

    constructor(uint32 val) {
        _value = val;
    }

    function setValue(uint32 val) public {
        _value = val;
    }

    function increment() public returns (uint32) {
        _value += 1;
        return _value;
    }

    function value() public view returns (uint32) {
        return _value;
    }
}
*/

const ABITest1 = `[
    {
      "inputs": [
        {
          "internalType": "uint32",
          "name": "val",
          "type": "uint32"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "constructor"
    },
    {
      "inputs": [],
      "name": "_value",
      "outputs": [
        {
          "internalType": "uint32",
          "name": "",
          "type": "uint32"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "increment",
      "outputs": [
        {
          "internalType": "uint32",
          "name": "",
          "type": "uint32"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint32",
          "name": "val",
          "type": "uint32"
        }
      ],
      "name": "setValue",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "value",
      "outputs": [
        {
          "internalType": "uint32",
          "name": "",
          "type": "uint32"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    }
  ]`
