package token

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

//go:embed erc20.abi.json
var erc20Abi []byte

//go:embed erc721.abi.json
var erc721Abi []byte

//go:embed erc1155.abi.json
var erc1155Abi []byte

type DecodedCallData struct {
	Signature string
	Name      string
	Inputs    []DecodedArgument
}

type DecodedArgument struct {
	SolType abi.Argument
	Value   interface{}
}

var (
	Erc20   = Init(erc20Abi)
	Erc721  = Init(erc721Abi)
	Erc1155 = Init(erc1155Abi)
)

func Init(bytesJson []byte) *abi.ABI {
	json, err := abi.JSON(bytes.NewReader(bytesJson))
	if err != nil {
		panic(err)
	}
	return &json
}

func ParseCallData(input []byte, abispec *abi.ABI) (*DecodedCallData, error) {
	// Validate the call data that it has the 4byte prefix and the rest divisible by 32 bytes
	if len(input) < 4 {
		return nil, fmt.Errorf("invalid call data, incomplete method signature (%d bytes < 4)", len(input))
	}
	argdata := input[4:]
	if len(argdata)%32 != 0 {
		return nil, fmt.Errorf("invalid call data; length should be a multiple of 32 bytes (was %d)", len(argdata))
	}
	// Validate the called method and upack the call data accordingly
	method, err := abispec.MethodById(input)
	if err != nil {
		return nil, err
	}
	values, err := method.Inputs.UnpackValues(argdata)
	if err != nil {
		return nil, fmt.Errorf("signature %q matches, but arguments mismatch: %v", method.String(), err)
	}

	// Everything valid, assemble the call infos for the signer
	decoded := DecodedCallData{Signature: method.Sig, Name: method.RawName}
	for i := 0; i < len(method.Inputs); i++ {
		decoded.Inputs = append(decoded.Inputs, DecodedArgument{
			SolType: method.Inputs[i],
			Value:   values[i],
		})
	}
	return &decoded, nil
}
