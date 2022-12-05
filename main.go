package main

import (
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
	// ############################ 🤡input ############################### \n
	fmt.Println()
	fmt.Println("############################ 🤡result ###############################")
	fmt.Println(string(result))
	fmt.Println()
	fmt.Println("🤡nonce: ", t.Nonce())
	fmt.Println("🤡hash: ", t.Hash())
	fmt.Println("🤡gasLimit: ", t.Gas())
	if t.Type() == types.LegacyTxType {
		fmt.Println("🤡gasPrice: ", t.GasPrice().String())
	} else {
		fmt.Println("🤡maxPriorityFeePerGas: ", t.GasTipCap().String())
		fmt.Println("🤡maxPriorityFeePerGas: ", t.GasFeeCap().String())
	}

	sender, err := types.NewEIP155Signer(big.NewInt(int64(chainId))).Sender(t)
	if err != nil {
		panic(err)
	}
	fmt.Println("🤡sender: ", sender.Hex())
}
