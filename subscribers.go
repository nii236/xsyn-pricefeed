package main

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Subscriber struct {
	Mainnet *ethclient.Client
	Goerli  *ethclient.Client
}

func (s *Subscriber) Start() {
	log.Info().Msg("start subscribers")
	go func() {
		ctx := context.Background()
		log.Info().Str("chain", "mainnet").Msg("start header listener")
		mainnetHead := make(chan *types.Header)
		headSub, err := s.Mainnet.SubscribeNewHead(ctx, mainnetHead)
		if err != nil {
			log.Err(err).Msg("subscribe head")
			return
		}
		defer headSub.Unsubscribe()
		for {
			select {
			case err := <-headSub.Err():
				log.Err(err).Msg("receive head")
				continue
			case head := <-mainnetHead:
				err = SetInt(KeyBlockHeightMainnet, int(head.Number.Int64()))
				if err != nil {
					log.Err(err).Msg("set head")
					continue
				}
			}
		}
	}()
	go func() {
		ctx := context.Background()
		log.Info().Str("chain", "goerli").Msg("start header listener")
		testnetHead := make(chan *types.Header)
		headSub, err := s.Goerli.SubscribeNewHead(ctx, testnetHead)
		if err != nil {
			log.Err(err).Msg("subscribe head")
			return
		}
		defer headSub.Unsubscribe()
		for {
			select {
			case err := <-headSub.Err():
				log.Err(err).Msg("receive head")
				continue
			case head := <-testnetHead:
				err = SetInt(KeyBlockHeightGoerli, int(head.Number.Int64()))
				if err != nil {
					log.Err(err).Msg("set head")
					continue
				}
			}
		}
	}()
}
