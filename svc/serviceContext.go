package svc

import (
	"decode-utils/chain"
	"decode-utils/token"
	_ "embed"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type ServiceContext struct {
	Erc20     *abi.ABI
	Erc721    *abi.ABI
	Erc1155   *abi.ABI
	RpcClient *chain.EthClient
	ChainID   int64
}

func NewServiceContext(chainName string) *ServiceContext {
	config := chain.NewConfig(chainName)
	return &ServiceContext{
		Erc20:     token.Erc20,
		Erc721:    token.Erc721,
		Erc1155:   token.Erc1155,
		ChainID:   config.ChainID,
		RpcClient: chain.NewClient(config.RpcURL),
	}
}
