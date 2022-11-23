package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Tickers struct {
	ScrapeMainnetETH  bool
	ScrapeMainnetSUPS bool
	ScrapeGoerliETH   bool
	ScrapeGoerliSUPS  bool

	*EthClient
	Mainnet        *ethclient.Client
	Goerli         *ethclient.Client
	SUPSAddr       common.Address
	GoerliSUPSAddr common.Address
}

const BaseMainnetBlock = 15879854
const BaseGoerliBlock = 7859764

func (t *Tickers) CatchUpMainnetETH() error {
	iter := 0
	for {
		log.Info().Int("tick", iter).Str("chain", "mainnet").Str("symbol", "eth").Msg("fast forwarding...")
		if !t.ScrapeMainnetETH {
			break
		}
		blockHeightMainnet, err := GetInt(KeyBlockHeightMainnet, BaseMainnetBlock)
		if err != nil {
			return fmt.Errorf("get BlockHeight: %w", err)
		}
		lastBlockMainnetETH, err := GetInt(KeyLastBlockMainnetEth, BaseMainnetBlock)
		if err != nil {
			return fmt.Errorf("get LastBlock: %w", err)
		}
		if lastBlockMainnetETH == blockHeightMainnet {
			break
		}

		err = t.TickMainnetEth()
		if err != nil {
			return fmt.Errorf("speedup tickblock: %w", err)
		}
		iter++
	}
	return nil
}
func (t *Tickers) CatchUpGoerliETH() error {
	iter := 0
	for {
		if !t.ScrapeGoerliETH {
			break
		}
		log.Info().Int("tick", iter).Str("chain", "goerli").Str("symbol", "eth").Msg("fast forwarding...")
		blockHeightGoerli, err := GetInt(KeyBlockHeightGoerli, BaseGoerliBlock)
		if err != nil {
			return fmt.Errorf("get BlockHeight: %w", err)
		}
		lastBlockGoerliETH, err := GetInt(KeyLastBlockGoerliEth, BaseGoerliBlock)
		if err != nil {
			return fmt.Errorf("get LastBlock: %w", err)
		}
		if lastBlockGoerliETH == blockHeightGoerli {
			break
		}

		err = t.TickGoerliETH()
		if err != nil {
			return fmt.Errorf("speedup tickblock: %w", err)
		}
		iter++
	}
	return nil
}
func (t *Tickers) CatchUpGoerliSUPS() error {
	iter := 0
	for {
		if !t.ScrapeGoerliSUPS {
			break
		}
		log.Info().Int("tick", iter).Str("chain", "goerli").Str("symbol", "sups").Msg("fast forwarding...")
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
		iter++
	}
	return nil
}
func (t *Tickers) CatchUpMainnetSUPS() error {
	iter := 0
	for {
		if !t.ScrapeMainnetSUPS {
			break
		}
		log.Info().Int("tick", iter).Str("chain", "mainnet").Str("symbol", "sups").Msg("fast forwarding...")
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
		iter++
	}
	return nil
}
func (t *Tickers) CatchUp() error {
	messages := make(chan string)
	wg := &sync.WaitGroup{}

	wgResult := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		worker := "mainnet_eth"
		log.Info().Str("worker", worker).Msg("start catchup worker")
		defer wg.Done()
		err := t.CatchUpMainnetETH()
		if err != nil {
			log.Err(err).Msg("catch up mainnet eth")
		}
		messages <- worker
	}(wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		worker := "goerli_eth"
		log.Info().Str("worker", worker).Msg("start catchup worker")
		defer wg.Done()
		err := t.CatchUpGoerliETH()
		if err != nil {
			log.Err(err).Msg("catch up goerli eth")
		}
		messages <- worker
	}(wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		worker := "goerli_sups"
		log.Info().Str("worker", worker).Msg("start catchup worker")
		defer wg.Done()
		err := t.CatchUpGoerliSUPS()
		if err != nil {
			log.Err(err).Msg("catch up goerli sups")
		}
		messages <- worker
	}(wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		worker := "mainnet_sups"
		log.Info().Str("worker", worker).Msg("start catchup worker")
		defer wg.Done()
		err := t.CatchUpMainnetSUPS()
		if err != nil {
			log.Err(err).Msg("catch up mainnet sups")
		}
		messages <- worker
	}(wg)

	wgResult.Add(1)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for msg := range messages {
			log.Info().Str("worker", msg).Msg("catchup worker finished")
		}
	}()
	wg.Wait()
	close(messages)
	wgResult.Wait()
	return nil
}

func (t *Tickers) Start() {
	err := t.TickBlockHeightMainnet()
	if err != nil {
		log.Err(err).Msg("tick block height mainnet")
		return
	}
	err = t.TickBlockHeightGoerli()
	if err != nil {
		log.Err(err).Msg("tick block height goerli")
	}
	err = t.TickPrice()
	if err != nil {
		log.Err(err).Msg("tick price")
		return
	}
	err = t.CatchUp()
	if err != nil {
		log.Err(err).Msg("catching up")
	}

	log.Info().Msg("starting tickers")
	tickerFast := time.NewTicker(12 * time.Second)
	tickerMedium := time.NewTicker(20 * time.Second)
	tickerSlow := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-tickerMedium.C:
			log.Info().Str("type", "medium").Msg("running ticker")
			if t.ScrapeMainnetSUPS {
				err = t.TickMainnetSUPS()
				if err != nil {
					log.Err(err).Msg("tick mainnet ETH")
				}
			}

			if t.ScrapeGoerliSUPS {
				err = t.TickGoerliSUPS()
				if err != nil {
					log.Err(err).Msg("tick goerli SUPS")
				}
			}

			if t.ScrapeMainnetETH {
				err = t.TickMainnetEth()
				if err != nil {
					log.Err(err).Msg("tick mainnet ETH")
				}
			}
			if t.ScrapeGoerliETH {
				err = t.TickGoerliETH()
				if err != nil {
					log.Err(err).Msg("tick goerli ETH")
				}
			}
		case <-tickerFast.C:
			log.Info().Str("type", "fast").Msg("running ticker")
			err := t.TickBlockHeightMainnet()
			if err != nil {
				log.Err(err).Msg("tick block height mainnet")
			}
			err = t.TickBlockHeightGoerli()
			if err != nil {
				log.Err(err).Msg("tick block height goerli")
			}
		case <-tickerSlow.C:
			log.Info().Str("type", "slow").Msg("running ticker")
			err := t.TickPrice()
			if err != nil {
				log.Err(err).Msg("tick price")
			}
		}
	}
}

func (t *Tickers) TickMainnetEth() error {
	blockHeightMainnet, err := GetInt(KeyBlockHeightMainnet, BaseMainnetBlock)
	if err != nil {
		return fmt.Errorf("start ticker: %w", err)
	}
	lastBlockETHMainnet, err := GetInt(KeyLastBlockMainnetEth, BaseMainnetBlock)
	if err != nil {
		return fmt.Errorf("start ticker: %w", err)
	}
	toBlock, err := t.TickEth(t.Mainnet, lastBlockETHMainnet, blockHeightMainnet, 1)
	if err != nil {
		return fmt.Errorf("tick eth Mainnet: %w", err)
	}
	err = SetInt(KeyLastBlockMainnetEth, int(toBlock))
	if err != nil {
		return fmt.Errorf("set latest block sups Mainnet: %w", err)
	}
	return nil
}
func (t *Tickers) TickGoerliETH() error {
	blockHeightGoerli, err := GetInt(KeyBlockHeightGoerli, BaseGoerliBlock)
	if err != nil {
		return fmt.Errorf("start ticker: %w", err)
	}
	lastBlockETHGoerli, err := GetInt(KeyLastBlockGoerliEth, BaseGoerliBlock)
	if err != nil {
		return fmt.Errorf("start ticker: %w", err)
	}
	toBlock, err := t.TickEth(t.Goerli, lastBlockETHGoerli, blockHeightGoerli, 5)
	if err != nil {
		return fmt.Errorf("tick eth goerli: %w", err)
	}
	err = SetInt(KeyLastBlockGoerliEth, int(toBlock))
	if err != nil {
		return fmt.Errorf("set latest block sups Goerli: %w", err)
	}
	return nil
}

func (t *Tickers) TickEth(client *ethclient.Client, lastBlock int, blockHeight int, chainID int64) (int64, error) {
	scrapeRange, err := GetInt(KeyScrapeRangeEth, 500)
	if err != nil {
		return 0, fmt.Errorf("get scrape range: %w", err)
	}

	fromBlock := int64(lastBlock - 50)
	toBlock := int64(lastBlock + scrapeRange)
	if toBlock > int64(blockHeight) {
		toBlock = int64(blockHeight)
	}
	if fromBlock < 0 {
		fromBlock = 0
	}

	log.Info().Int64("from_block", fromBlock).Int64("chain_id", chainID).Str("symbol", "eth").Msg("scraping transfers")

	whitelisted, err := WhitelistedAddresses(int(chainID))
	if err != nil {
		return 0, fmt.Errorf("scrape transfers: %w", err)
	}

	total, err := ScrapeETH(client, fromBlock, toBlock, whitelisted, chainID)
	if err != nil {
		return 0, fmt.Errorf("scrape transfers: %w", err)
	}
	log.Info().
		Int64("last_block", int64(lastBlock)).
		Int64("from_block", fromBlock).
		Int64("to_block", toBlock).
		Int("block_height", blockHeight).
		Int64("chain_id", chainID).
		Int("total", total).
		Str("symbol", "eth").
		Msg("scraped transfers")
	return toBlock, nil
}

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
	scrapeRange, err := GetInt(KeyScrapeRangeEth, 5000)
	if err != nil {
		return 0, fmt.Errorf("get scrape range: %w", err)
	}
	fromBlock := int64(lastBlock - 50)
	toBlock := int64(lastBlock + scrapeRange)
	if toBlock > int64(blockHeight) {
		toBlock = int64(blockHeight)
	}
	if fromBlock < 0 {
		fromBlock = 0
	}

	log.Info().Int64("from_block", fromBlock).Int64("chain_id", chainID).Str("symbol", "sups").Msg("scraping transfers")

	total, err := ScrapeSUPS(client, fromBlock, toBlock, chainID, tokenAddr, "SUPS", decimals)
	if err != nil {
		return 0, fmt.Errorf("scrape transfers: %w", err)
	}
	log.Info().
		Int64("last_block", int64(lastBlock)).
		Int64("from_block", fromBlock).
		Int64("to_block", toBlock).
		Int("block_height", blockHeight).
		Int64("chain_id", chainID).
		Int("total", total).
		Str("symbol", "sups").
		Msg("scraped transfers")
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
