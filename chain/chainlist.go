package chain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Chain struct {
	Name    string   `json:"name"`
	ChainID int64    `json:"chainId"`
	RpcURL  []string `json:"rpc"`
}

func (c *Chain) CheckRpcURL(url string) bool {
	reqBody := []byte(`{"jsonrpc": "2.0", "method": "eth_blockNumber", "params": [], "id": 1}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var result struct {
		JsonRPC string `json:"jsonrpc"`
		Result  string `json:"result"`
		Id      int    `json:"id"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return false
	}

	return result.JsonRPC == "2.0" && result.Result != ""
}

func GetChain(chainID int64) *Chain {
	resp, err := http.Get("https://chainid.network/chains.json")
	if err != nil {
		fmt.Println("无法获取 JSON 数据：", err)
		return nil
	}
	defer resp.Body.Close()

	var chains []Chain
	err = json.NewDecoder(resp.Body).Decode(&chains)
	if err != nil {
		fmt.Println("无法解析 JSON 数据：", err)
		return nil
	}

	for _, chain := range chains {
		if chain.ChainID != chainID {
			continue
		}

		for _, url := range chain.RpcURL {
			if !strings.Contains(url, "INFURA_API_KEY") && chain.CheckRpcURL(url) {
				return &Chain{
					Name:    chain.Name,
					ChainID: chain.ChainID,
					RpcURL:  []string{url},
				}
			}
		}

		break
	}

	return nil
}
