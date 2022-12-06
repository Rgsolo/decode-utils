package main

import (
	"decode-utils/token"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"math/big"
	"strings"
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Println("🙁err args")
		return
	}

	var chainId int
	switch strings.ToLower(args[0]) {
	case "eth":
		chainId = eth
	case "bsc":
		chainId = bsc
	case "polygon":
		chainId = polygon
	default:
		fmt.Println("🙁not support chain")
		return
	}

	decode, err := hexutil.Decode(args[1])
	if err != nil {
		panic(err)
	}
	t := new(types.Transaction)
	err = t.UnmarshalBinary(decode)
	if err != nil {
		panic(err)
	}
	result, _ := t.MarshalJSON()
	fmt.Println()
	fmt.Println("############################ 🤡result ###############################")
	fmt.Println(string(result))
	fmt.Println()
	fmt.Println("🌱nonce: ", t.Nonce())
	fmt.Println("🌱hash: ", t.Hash())
	fmt.Println("🌱gasLimit: ", t.Gas())
	if t.Type() == types.LegacyTxType {
		fmt.Println("🌱gasPrice: ", t.GasPrice().String())
	} else {
		fmt.Println("🌱maxPriorityFeePerGas: ", t.GasTipCap().String())
		fmt.Println("🌱maxPriorityFeePerGas: ", t.GasFeeCap().String())
	}

	fmt.Println()
	data, err := token.ParseCallData(t.Data(), token.Erc20)
	if err != nil {
		data, err = token.ParseCallData(t.Data(), token.Erc721)
		if err != nil {
			data, err = token.ParseCallData(t.Data(), token.Erc1155)
			if err != nil {
				fmt.Println("🙁not support contract")
			} else {
				fmt.Printf("erc1155: %s \n", data.Signature)
			}
		} else {
			fmt.Printf("erc721: %s \n", data.Signature)
		}
	} else {
		fmt.Printf("erc20: %s \n", data.Signature)
	}
	for _, input := range data.Inputs {
		fmt.Printf("🌱%s[%s]: %s \n", input.SolType.Name, input.SolType.Type, input.Value)
	}

	sender, err := types.NewEIP155Signer(big.NewInt(int64(chainId))).Sender(t)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println("🤡sender: ", sender.Hex())
}
