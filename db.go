package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/georgysavva/scany/v2/pgxscan"
)

var ErrNotImplemented = errors.New("not implemented")
var conn *pgxpool.Pool

type KVKey string

const KeyScrapeRangeLookbackEth KVKey = "scrape_range_lookback_eth"
const KeyScrapeRangeLookbackSups KVKey = "scrape_range_lookback_sups"

const KeyScrapeRangeEth KVKey = "scrape_range_eth"
const KeyScrapeRangeSups KVKey = "scrape_range_sups"
const KeyBlockHeightGoerli KVKey = "block_height_goerli"
const KeyLastBlockGoerliSups KVKey = "last_block_goerli_sups"
const KeyLastBlockGoerliEth KVKey = "last_block_goerli_eth"
const KeyBlockHeightMainnet KVKey = "block_height_mainnet"
const KeyLastBlockMainnetSups KVKey = "last_block_mainnet_sups"
const KeyLastBlockMainnetEth KVKey = "last_block_mainnet_eth"

func Connect(connString string) error {
	var err error
	conn, err = pgxpool.New(context.TODO(), connString)
	if err != nil {
		return fmt.Errorf("conn: %w", err)
	}
	return nil

}

func Exists(key KVKey) (bool, error) {
	q := `SELECT COUNT id FROM kv WHERE key = $1`
	var result int
	err := pgxscan.Get(context.TODO(), conn, &result, q, key)
	if err != nil {
		return false, fmt.Errorf("set block: %w", err)
	}
	return true, nil
}

func Set(key KVKey, value string) error {
	q := `INSERT INTO kv (key, value) VALUES ($1, $2) ON CONFLICT (key) DO UPDATE SET value=EXCLUDED.value;`

	_, err := conn.Exec(context.TODO(), q, key, value)
	if err != nil {
		return fmt.Errorf("set block: %w", err)
	}

	return nil
}

func GetInt(key KVKey, defaultValue int) (int, error) {
	resultStr, err := Get(key, defaultValue)
	if err != nil {
		return 0, fmt.Errorf("get int: %w", err)
	}
	result, err := strconv.Atoi(resultStr)
	if err != nil {
		return 0, fmt.Errorf("convert int: %w", err)
	}
	return result, nil
}
func SetInt(key KVKey, value int) error {
	return Set(key, strconv.Itoa(value))
}

func WhitelistedAddresses(chainID int) ([]common.Address, error) {
	q := `SELECT address FROM whitelisted_addresses WHERE chain_id = $1`
	resultStr := []string{}
	err := pgxscan.Select(context.TODO(), conn, &resultStr, q, chainID)
	if err != nil {
		return nil, fmt.Errorf("set block: %w", err)
	}
	result := []common.Address{}
	for _, addrStr := range resultStr {
		result = append(result, common.HexToAddress(addrStr))
	}

	return result, nil
}

func Get(key KVKey, defaultValue ...int) (string, error) {
	q := `SELECT value FROM kv WHERE key = $1 LIMIT 1`
	var result string
	err := pgxscan.Get(context.TODO(), conn, &result, q, key)
	if err != nil && errors.Is(err, pgx.ErrNoRows) && len(defaultValue) > 0 {
		return strconv.Itoa(defaultValue[0]), SetInt(key, defaultValue[0])
	}
	if err != nil {
		return "", fmt.Errorf("set block: %w", err)
	}

	return result, nil
}

func AddPrice(price *PriceResponse) error {
	if price.SUPSUSD == decimal.Zero {
		price.SUPSUSD = decimal.NewFromFloat(0.8)
	}
	q := `INSERT INTO prices (sups_price_cents, eth_price_cents, bnb_price_cents) VALUES ($1, $2, $3)`
	_, err := conn.Exec(context.TODO(), q, price.SUPSUSD, price.ETHUSD, price.BNBUSD)
	if err != nil {
		return fmt.Errorf("add price: %w", err)
	}
	return nil
}
func AddTransfer(transfer *Transfer) error {
	q := `INSERT INTO transfers	(block, log_index, chain_id, contract, symbol, decimals, tx_id, from_address, to_address, amount, timestamp) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := conn.Exec(context.TODO(), q,
		transfer.Block,
		transfer.LogIndex,
		transfer.ChainID,
		transfer.Contract.Hex(),
		transfer.Symbol,
		transfer.Decimals,
		transfer.TxID.Hex(),
		transfer.FromAddress.Hex(),
		transfer.ToAddress.Hex(),
		transfer.Amount.String(),
		transfer.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert transfer: %w", err)
	}
	return nil

}
func Transfers(symbol string, blockHeight int, sinceBlock int, chainID int) ([]*TransferAPIResponse, error) {
	q := `SELECT * FROM transfers WHERE block > $1 AND chain_id = $2 AND symbol = $3 ORDER BY block DESC`
	resultDB := []*TransferRecord{}
	err := pgxscan.Select(context.TODO(), conn, &resultDB, q, sinceBlock, chainID, strings.ToUpper(symbol))
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("set block: %w", err)
	}

	result := []*TransferAPIResponse{}

	for _, record := range resultDB {
		confirmations := blockHeight - int(record.Block)
		result = append(result, &TransferAPIResponse{
			TxHash:          record.TxID,
			LogIndex:        record.LogIndex,
			Time:            record.CreatedAt.Unix(),
			Chain:           record.ChainID,
			BlockNumber:     record.Block,
			Confirmations:   confirmations,
			FromAddress:     record.FromAddress,
			ToAddress:       record.ToAddress,
			ContractAddress: record.Contract,
			Value:           record.Amount.Shift(-18).String(),
			ValueInt:        record.Amount.String(),
			Timestamp:       record.Timestamp,
			ValueDecimals:   record.Decimals,
			Symbol:          record.Symbol,
		})

	}
	return result, nil
}

type TransferRecord struct {
	ID          uuid.UUID
	Block       uint64
	LogIndex    uint
	ChainID     int64
	Contract    string
	Symbol      string
	Decimals    int
	TxID        string
	FromAddress string
	ToAddress   string
	Amount      decimal.Decimal
	Timestamp   int64
	CreatedAt   time.Time
}
type TransferAPIResponse struct {
	TxHash          string `json:"tx_hash"`
	LogIndex        uint   `json:"log_index"`
	Time            int64  `json:"time"`
	Chain           int64  `json:"chain"`
	BlockNumber     uint64 `json:"block_number"`
	Confirmations   int    `json:"confirmations"`
	FromAddress     string `json:"from_address"`
	ToAddress       string `json:"to_address"`
	ContractAddress string `json:"contract_address"`
	Value           string `json:"value"`
	ValueInt        string `json:"value_int"`
	Timestamp       int64  `json:"timestamp"`
	ValueDecimals   int    `json:"value_decimals"`
	Symbol          string `json:"symbol"`
}
