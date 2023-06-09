package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"xsyn-pricefeed/ethusd"
	"xsyn-pricefeed/supseth"

	"github.com/gomarkdown/markdown"

	chiprometheus "xsyn-pricefeed/middleware"

	_ "github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/docgen"
	"github.com/shopspring/decimal"
)

var log zerolog.Logger

func main() {

	app := &cli.App{

		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "serve price feed API",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "log_format", Value: "console", Usage: "log formatting (json or console)", EnvVars: []string{"LOG_FORMAT"}},
					&cli.IntFlag{Name: "ttl_seconds", Value: 300, Usage: "seconds to cache the responses", EnvVars: []string{"TTL_SECONDS"}},
					&cli.IntFlag{Name: "port", Value: 8080, Usage: "Server port to host on", EnvVars: []string{"PORT"}},
					&cli.StringFlag{Name: "rpc_url", Required: true, Usage: "ETH node RPC URL", EnvVars: []string{"RPC_URL"}},
					&cli.StringFlag{Name: "goerli_rpc_url", Required: true, Usage: "Goerli ETH node RPC URL", EnvVars: []string{"GOERLI_RPC_URL"}},
					&cli.StringFlag{Name: "db_url", Required: true, Usage: "Database connection string", EnvVars: []string{"DATABASE_URL"}},
					&cli.StringFlag{Name: "token_addr", Value: "0xCF39360b26a7E54f6c456E69640671Fc5e774FA2", Usage: "Set the token addr (mainnet)", EnvVars: []string{"TOKEN_ADDR"}},
					&cli.StringFlag{Name: "goerli_token_addr", Value: "0xfF30d2c046AEb5FA793138265Cc586De814d0040", Usage: "Set the token addr (goerli)", EnvVars: []string{"GOERLI_TOKEN_ADDR"}},
					&cli.BoolFlag{Name: "scrape_mainnet_eth", Value: true, Usage: "Scrape mainnet eth txes", EnvVars: []string{"SCRAPE_MAINNET_ETH"}},
					&cli.BoolFlag{Name: "scrape_mainnet_sups", Value: true, Usage: "Scrape mainnet sups txes", EnvVars: []string{"SCRAPE_MAINNET_SUPS"}},
					&cli.BoolFlag{Name: "scrape_goerli_eth", Value: true, Usage: "Scrape goerli eth txes", EnvVars: []string{"SCRAPE_GOERLI_ETH"}},
					&cli.BoolFlag{Name: "scrape_goerli_sups", Value: true, Usage: "Scrape goerli sups txes", EnvVars: []string{"SCRAPE_GOERLI_SUPS"}},
				},
				Action: func(c *cli.Context) error {
					logFormat := c.String("log_format")
					ttlSeconds := c.Int("ttl_seconds")
					rpcURL := c.String("rpc_url")
					goerliRpcUrl := c.String("goerli_rpc_url")
					port := c.Int("port")
					dbURL := c.String("db_url")
					tokenAddr := c.String("token_addr")
					goerliTokenAddr := c.String("goerli_token_addr")
					log = zerolog.New(os.Stdout).With().Caller().Logger()
					if logFormat == "console" {
						log = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
					}
					err := Connect(dbURL)
					if err != nil {
						return fmt.Errorf("connect db: %w", err)
					}
					mainnetClient, err := ethclient.Dial(rpcURL)
					if err != nil {
						return fmt.Errorf("dial eth node %s: %w", rpcURL, err)
					}
					goerliClient, err := ethclient.Dial(goerliRpcUrl)
					if err != nil {
						return fmt.Errorf("dial goerli eth node %s: %w", rpcURL, err)
					}

					ethusdAddr := common.HexToAddress("0x5f4eC3Df9cbd43714FE2740f5E3616155c5b8419")
					bnbethAddr := common.HexToAddress("0x14e613ac84a31f709eadbdf89c6cc390fdc9540a")
					supethAddr := common.HexToAddress("0xa1e5dc01359c2920c096f0091fc7f0bf69812ca7")

					bnbethContract, err := ethusd.NewEthusd(bnbethAddr, mainnetClient)
					if err != nil {
						return fmt.Errorf("create ethusd contract: %w", err)
					}
					ethusdContract, err := ethusd.NewEthusd(ethusdAddr, mainnetClient)
					if err != nil {
						return fmt.Errorf("create ethusd contract: %w", err)
					}
					supsethContract, err := supseth.NewSupseth(supethAddr, mainnetClient)
					if err != nil {
						return fmt.Errorf("create supseth contract: %w", err)
					}
					ethC := &EthClient{mainnetClient, ethusdContract, supsethContract, bnbethContract}
					if err != nil {
						return fmt.Errorf("connect db: %w", err)
					}

					t := &Tickers{
						c.Bool("scrape_mainnet_eth"),
						c.Bool("scrape_mainnet_sups"),
						c.Bool("scrape_goerli_eth"),
						c.Bool("scrape_goerli_sups"),
						ethC,
						ethC.Client,
						goerliClient,
						common.HexToAddress(tokenAddr),
						common.HexToAddress(goerliTokenAddr),
					}
					go t.Start()

					s := &Subscriber{mainnetClient, goerliClient, t}
					s.Start()

					return Serve(ethC, rpcURL, port, ttlSeconds)
				},
			},
			{
				Name: "scrape",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "rpc_url", Required: true, Usage: "Infura or Eth node URL", EnvVars: []string{"RPC_URL"}},
					&cli.IntFlag{Name: "from_block", Value: 15879854, Usage: "Set the from block", EnvVars: []string{"FROM_BLOCK"}},
					&cli.IntFlag{Name: "to_block", Value: 15974754, Usage: "Set the to block", EnvVars: []string{"TO_BLOCK"}},
					&cli.IntFlag{Name: "chain_id", Value: 1, Usage: "Set the chain id", EnvVars: []string{"CHAIN_ID"}},
					&cli.IntFlag{Name: "token_decimals", Value: 18, Usage: "Set the token decimals", EnvVars: []string{"TOKEN_DECIMALS"}},
					&cli.StringFlag{Name: "token_addr", Value: "0xCF39360b26a7E54f6c456E69640671Fc5e774FA2", Usage: "Set the token addr", EnvVars: []string{"TOKEN_ADDR"}},
					&cli.StringFlag{Name: "token_symbol", Value: "SUPS", Usage: "Set the token symbol", EnvVars: []string{"TOKEN_SYMBOL"}},
				},
				Usage: "Run SUPS scraper",
				Action: func(c *cli.Context) error {
					rpcUrl := c.String("rpc_url")
					fromBlock := c.Int("from_block")
					toBlock := c.Int("to_block")
					chainId := c.Int("chain_id")
					tokenDecimals := c.Int("token_decimals")
					tokenAddr := c.String("token_addr")
					tokenSymbol := c.String("token_symbol")
					client, err := ethclient.Dial(rpcUrl)
					if err != nil {
						return fmt.Errorf("dial eth node: %w", err)
					}

					log.Info().
						Int("from_block", fromBlock).
						Int("to_block", toBlock).
						Int("chain_id", chainId).
						Int("token_decimals", tokenDecimals).
						Str("token_addr", tokenAddr).
						Str("token_symbol", tokenSymbol).
						Msg("scrape")

					_, err = ScrapeSUPS(client, int64(fromBlock), int64(toBlock), int64(chainId), common.HexToAddress(tokenAddr), tokenSymbol, tokenDecimals)
					return err
				},
			}},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("run cli")
	}

}

func LoggerMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		t1 := time.Now()
		defer func() {
			log.Info().
				Str("ip", r.RemoteAddr).
				Str("method", r.Method).
				Str("path", r.URL.String()).
				Int("status", ww.Status()).
				Int("bytes", ww.BytesWritten()).
				Str("duration", fmt.Sprintf("%dus", int(time.Since(t1)/time.Microsecond))).
				Msg("api call")
		}()

		next.ServeHTTP(ww, r)
	}
	return http.HandlerFunc(fn)
}

func Serve(ethC *EthClient, rpcURL string, port int, ttlSeconds int) error {

	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(10000000),
	)
	if err != nil {
		return fmt.Errorf("memcached middleware: %w", err)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(time.Duration(ttlSeconds)*time.Second),
		cache.ClientWithRefreshKey("opn"),
	)
	if err != nil {
		return fmt.Errorf("memcached client: %w", err)
	}

	c := &Controller{ethC}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(LoggerMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(chiprometheus.NewPatternMiddleware("xsyn-pricefeed"))

	r.Handle("/metrics", promhttp.Handler())
	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		md := docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{ProjectPath: "xsyn-pricefeed"})
		output := markdown.ToHTML([]byte(md), nil, nil)

		w.Write(output)
	})
	r.Get("/api/transfers/{chain}/{symbol}", http.HandlerFunc(c.Transfers))
	r.Get("/api/check", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	r.Get("/api/prices", cacheClient.Middleware(http.HandlerFunc(c.PricesHandler)).ServeHTTP)
	r.Get("/api/eth_price", cacheClient.Middleware(http.HandlerFunc(c.Eth)).ServeHTTP)
	r.Get("/api/bnb_price", cacheClient.Middleware(http.HandlerFunc(c.Bnb)).ServeHTTP)
	r.Get("/api/sups_price", cacheClient.Middleware(http.HandlerFunc(c.Sups)).ServeHTTP)
	log.Info().Int("port", port).Msg("Running server")

	return http.ListenAndServe(":"+fmt.Sprintf("%d", port), r)
}

func (c *Controller) Transfers(w http.ResponseWriter, r *http.Request) {

	symbol := chi.URLParam(r, "symbol")
	chain := chi.URLParam(r, "chain")
	if symbol != "sups" && symbol != "eth" {
		http.Error(w, "sups or eth only", http.StatusInternalServerError)
		return
	}

	var err error
	sinceBlockStr := r.URL.Query().Get("since_block")
	sinceBlock := 0
	if sinceBlockStr != "" {
		sinceBlock, err = strconv.Atoi(sinceBlockStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	blockheight := 0
	chainID := 0
	switch chain {
	case "mainnet":
		chainID = 1
		blockheight, err = GetInt(KeyBlockHeightMainnet, BaseMainnetBlock)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "goerli":
		chainID = 5
		blockheight, err = GetInt(KeyBlockHeightGoerli, BaseGoerliBlock)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	result, err := Transfers(symbol, blockheight, sinceBlock, chainID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type Controller struct {
	*EthClient
}

type SingleResponse struct {
	Time int64  `json:"time"`
	Usd  string `json:"usd"`
}

func (c *Controller) Eth(w http.ResponseWriter, r *http.Request) {
	price, err := c.ETHUSD()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	time := time.Now().Unix()
	resp := &SingleResponse{time, price.Div(decimal.NewFromInt(100)).String()}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) Bnb(w http.ResponseWriter, r *http.Request) {
	price, err := c.BNBUSD()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	time := time.Now().Unix()
	resp := &SingleResponse{time, price.Div(decimal.NewFromInt(100)).String()}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (c *Controller) Sups(w http.ResponseWriter, r *http.Request) {
	price, err := c.SUPSUSD()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	time := time.Now().Unix()
	resp := &SingleResponse{time, price.Div(decimal.NewFromInt(100)).String()}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type PriceResponse struct {
	Time    int64           `json:"time"`
	SUPSUSD decimal.Decimal `json:"sups_usd_cents"`
	ETHUSD  decimal.Decimal `json:"eth_usd_cents"`
	BNBUSD  decimal.Decimal `json:"bnb_usd_cents"`
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
	bnbusd, err := c.BNBUSD()
	if err != nil {
		log.Err(err).Msg("get ethusd price")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := &PriceResponse{time.Now().Unix(), supsusd, ethusd, bnbusd}
	if result.SUPSUSD.IsZero() {
		result.SUPSUSD = decimal.NewFromFloat(0.8)
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Err(err).Msg("marshal json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type EthClient struct {
	Client  *ethclient.Client
	Ethusd  *ethusd.Ethusd
	Supseth *supseth.Supseth
	Bnbusd  *ethusd.Ethusd
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

func (c *EthClient) BNBUSD() (decimal.Decimal, error) {
	result, err := c.Bnbusd.LatestRoundData(&bind.CallOpts{})
	if err != nil {
		return decimal.Zero, fmt.Errorf("query bnbusd: %w", err)
	}
	return decimal.NewFromBigInt(result.Answer, -6), nil
}
