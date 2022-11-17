package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"xsyn-pricefeed/ethusd"
	"xsyn-pricefeed/supseth"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shopspring/decimal"
)

var log zerolog.Logger

func main() {
	app := &cli.App{
		Name:  "serve",
		Usage: "serve price feed API",
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "port", Value: 8080, Usage: "Server port to host on"},
			&cli.StringFlag{Name: "rpc_url", Required: true, Usage: "ETH node RPC URL"},
		},
		Action: func(c *cli.Context) error {
			rpcURL := c.String("rpc_url")
			port := c.Int("port")
			return Serve(rpcURL, port)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("run cli")
	}

}

func Serve(rpcURL string, port int) error {
	log = zerolog.New(os.Stdout)

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("dial eth node: %w", err)
	}
	ethusdAddr := common.HexToAddress("0x5f4eC3Df9cbd43714FE2740f5E3616155c5b8419")
	supethAddr := common.HexToAddress("0xa1e5dc01359c2920c096f0091fc7f0bf69812ca7")

	ethusdContract, err := ethusd.NewEthusd(ethusdAddr, client)
	if err != nil {
		return fmt.Errorf("create ethusd contract: %w", err)
	}
	supsethContract, err := supseth.NewSupseth(supethAddr, client)
	if err != nil {
		return fmt.Errorf("create supseth contract: %w", err)
	}

	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(10000000),
	)
	if err != nil {
		return fmt.Errorf("memcached middleware: %w", err)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(1*time.Minute),
		cache.ClientWithRefreshKey("opn"),
	)
	if err != nil {
		return fmt.Errorf("memcached client: %w", err)
	}

	ethC := &EthClient{ethusdContract, supsethContract}
	c := &Controller{ethC}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/api/prices", cacheClient.Middleware(http.HandlerFunc(c.PricesHandler)).ServeHTTP)
	log.Info().Int("port", port).Msg("Running server")
	return http.ListenAndServe(":"+fmt.Sprintf("%d", port), r)
}

type Controller struct {
	*EthClient
}
type PriceResponse struct {
	SUPSUSD decimal.Decimal `json:"sups_usd_cents"`
	ETHUSD  decimal.Decimal `json:"eth_usd_cents"`
}

func (c *Controller) PricesHandler(w http.ResponseWriter, r *http.Request) {
	supsusd, err := c.SUPSUSD()
	if err != nil {
		log.Err(err).Msg("get supsusd price")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ethusd, err := c.ETHUSD()
	if err != nil {
		log.Err(err).Msg("get ethusd price")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := &PriceResponse{supsusd, ethusd}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Err(err).Msg("marshal json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type EthClient struct {
	*ethusd.Ethusd
	*supseth.Supseth
}

func (c *EthClient) SUPSUSD() (decimal.Decimal, error) {
	ethusdPrice, err := c.ETHUSD()
	if err != nil {
		return decimal.Zero, fmt.Errorf("query supsusd: %w", err)
	}

	result, err := c.Supseth.Slot0(&bind.CallOpts{})
	if err != nil {
		return decimal.Zero, fmt.Errorf("query slot0: %w", err)
	}

	sqrtprice := decimal.NewFromBigInt(result.SqrtPriceX96, 0)

	supsEthPrice := sqrtprice.Pow(decimal.NewFromInt(2)).
		Div(decimal.NewFromInt(2).Pow(decimal.NewFromInt(192)))

	supsUsdPrice := ethusdPrice.Div(supsEthPrice)

	return supsUsdPrice, nil
}
func (c *EthClient) ETHUSD() (decimal.Decimal, error) {
	result, err := c.Ethusd.LatestRoundData(&bind.CallOpts{})
	if err != nil {
		return decimal.Zero, fmt.Errorf("query ethusd: %w", err)
	}
	return decimal.NewFromBigInt(result.Answer, -6), nil
}
