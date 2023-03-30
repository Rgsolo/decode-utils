package main

import (
	"context"
	"decode-utils/chain"
	"decode-utils/token"
	"flag"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
)

func main() {
	// 解析命令行参数
	flag.Parse()

	args := flag.Args()

	isSend := false
	var signedTx string
	if args[0] == "send" {
		isSend = true
		signedTx = args[1]
	} else {
		signedTx = args[0]
	}
	decode, err := hexutil.Decode(signedTx)
	if err != nil {
		panic(err)
	}
	transaction := new(types.Transaction)
	err = transaction.UnmarshalBinary(decode)
	if err != nil {
		panic(err)
	}

	getChain := chain.GetChain(transaction.ChainId().Int64())

	result, _ := transaction.MarshalJSON()
	fmt.Println()
	fmt.Println("############################ 🤡result ###############################")
	fmt.Println(string(result))
	fmt.Println()
	fmt.Println("############################ 🤡transaction details ###############################")
	fmt.Println("🌱chain name: ", getChain.Name)
	fmt.Println("🌱chain id: ", getChain.ChainID)
	fmt.Println("🌱chain rpc url: ", getChain.RpcURL[0])
	fmt.Println("🌱nonce: ", transaction.Nonce())
	fmt.Println("🌱hash: ", transaction.Hash())
	gasLimit := decimal.NewFromInt(int64(transaction.Gas()))
	fmt.Println("🌱gasLimit: ", gasLimit)
	fee := decimal.Zero
	var sender common.Address
	if transaction.Type() == types.LegacyTxType {
		fmt.Println("🌱gasPrice: ", transaction.GasPrice().String())
		fee = decimal.NewFromBigInt(transaction.GasPrice(), -18).Mul(gasLimit)

		sender, err = types.NewEIP155Signer(big.NewInt(getChain.ChainID)).Sender(transaction)
		if err != nil {
			panic(err)
		}
	} else {
		fee = decimal.NewFromBigInt(transaction.GasFeeCap(), -18).Mul(gasLimit)
		fmt.Println("🌱maxPriorityFeePerGas: ", transaction.GasTipCap().String())
		fmt.Println("🌱maxFeePerGas: ", transaction.GasFeeCap().String())

		sender, err = types.NewLondonSigner(big.NewInt(getChain.ChainID)).Sender(transaction)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("🌱fee: ", fee)
	value := decimal.NewFromBigInt(transaction.Value(), -18)
	fmt.Printf("🌱value: %s\n", value)
	fmt.Println()

	if len(transaction.Data()) != 0 {
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
	}

	fmt.Println()
	fmt.Println("############################ 🤡sender information ###############################")

	fmt.Println("🤡sender: ", sender.Hex())

	client := chain.NewClient(getChain.RpcURL[0])
	nextNonce, err := client.GetNonce(sender)
	if err != nil {
		panic(err)
	}
	fmt.Printf("🤡next nonce : %d\n", nextNonce)

	balanceAt, err := client.Client.BalanceAt(context.Background(), sender, nil)
	if err != nil {
		panic(err)
	}
	balance := decimal.NewFromBigInt(balanceAt, -18)
	fmt.Printf("🤡balance : %s\n", balance.String())

	if balance.Cmp(fee.Add(value)) >= 0 {
		fmt.Print("♓ balance is enough~ \n")
	} else {
		fmt.Printf("♓ balance is not enough: %s\n", balance.Sub(fee.Add(value)))
	}

	if isSend {
		err = client.SendTx(transaction)
		if err != nil {
			panic(err)
		}
		fmt.Println("🤡send success")
	}
}
