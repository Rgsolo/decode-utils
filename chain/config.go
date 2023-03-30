package chain

type Config struct {
	ChainName string // chain name
	ChainID   int64  // chainID
	RpcURL    string // rpc url
}

var CMap = map[int64]*Config{
	1:          {ChainName: "ETH", ChainID: 1, RpcURL: "https://rpc.ankr.com/eth"},
	56:         {ChainName: "BSC", ChainID: 56, RpcURL: "https://rpc.ankr.com/bsc"},
	137:        {ChainName: "POLYGON", ChainID: 137, RpcURL: "https://rpc.ankr.com/polygon"},
	42161:      {ChainName: "Arbitrum", ChainID: 42161, RpcURL: "https://rpc.ankr.com/arbitrum"},
	43114:      {ChainName: "Avalanche", ChainID: 43114, RpcURL: "https://rpc.ankr.com/avalanche"},
	10:         {ChainName: "Optimism", ChainID: 10, RpcURL: "https://rpc.ankr.com/optimism"},
	250:        {ChainName: "Fantom", ChainID: 250, RpcURL: "https://rpc.ankr.com/fantom"},
	42220:      {ChainName: "Celo", ChainID: 42220, RpcURL: "https://rpc.ankr.com/celo"},
	1284:       {ChainName: "Moonbeam", ChainID: 1284, RpcURL: "https://rpc.ankr.com/moonbeam"},
	592:        {ChainName: "Astar", ChainID: 592, RpcURL: "https://evm.astar.network"},
	1370:       {ChainName: "Rama", ChainID: 1370, RpcURL: "https://blockchain.ramestta.com"},
	14:         {ChainName: "FLR", ChainID: 14, RpcURL: "https://flare-api.flare.network/ext/C/rpc"},
	1501795822: {ChainName: "APEX", ChainID: 1501795822, RpcURL: "https://rpc.theapexchain.org"},
	5500:       {ChainName: "GODE", ChainID: 5500, RpcURL: "https://rpc.godechain.net"},
	256256:     {ChainName: "CMP", ChainID: 256256, RpcURL: "https://mainnet.block.caduceus.foundation"},
}

func GetConfigByChainID(chainID int64) *Config {
	return CMap[chainID]
}
