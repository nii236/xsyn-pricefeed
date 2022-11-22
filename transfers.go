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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

func Scrape(client *ethclient.Client, fromBlock int64, toBlock int64, chainID int64, tokenAddr common.Address, tokenSymbol string, tokenDecimals int) error {
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
		return fmt.Errorf("contract abi: %w", err)
	}
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		return fmt.Errorf("dial eth node: %w", err)
	}
	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	total := 0
	for _, vLog := range logs {
		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			ev, err := contractAbi.Unpack("Transfer", vLog.Data)
			if err != nil {
				return fmt.Errorf("unpack log: %w", err)
			}
			from := common.HexToAddress(vLog.Topics[1].Hex())
			to := common.HexToAddress(vLog.Topics[2].Hex())
			amtBig := ev[0].(*big.Int)
			amt := decimal.NewFromBigInt(amtBig, 0)

			result := &Transfer{vLog.BlockNumber, vLog.Index, chainID, tokenAddr, tokenSymbol, tokenDecimals, vLog.TxHash, from, to, amt}
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
	return nil
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
}
