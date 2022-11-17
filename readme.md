# XSYN Pricefeed

A REST API that returns the price of Ether and SUPS.

To regenerate the contract bindings, install `solc` and `abigen`:

```
cd ethusd
solc --abi ethusd.sol -o .
abigen --abi=AggregatorV3Interface.abi --pkg=ethusd --out=ethusd.go

cd ..
cd supseth
solc --abi supseth.sol -o .
abigen --abi=IUniswapV3PoolState.abi.abi --pkg=supseth --out=supseth.go
```

To run:

```
go run main.go --rpc_url {{RPC_URL}}
```
