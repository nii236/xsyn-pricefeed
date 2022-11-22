package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"xsyn-pricefeed/erc20"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

func ContainsAddress(wantAddr common.Address, whitelistedAddr []common.Address) bool {
	for _, addr := range whitelistedAddr {
		if addr == wantAddr {
			return true
		}
		continue
	}
	return false
}

func ScrapeETH(client *ethclient.Client, fromBlock int64, toBlock int64, whitelistedAddr []common.Address, chainID int64) (int, error) {
	total := 0
	for blockNumber := fromBlock; blockNumber < toBlock; blockNumber++ {
		block, err := client.BlockByNumber(context.TODO(), big.NewInt(blockNumber))
		if err != nil {
			return 0, fmt.Errorf("scrape eth get block: %w", err)
		}
		timestamp := block.Time()
		txes := block.Transactions()
		for i, tx := range txes {
			if tx.To() != nil && ContainsAddress(*tx.To(), whitelistedAddr) {
				msg, err := tx.AsMessage(types.LatestSignerForChainID(big.NewInt(chainID)), nil)
				if err != nil {
					log.Err(err).
						Uint64("block", block.NumberU64()).
						Int64("chain_id", chainID).
						Str("symbol", "ETH").
						Int("decimals", 18).
						Str("tx_id", tx.Hash().Hex()).
						Msg("eth native transfer")
					continue
				}
				to := common.HexToAddress("0x0")
				if msg.To() != nil {
					to = *msg.To()
				}

				result := &Transfer{
					Block:       block.Number().Uint64(),
					LogIndex:    uint(i),
					Symbol:      "ETH",
					Decimals:    18,
					ChainID:     chainID,
					TxID:        tx.Hash(),
					FromAddress: msg.From(),
					ToAddress:   to,
					Amount:      decimal.NewFromBigInt(msg.Value(), 0),
					CreatedAt:   timestamp,
				}

				err = AddTransfer(result)
				if err != nil {
					log.Err(err).
						Uint64("block", result.Block).
						Int64("chain_id", result.ChainID).
						Str("contract", result.Contract.Hex()).
						Str("symbol", result.Symbol).
						Int("decimals", result.Decimals).
						Str("tx_id", result.TxID.Hex()).
						Str("from_address", result.FromAddress.Hex()).
						Str("to_address", result.ToAddress.Hex()).
						Str("amount", result.Amount.String()).
						Msg("insert transfer")
					continue
				}
				total++
			}
		}
	}
	return total, nil
}
func ScrapeSUPS(client *ethclient.Client, fromBlock int64, toBlock int64, chainID int64, tokenAddr common.Address, tokenSymbol string, tokenDecimals int) (int, error) {
	total := 0
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
		Addresses: []common.Address{tokenAddr},
		Topics: [][]common.Hash{
			{common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")}, // Transfer(address,address,uint256)
		},
	}
	contractAbi, err := abi.JSON(strings.NewReader(string(erc20.Erc20ABI)))
	if err != nil {
		return 0, fmt.Errorf("contract abi: %w", err)
	}
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		return 0, fmt.Errorf("dial eth node: %w", err)
	}
	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	for _, vLog := range logs {
		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			ev, err := contractAbi.Unpack("Transfer", vLog.Data)
			if err != nil {
				return 0, fmt.Errorf("unpack log: %w", err)
			}
			from := common.HexToAddress(vLog.Topics[1].Hex())
			to := common.HexToAddress(vLog.Topics[2].Hex())
			amtBig := ev[0].(*big.Int)
			amt := decimal.NewFromBigInt(amtBig, 0)
			block, err := client.BlockByHash(context.TODO(), vLog.BlockHash)
			if err != nil {
				log.Err(err).Msg("get block")
				continue
			}
			result := &Transfer{vLog.BlockNumber, vLog.Index, chainID, tokenAddr, tokenSymbol, tokenDecimals, vLog.TxHash, from, to, amt, block.Time()}
			err = AddTransfer(result)
			if err != nil {
				log.Err(err).
					Uint64("block", result.Block).
					Int64("chain_id", result.ChainID).
					Str("contract", result.Contract.Hex()).
					Str("symbol", result.Symbol).
					Int("decimals", result.Decimals).
					Str("tx_id", result.TxID.Hex()).
					Str("from_address", result.FromAddress.Hex()).
					Str("to_address", result.ToAddress.Hex()).
					Str("amount", result.Amount.String()).
					Msg("insert transfer")
				continue
			}
			total++
		}
	}
	return total, nil
}

type Transfer struct {
	Block       uint64
	LogIndex    uint
	ChainID     int64
	Contract    common.Address
	Symbol      string
	Decimals    int
	TxID        common.Hash
	FromAddress common.Address
	ToAddress   common.Address
	Amount      decimal.Decimal
	CreatedAt   uint64
}
