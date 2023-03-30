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
	ChainName string
}

func NewServiceContext(chainId int64) *ServiceContext {
	config := chain.GetConfigByChainID(chainId)
	return &ServiceContext{
		Erc20:     token.Erc20,
		Erc721:    token.Erc721,
		Erc1155:   token.Erc1155,
		ChainID:   config.ChainID,
		ChainName: config.ChainName,
		RpcClient: chain.NewClient(config.RpcURL),
	}
}
