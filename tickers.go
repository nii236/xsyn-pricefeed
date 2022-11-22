package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Tickers struct {
	*EthClient
	Mainnet        *ethclient.Client
	Goerli         *ethclient.Client
	SUPSAddr       common.Address
	GoerliSUPSAddr common.Address
}

const BaseMainnetBlock = 15879854
const BaseGoerliBlock = 7859764

func (t *Tickers) CatchUp() error {
	for {
		blockHeightMainnet, err := GetInt(KeyBlockHeightMainnet, BaseMainnetBlock)
		if err != nil {
			return fmt.Errorf("get BlockHeight: %w", err)
		}
		lastBlockMainnetSUPS, err := GetInt(KeyLastBlockMainnetSups, BaseMainnetBlock)
		if err != nil {
			return fmt.Errorf("get LastBlock: %w", err)
		}
		if lastBlockMainnetSUPS == blockHeightMainnet {
			break
		}

		err = t.TickMainnetSUPS()
		if err != nil {
			return fmt.Errorf("speedup tickblock: %w", err)
		}
	}

	for {
		blockHeightGoerli, err := GetInt(KeyBlockHeightGoerli, BaseGoerliBlock)
		if err != nil {
			return fmt.Errorf("get BlockHeight: %w", err)
		}
		lastBlockGoerliSUPS, err := GetInt(KeyLastBlockGoerliSups, BaseGoerliBlock)
		if err != nil {
			return fmt.Errorf("get LastBlock: %w", err)
		}
		if lastBlockGoerliSUPS == blockHeightGoerli {
			break
		}

		err = t.TickGoerliSUPS()
		if err != nil {
			return fmt.Errorf("speedup tickblock: %w", err)
		}
	}

	return nil
}

func (t *Tickers) Start() {
	log.Info().Msg("fast forwarding...")
	err := t.CatchUp()
	if err != nil {
		log.Err(err).Msg("catching up")
	}

	err = t.TickPrice()
	if err != nil {
		log.Err(err).Msg("tick price")
		return
	}
	err = t.TickGoerliSUPS()
	if err != nil {
		log.Err(err).Msg("tick price")
		return
	}
	err = t.TickMainnetSUPS()
	if err != nil {
		log.Err(err).Msg("tick price")
		return
	}
	err = t.TickBlockHeightMainnet()
	if err != nil {
		log.Err(err).Msg("tick block height")
		return
	}

	log.Info().Msg("starting tickers")
	tickerFast := time.NewTicker(12 * time.Second)
	tickerMedium := time.NewTicker(20 * time.Second)
	tickerSlow := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-tickerMedium.C:
			err = t.TickMainnetSUPS()
			if err != nil {
				log.Err(err).Msg("tick testnet ETH")
			}
			err = t.TickGoerliSUPS()
			if err != nil {
				log.Err(err).Msg("tick goerli SUPS")
			}
		case <-tickerFast.C:
			err := t.TickBlockHeightMainnet()
			if err != nil {
				log.Err(err).Msg("tick block height mainnet")
			}
			err = t.TickBlockHeightGoerli()
			if err != nil {
				log.Err(err).Msg("tick block height goerli")
			}
		case <-tickerSlow.C:
			err := t.TickPrice()
			if err != nil {
				log.Err(err).Msg("tick price")
			}
		}
	}
}

func (t *Tickers) TickMainnetEth() error { return ErrNotImplemented }
func (t *Tickers) TickTestnetEth() error { return ErrNotImplemented }
func (t *Tickers) TickGoerliSUPS() error {
	blockHeightGoerli, err := GetInt(KeyBlockHeightGoerli, BaseGoerliBlock)
	if err != nil {
		return fmt.Errorf("start ticker: %w", err)
	}
	lastBlockSUPSGoerli, err := GetInt(KeyLastBlockGoerliSups, BaseGoerliBlock)
	if err != nil {
		return fmt.Errorf("start ticker: %w", err)
	}

	toBlock, err := t.TickSUPS(t.Goerli, lastBlockSUPSGoerli, blockHeightGoerli, 5, 18, t.GoerliSUPSAddr)
	if err != nil {
		return fmt.Errorf("tick sups Goerli: %w", err)
	}
	err = SetInt(KeyLastBlockGoerliSups, int(toBlock))
	if err != nil {
		return fmt.Errorf("set latest block sups Goerli: %w", err)
	}
	return nil
}
func (t *Tickers) TickMainnetSUPS() error {
	blockHeightMainnet, err := GetInt(KeyBlockHeightMainnet, BaseMainnetBlock)
	if err != nil {
		return fmt.Errorf("start ticker: %w", err)
	}
	lastBlockSUPSMainnet, err := GetInt(KeyLastBlockMainnetSups, BaseMainnetBlock)
	if err != nil {
		return fmt.Errorf("start ticker: %w", err)
	}

	toBlock, err := t.TickSUPS(t.Mainnet, lastBlockSUPSMainnet, blockHeightMainnet, 1, 18, t.SUPSAddr)
	if err != nil {
		return fmt.Errorf("tick sups mainnet: %w", err)
	}
	err = SetInt(KeyLastBlockMainnetSups, int(toBlock))
	if err != nil {
		return fmt.Errorf("set latest block sups mainnet: %w", err)
	}
	return nil
}

func (t *Tickers) TickSUPS(client *ethclient.Client, lastBlock int, blockHeight int, chainID int64, decimals int, tokenAddr common.Address) (int64, error) {
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

	err := ScrapeSUPS(client, fromBlock, toBlock, chainID, tokenAddr, "SUPS", decimals)
	if err != nil {
		return 0, fmt.Errorf("scrape transfers: %w", err)
	}

	return toBlock, nil
}
func (t *Tickers) TickBlockHeightGoerli() error {
	height, err := t.Goerli.BlockNumber(context.TODO())
	if err != nil {
		return fmt.Errorf("block height: %w", err)
	}
	err = SetInt(KeyBlockHeightGoerli, int(height))
	if err != nil {
		return fmt.Errorf("set block height: %w", err)
	}
	log.Info().Int("block_height_goerli", int(height)).Msg("scraping block height")
	return nil
}
func (t *Tickers) TickBlockHeightMainnet() error {
	height, err := t.Mainnet.BlockNumber(context.TODO())
	if err != nil {
		return fmt.Errorf("block height: %w", err)
	}
	err = SetInt(KeyBlockHeightMainnet, int(height))
	if err != nil {
		return fmt.Errorf("set block height: %w", err)
	}
	log.Info().Int("block_height_mainnet", int(height)).Msg("scraping block height")
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
