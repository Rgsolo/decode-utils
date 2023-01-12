package chain

import (
	"strings"
)

type Config struct {
	ChainID int64  // chainID
	RpcURL  string // rpc url
}

var (
	ETH = Config{
		ChainID: 1,
		RpcURL:  "https://rpc.ankr.com/eth",
	}
	BSC = Config{
		ChainID: 56,
		RpcURL:  "https://rpc.ankr.com/bsc",
	}

	POLYGON = Config{
		ChainID: 137,
		RpcURL:  "https://rpc.ankr.com/polygon",
	}
	Arbitrum = Config{
		ChainID: 42161,
		RpcURL:  "https://rpc.ankr.com/arbitrum",
	}
	Avalanche = Config{
		ChainID: 43114,
		RpcURL:  "https://rpc.ankr.com/avalanche",
	}
	Optimism = Config{
		ChainID: 10,
		RpcURL:  "https://rpc.ankr.com/optimism",
	}
	Fantom = Config{
		ChainID: 250,
		RpcURL:  "https://rpc.ankr.com/fantom",
	}
	Celo = Config{
		ChainID: 42220,
		RpcURL:  "https://rpc.ankr.com/celo",
	}
	Moonbeam = Config{
		ChainID: 1284,
		RpcURL:  "https://rpc.ankr.com/moonbeam",
	}
	Astar = Config{
		ChainID: 592,
		RpcURL:  "https://evm.astar.network",
	}
	Rama = Config{
		ChainID: 1370,
		RpcURL:  "https://blockchain.ramestta.com",
	}
)

func NewConfig(name string) *Config {
	var c Config
	switch strings.ToLower(name) {
	case "eth", "ethereum":
		c = ETH
	case "bsc":
		c = BSC
	case "polygon", "pol":
		c = POLYGON
	case "arbitrum", "arb":
		c = Arbitrum
	case "avax", "avalanche":
		c = Avalanche
	case "op", "optimism":
		c = Optimism
	case "ftm", "fantom":
		c = Fantom
	case "celo":
		c = Celo
	case "moonbeam":
		c = Moonbeam
	case "astar":
		c = Astar
	case "rama":
		c = Rama

	default:
		panic("🙁not support chain")
	}
	return &c
}
