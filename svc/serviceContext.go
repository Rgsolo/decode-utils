package svc

import (
	_ "embed"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type ServiceContext struct {
	Erc20   *abi.ABI
	Erc721  *abi.ABI
	Erc1155 *abi.ABI
}

//
//func NewServiceContext() *ServiceContext {
//
//	return &ServiceContext{Erc20: token.Erc20, Erc721: token.Erc721, Erc1155: token.Erc1155}
//}
