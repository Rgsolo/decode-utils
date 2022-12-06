package chain

import (
	"strings"
)

type Config struct {
	ChainID int64  // chainID
	Name    string // lowercase name
	RpcURL  string // rpc url
}

var (
	ETH = Config{
		ChainID: 1,
		Name:    "eth",
		RpcURL:  "https://rpc.ankr.com/eth",
	}
	BSC = Config{
		ChainID: 56,
		Name:    "bsc",
		RpcURL:  "https://rpc.ankr.com/bsc",
	}

	POLYGON = Config{
		ChainID: 137,
		Name:    "polygon",
		RpcURL:  "https://rpc.ankr.com/polygon",
	}
)

func NewConfig(name string) *Config {
	var newConfig Config
	switch strings.ToLower(name) {
	case ETH.Name:
		newConfig = ETH
	case BSC.Name:
		newConfig = BSC
	case POLYGON.Name:
		newConfig = POLYGON
	default:
		panic("🙁not support chain")
	}
	return &newConfig
}
