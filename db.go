package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/georgysavva/scany/v2/pgxscan"
)

var ErrNotImplemented = errors.New("not implemented")
var conn *pgxpool.Pool

type KVKey string

const KeyLastBlock KVKey = "last_block"
const KeyBlockHeight KVKey = "block_height"

const KeyLastBlockTestnet KVKey = "last_block_testnet"
const KeyBlockHeightTestnet KVKey = "block_height_testnet"

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

func GetInt(key string) (int, error) {
	resultStr, err := Get(key)
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

func Get(key string) (string, error) {
	q := `SELECT value FROM kv WHERE key = $1 LIMIT 1`
	var result string
	err := pgxscan.Get(context.TODO(), conn, &result, q, key)
	if err != nil {
		return "", fmt.Errorf("set block: %w", err)
	}

	return result, nil
}

func AddPrice(price *PriceResponse) error {
	q := `INSERT INTO prices (sups_price_cents, eth_price_cents, bnb_price_cents) VALUES ($1, $2, $3)`
	_, err := conn.Exec(context.TODO(), q, price.SUPSUSD, price.ETHUSD, price.BNBUSD)
	if err != nil {
		return fmt.Errorf("add price: %w", err)
	}
	return nil
}
func AddTransfer(transfer *Transfer) error {
	q := `INSERT INTO transfers	(block, log_index, chain_id, contract, symbol, decimals, tx_id, from_address, to_address, amount) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

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
	)
	if err != nil {
		return fmt.Errorf("insert transfer: %w", err)
	}
	return nil

}
func Transfers(fromBlock int, toAddr common.Address) ([]*TransferAPIResponse, error) {
	q := `SELECT * FROM transfers WHERE block > $1 AND to_address = $2`
	resultDB := []*TransferRecord{}
	err := pgxscan.Select(context.TODO(), conn, &resultDB, q, fromBlock, toAddr.Hex())
	if err != nil {
		return nil, fmt.Errorf("set block: %w", err)
	}

	result := []*TransferAPIResponse{}
	blockHeightStr, err := Get(string(KeyBlockHeight))
	if err != nil {
		return nil, fmt.Errorf("get block height: %w", err)
	}
	blockHeight, err := strconv.Atoi(blockHeightStr)
	if err != nil {
		return nil, fmt.Errorf("convert block height: %w", err)
	}

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
			ValueDecimals:   record.Decimals,
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
	ValueDecimals   int    `json:"value_decimals"`
}
