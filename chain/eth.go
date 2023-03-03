package chain

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"math/big"
	"time"
)

type EthClient struct {
	Client  *ethclient.Client
	timeout time.Duration
}

func NewClient(url string) *EthClient {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client, err := ethclient.DialContext(ctx, url)
	if err != nil {
		log.Fatal(err)
	}

	return &EthClient{
		Client:  client,
		timeout: time.Minute * 10,
	}
}

func (c *EthClient) GetBlockNumber() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	blockNumber, err := c.Client.BlockNumber(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get block number")
	}
	return blockNumber, nil
}

// GetNonce get current nonce
func (c *EthClient) GetNonce(address common.Address) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	nonce, err := c.Client.NonceAt(ctx, address, nil)
	if err != nil {
		return 0, errors.Wrap(err, "failed to retrieve account nonce")
	}

	return nonce, nil
}

func (c *EthClient) GetGasPrice() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	gasPrice, err := c.Client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to suggest gas price")
	}

	return gasPrice, err
}

func (c *EthClient) GetGasLimit(msg ethereum.CallMsg) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	tempGasLimit, err := c.Client.EstimateGas(ctx, msg)
	if err != nil {
		return 0, errors.Wrap(err, "failed to estimate gas limit")
	}

	return uint64(float64(tempGasLimit) * 1.3), nil
}

func (c *EthClient) GetChainId() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	chainID, err := c.Client.NetworkID(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get chain id")
	}

	return chainID, nil
}

func (c *EthClient) GetTxByHash(hash common.Hash) (*types.Transaction, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	tx, isPending, err := c.Client.TransactionByHash(ctx, hash)
	if err != nil {
		if err.Error() == "not found" {
			return nil, false, err
		}
		return nil, false, errors.Wrapf(err, "failed to get tx by hash:%s", hash.String())
	}

	return tx, isPending, err
}

func (c *EthClient) GetTxReceipt(hash *common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	receipt, err := c.Client.TransactionReceipt(ctx, *hash)
	if err != nil {
		return nil, err
	}
	return receipt, err
}

func (c *EthClient) SendTx(tx *types.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	err := c.Client.SendTransaction(ctx, tx)
	if err != nil {
		err = errors.Wrap(err, "fail send transaction")
	}
	return err
}

func (c *EthClient) EthCall(payload []byte, ctrAdr common.Address, blockNumber *big.Int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	msg := ethereum.CallMsg{
		To:   &ctrAdr,
		Data: payload,
	}
	output, err := c.Client.CallContract(ctx, msg, blockNumber)
	return output, err
}

type CallMethodOpts struct {
	Nonce    uint64
	Value    *big.Int
	GasPrice *big.Int
	GasLimit uint64
	ChainId  *big.Int
}

func privateToAddress(priKey string) (*common.Address, *ecdsa.PrivateKey, error) {
	privateKey, err := crypto.HexToECDSA(priKey)
	if err != nil {
		return nil, nil, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &address, privateKey, err
}

func (c *EthClient) BuildCallMethodTxWithPayload(payload []byte, contractAddress string, priKey string, opts *CallMethodOpts) (*types.Transaction, error) {
	address, privateKey, err := privateToAddress(priKey)
	if err != nil {
		return nil, err
	}

	var (
		value    = big.NewInt(0)
		gasPrice *big.Int
		gasLimit uint64
		nonce    uint64
		chainId  *big.Int
	)

	if opts != nil {
		if opts.Value != nil {
			value = opts.Value
		}
		if opts.GasPrice != nil {
			gasPrice = opts.GasPrice
		}
		if opts.GasLimit != 0 {
			gasLimit = opts.GasLimit
		}
		if opts.Nonce != 0 {
			nonce = opts.Nonce
		}
		if opts.ChainId != nil {
			chainId = opts.ChainId
		}
	}

	if nonce == 0 {
		nonce, err = c.GetNonce(*address)
		if err != nil {
			return nil, err
		}
	}

	if gasPrice == nil {
		gasPrice, err = c.GetGasPrice()
		if err != nil {
			return nil, err
		}
	}

	ctrAdr := common.HexToAddress(contractAddress)
	if gasLimit == 0 {
		msg := ethereum.CallMsg{
			From:     *address,
			To:       &ctrAdr,
			GasPrice: gasPrice,
			Value:    value,
			Data:     payload,
		}
		gasLimit, err = c.GetGasLimit(msg)
		if err != nil {
			return nil, err
		}
	}

	unSignedTx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &ctrAdr,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     payload,
	})
	signedTx, err := types.SignTx(unSignedTx, types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "fail to sign tx")
	}

	return signedTx, nil
}

func (c *EthClient) Sender(chainId *big.Int, tx *types.Transaction) (common.Address, error) {
	signer := types.NewEIP155Signer(chainId)
	sender, err := signer.Sender(tx)
	return sender, err
}

func (c *EthClient) CheckTxPacked(txHash common.Hash) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		receipt, err := c.Client.TransactionReceipt(ctx, txHash)
		cancel()
		if err != nil {
			if err.Error() == "not found" {
				time.Sleep(time.Second * 1)
				continue
			}
			logx.Must(err)
		}
		if receipt.BlockNumber != nil {
			if receipt.Status == 0 {
				err = errors.Errorf("transaction: %s execute failed", txHash.String())
				logx.Must(err)
			}
			break
		}
		time.Sleep(time.Second * 1)
	}
}
