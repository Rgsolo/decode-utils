package main

import (
	"context"
	"decode-utils/svc"
	"decode-utils/token"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
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
	fmt.Println("############################ 🤡transaction details###############################")
	fmt.Println("🌱nonce: ", transaction.Nonce())
	fmt.Println("🌱hash: ", transaction.Hash())
	gasLimit := decimal.NewFromInt(int64(transaction.Gas()))
	fmt.Println("🌱gasLimit: ", gasLimit)
	fee := decimal.Zero
	var sender common.Address
	if transaction.Type() == types.LegacyTxType {
		fmt.Println("🌱gasPrice: ", transaction.GasPrice().String())
		fee = decimal.NewFromBigInt(transaction.GasPrice(), -18).Mul(gasLimit)

		sender, err = types.NewEIP155Signer(big.NewInt(svcCtx.ChainID)).Sender(transaction)
		if err != nil {
			panic(err)
		}
	} else {
		fee = decimal.NewFromBigInt(transaction.GasFeeCap(), -18).Mul(gasLimit)
		fmt.Println("🌱maxPriorityFeePerGas: ", transaction.GasTipCap().String())
		fmt.Println("🌱maxFeePerGas: ", transaction.GasFeeCap().String())

		sender, err = types.NewLondonSigner(big.NewInt(svcCtx.ChainID)).Sender(transaction)
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
	fmt.Println("############################ 🤡sender information###############################")

	fmt.Println("🤡sender: ", sender.Hex())
	nextNonce, err := svcCtx.RpcClient.GetNonce(sender)
	if err != nil {
		panic(err)
	}
	fmt.Printf("🤡next nonce : %d\n", nextNonce)

	balanceAt, err := svcCtx.RpcClient.Client.BalanceAt(context.Background(), sender, nil)
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
}
