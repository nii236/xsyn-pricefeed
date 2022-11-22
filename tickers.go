package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Tickers struct {
	*EthClient
	SUPSAddr common.Address
}

func (t *Tickers) CatchUp() error {
	for {
		blockHeight, err := GetInt(string(KeyBlockHeight))
		if err != nil {
			return fmt.Errorf("get BlockHeight: %w", err)
		}
		lastBlock, err := GetInt(string(KeyLastBlock))
		if err != nil {
			return fmt.Errorf("get LastBlock: %w", err)
		}
		if lastBlock == blockHeight {
			break
		}
		err = t.TickBlock()
		if err != nil {
			return fmt.Errorf("speedup tickblock: %w", err)
		}
	}
	return nil
}

func (t *Tickers) Start() {
	blockHeight, err := GetInt(string(KeyBlockHeight))
	if err != nil {
		log.Fatal().Err(err).Msg("start ticker")
	}
	lastBlock, err := GetInt(string(KeyLastBlock))
	if err != nil {
		log.Fatal().Err(err).Msg("start ticker")
	}

	if lastBlock < blockHeight {
		log.Info().Msg("fast forwarding...")
		err = t.CatchUp()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to catchup")
		}
	}

	log.Info().Msg("starting tickers")
	err = t.TickPrice()
	if err != nil {
		log.Err(err).Msg("tick price")
	}
	err = t.TickBlock()
	if err != nil {
		log.Err(err).Msg("tick block")
	}
	err = t.TickBlockHeight()
	if err != nil {
		log.Err(err).Msg("tick block height")
	}

	tb := time.NewTicker(20 * time.Second)
	tp := time.NewTicker(60 * time.Second)
	tbh := time.NewTicker(12 * time.Second)

	for {
		select {
		case <-tbh.C:
			err := t.TickBlockHeight()
			if err != nil {
				log.Err(err).Msg("tick block")
			}
		case <-tb.C:
			err := t.TickBlock()
			if err != nil {
				log.Err(err).Msg("tick block")
			}
		case <-tp.C:
			err := t.TickPrice()
			if err != nil {
				log.Err(err).Msg("tick price")
			}
		}
	}
}

func (t *Tickers) TickBlock() error {
	lastBlock, err := GetInt(string(KeyLastBlock))
	if err != nil {
		return fmt.Errorf("get last block: %w", err)
	}

	blockHeight, err := GetInt(string(KeyBlockHeight))
	if err != nil {
		return fmt.Errorf("get block height: %w", err)
	}
	fromBlock := int64(lastBlock - 50)
	toBlock := int64(lastBlock + 9000)
	if toBlock > int64(blockHeight) {
		toBlock = int64(blockHeight)
	}
	if fromBlock < 0 {
		fromBlock = 0
	}

	log.Info().
		Int64("last_block", int64(lastBlock)).
		Int64("from_block", fromBlock).
		Int64("to_block", toBlock).
		Int64("block_height", int64(blockHeight)).
		Msg("scraping transfers")

	err = Scrape(t.EthClient.Client, fromBlock, toBlock, 1, t.SUPSAddr, "SUPS", 18)
	if err != nil {
		return fmt.Errorf("scrape transfers: %w", err)
	}

	err = SetInt(KeyLastBlock, int(toBlock))
	if err != nil {
		return fmt.Errorf("set last block: %w", err)
	}

	return nil
}

func (t *Tickers) TickBlockHeight() error {
	height, err := t.Client.BlockNumber(context.TODO())
	if err != nil {
		return fmt.Errorf("block height: %w", err)
	}
	err = SetInt(KeyBlockHeight, int(height))
	if err != nil {
		return fmt.Errorf("set block height: %w", err)
	}
	log.Info().Int("block_height", int(height)).Msg("scraping block height")
	return nil
}

func (t *Tickers) TickPrice() error {
	log.Info().Msg("scraping prices")
	supsusd, err := t.SUPSUSD()
	if err != nil {
		return fmt.Errorf("get supsusd price: %w", err)
	}
	ethusd, err := t.ETHUSD()
	if err != nil {
		return fmt.Errorf("get ethusd price: %w", err)
	}
	bnbusd, err := t.BNBUSD()
	if err != nil {
		return fmt.Errorf("get ethusd price: %w", err)
	}
	result := &PriceResponse{time.Now().Unix(), supsusd, ethusd, bnbusd}
	return AddPrice(result)
}
