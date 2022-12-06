package main

import (
	"decode-utils/svc"
	"decode-utils/token"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"math/big"
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Println("🙁err args")
		return
	}

	svcCtx := svc.NewServiceContext(args[0])

	decode, err := hexutil.Decode(args[1])
	if err != nil {
		panic(err)
	}
	transaction := new(types.Transaction)
	err = transaction.UnmarshalBinary(decode)
	if err != nil {
		panic(err)
	}
	result, _ := transaction.MarshalJSON()
	fmt.Println()
	fmt.Println("############################ 🤡result ###############################")
	fmt.Println(string(result))
	fmt.Println()
	fmt.Println("🌱nonce: ", transaction.Nonce())
	fmt.Println("🌱hash: ", transaction.Hash())
	fmt.Println("🌱gasLimit: ", transaction.Gas())
	if transaction.Type() == types.LegacyTxType {
		fmt.Println("🌱gasPrice: ", transaction.GasPrice().String())
	} else {
		fmt.Println("🌱maxPriorityFeePerGas: ", transaction.GasTipCap().String())
		fmt.Println("🌱maxPriorityFeePerGas: ", transaction.GasFeeCap().String())
	}

	fmt.Println()
	data, err := token.ParseCallData(transaction.Data(), token.Erc20)
	if err != nil {
		data, err = token.ParseCallData(transaction.Data(), token.Erc721)
		if err != nil {
			data, err = token.ParseCallData(transaction.Data(), token.Erc1155)
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

	sender, err := types.NewEIP155Signer(big.NewInt(svcCtx.ChainID)).Sender(transaction)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println("🤡sender: ", sender.Hex())

	nonce, err := svcCtx.RpcClient.GetNonce(sender)
	if err != nil {
		panic(err)
	}
	fmt.Printf("🤡next nonce : %d\n", nonce)
}
